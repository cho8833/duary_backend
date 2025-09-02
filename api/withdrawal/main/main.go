package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/event"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/scheduler"
	"github.com/cho8833/duary_lambda/shared"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/withdrawal/main/main.go && chmod 755 bootstrap && zip  build/package/withdrawal.zip bootstrap && rm bootstrap
*/
// TODO: 소셜 로그인 연결 해제 필요
func withdrawal(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	authContext := shared.NewAuthContext(request)

	socialId := authContext.SocialId
	provider := authContext.Provider
	coupleId := authContext.CoupleId

	schedulerClient, err := scheduler.GetSchedulerClient()
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	dbClient, err := model.GetDynamoDBClient()
	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	memberRepo := member.NewRepository(dbClient)
	memberSvc := member.NewService(memberRepo)
	coupleRepo := couple.NewRepository(dbClient)
	coupleSvc := couple.NewService(coupleRepo)
	eventRepo := event.NewRepository(dbClient)
	eventSvc := event.NewService(eventRepo)

	schedulerHelper := scheduler.NewEventBridgeSchedulerHelper(schedulerClient)

	targetCouple, svcErr := coupleSvc.FindById(*coupleId)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	//******************* disconnect couple 과 거의 비슷하나 remove member.Couple -> delete Member 가 다름 ******************
	transaction := model.NewWriteTransaction(dbClient)
	transaction.BeginTransaction()

	// delete member
	svcErr = memberSvc.DeleteTransaction(*socialId, *provider, transaction)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	// delete couple if connectedMember is alone,
	//else remove requested member from couple
	if len(targetCouple.ConnectedMemberIds) == 1 {
		svcErr = coupleSvc.Delete(*targetCouple.Id, transaction)
		if svcErr != nil {
			return shared.LambdaAppErrorResponse(svcErr), nil
		}
	} else {
		// remove requested Member from couple
		var remainMember []member.Member
		for _, m := range targetCouple.Members {
			if m.SocialId != *socialId {
				remainMember = append(remainMember, m)
			}
		}
		if len(remainMember) == 0 {
			remainMember = make([]member.Member, 0)
		}
		updateCoupleReq := &couple.UpdateCoupleReq{
			Id:      targetCouple.Id,
			Members: remainMember,
		}
		_, svcErr = coupleSvc.Update(updateCoupleReq, transaction)
		if svcErr != nil {
			return shared.LambdaAppErrorResponse(svcErr), nil
		}
	}

	// delete anniversary event
	svcErr = eventSvc.DeleteByCoupleIdAndTypeTransaction(*coupleId, event.Anniversary, transaction)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}
	svcErr = eventSvc.DeleteByCoupleIdAndTypeTransaction(*coupleId, event.Birthday, transaction)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	_, err = transaction.Execute()
	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.DBUpdateError{}), nil
	}

	// delete anniversary notification schedule
	svcErr = schedulerHelper.DeleteEventSchedule(schedulerHelper.Get100DayAnniversaryScheduleName(*coupleId))
	if svcErr != nil {
		log.Printf(svcErr.Error())
	}
	svcErr = schedulerHelper.DeleteEventSchedule(schedulerHelper.GetYearlyAnniversaryScheduleName(*coupleId))
	if svcErr != nil {
		log.Printf(svcErr.Error())
	}
	for _, mId := range targetCouple.ConnectedMemberIds {
		svcErr = schedulerHelper.DeleteEventSchedule(schedulerHelper.GetBirthdayScheduleName(*coupleId, mId))
		if svcErr != nil {
			log.Printf(svcErr.Error())
		}
	}

	return shared.LambdaResponseWithData(nil), nil
}

func main() {
	lambda.Start(withdrawal)
}
