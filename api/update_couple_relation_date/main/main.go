package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/event"
	"github.com/cho8833/duary_lambda/scheduler"
	"github.com/cho8833/duary_lambda/shared"
	"log"
	"time"
)

/*
GCO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/update_couple_relation_date/main/main.go && chmod 755 bootstrap && zip build/package/update_couple_relation_date.zip bootstrap && rm bootstrap
*/

func updateCoupleRelationDate(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

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

	schedulerHelper := scheduler.NewEventBridgeSchedulerHelper(schedulerClient)

	coupleRepo := couple.NewRepository(dbClient)
	eventRepo := event.NewRepository(dbClient)
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

	// 기념일 제거 후 재생성
	eventType := event.Anniversary
	svcErr = eventSvc.DeleteByCoupleIdAndTypeTransaction(*coupleId, eventType, transaction)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}
	loginMemberId := *socialId + "-" + *provider

	firstMetDaySaveReq := eventSvc.GenerateFirstMetDay(*updatedCouple.Id, loginMemberId, *updatedCouple.RelationDate)
	yearlyAnniversarySaveReq := eventSvc.GenerateYearlyAnniversary(*updatedCouple.Id, loginMemberId, *updatedCouple.RelationDate)
	day100AnniversarySaveReq := eventSvc.Generate100Anniversary(*updatedCouple.Id, loginMemberId, *updatedCouple.RelationDate)

	yearlyVO, svcErr := eventSvc.SaveTransaction(yearlyAnniversarySaveReq, transaction)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}
	day100VO, svcErr := eventSvc.SaveTransaction(day100AnniversarySaveReq, transaction)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}
	_, svcErr = eventSvc.SaveTransaction(firstMetDaySaveReq, transaction)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	output, err := transaction.Execute()
	if err != nil {
		log.Printf(err.Error())
		log.Printf(shared.ToString(output))
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	// 기념일 알림 schedule 업데이트
	svcErr = schedulerHelper.UpdateAnniversarySchedule(*day100VO, updatedCouple.Members)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}
	svcErr = schedulerHelper.UpdateAnniversarySchedule(*yearlyVO, updatedCouple.Members)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	return shared.LambdaResponseWithData(updatedCouple), nil
}

func main() {
	lambda.Start(updateCoupleRelationDate)
}
