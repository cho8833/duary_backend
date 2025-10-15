package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_backend/model"
	"github.com/cho8833/duary_backend/model/couple"
	"github.com/cho8833/duary_backend/model/event"
	"github.com/cho8833/duary_backend/scheduler"
	"github.com/cho8833/duary_backend/shared"
	"log"
	"time"
)

/*
GCO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/update_couple/main.go && chmod 755 bootstrap && zip build/package/update_couple.zip bootstrap && rm bootstrap
*/

func updateCouple(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	authContext := shared.NewAuthContext(request)
	coupleId := authContext.CoupleId
	if coupleId == nil {
		log.Printf("coupleId is nil")
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}
	socialId := authContext.SocialId
	if socialId == nil {
		log.Printf("socialId is nil")
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}
	provider := authContext.Provider
	if provider == nil {
		log.Printf("provider is nil")
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	dbClient, err := model.GetDynamoDBClient()
	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	schedulerClient, err := scheduler.GetSchedulerClient()
	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	stage := request.QueryStringParameters["stage"]

	schedulerHelper := scheduler.NewEventBridgeSchedulerHelper(schedulerClient)

	coupleRepo := couple.NewRepository(dbClient, stage)
	eventRepo := event.NewRepository(dbClient, stage)
	coupleSvc := couple.NewService(coupleRepo)
	eventSvc := event.NewService(eventRepo)

	reqString := request.QueryStringParameters["relationDate"]
	updateRelationDate, err := time.Parse(time.RFC3339, reqString)
	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	transaction := model.NewWriteTransaction(dbClient)
	transaction.BeginTransaction()

	// 커플 정보 업데이트
	req := &couple.UpdateCoupleReq{
		RelationDate: &updateRelationDate,
		Id:           authContext.CoupleId,
	}
	updatedCouple, svcErr := coupleSvc.Update(req, transaction)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	// 커플이 연결된 경우에만 Event 를 재생성함
	var yearlyVO *event.VO
	var day100VO *event.VO
	if len(updatedCouple.Members) > 1 {
		// 기념일 제거
		// 기념일 Event 는 DB 에 없을 수도 있지만, 삭제 대상 Resource 가 없다고 해서 Error 가 발생하지 않음
		eventType := event.Anniversary
		svcErr = eventSvc.DeleteByCoupleIdAndTypeTransaction(*coupleId, eventType, transaction)
		if svcErr != nil {
			return shared.LambdaAppErrorResponse(svcErr), nil
		}
		loginMemberId := *socialId + "-" + *provider

		// 기념일 재생성
		firstMetDaySaveReq := eventSvc.GenerateFirstMetDay(*updatedCouple.Id, loginMemberId, *updatedCouple.RelationDate)
		yearlyAnniversarySaveReq := eventSvc.GenerateYearlyAnniversary(*updatedCouple.Id, loginMemberId, *updatedCouple.RelationDate)
		day100AnniversarySaveReq := eventSvc.Generate100Anniversary(*updatedCouple.Id, loginMemberId, *updatedCouple.RelationDate)

		yearlyVO, svcErr = eventSvc.SaveTransaction(yearlyAnniversarySaveReq, transaction)
		if svcErr != nil {
			return shared.LambdaAppErrorResponse(svcErr), nil
		}
		day100VO, svcErr = eventSvc.SaveTransaction(day100AnniversarySaveReq, transaction)
		if svcErr != nil {
			return shared.LambdaAppErrorResponse(svcErr), nil
		}
		_, svcErr = eventSvc.SaveTransaction(firstMetDaySaveReq, transaction)
		if svcErr != nil {
			return shared.LambdaAppErrorResponse(svcErr), nil
		}
	}

	output, err := transaction.Execute()
	if err != nil {
		log.Printf(err.Error())
		log.Printf(shared.ToString(output))
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	// 기념일 알림 schedule 업데이트
	if day100VO != nil {
		_ = schedulerHelper.UpdateAnniversarySchedule(*day100VO, updatedCouple.Members, stage)
	}
	if yearlyVO != nil {
		_ = schedulerHelper.UpdateAnniversarySchedule(*yearlyVO, updatedCouple.Members, stage)
	}

	return shared.LambdaResponseWithData(updatedCouple), nil
}

func main() {
	lambda.Start(updateCouple)
}
