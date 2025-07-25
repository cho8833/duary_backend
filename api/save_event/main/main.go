package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/event"
	"github.com/cho8833/duary_lambda/scheduler"
	"github.com/cho8833/duary_lambda/shared"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/save_event/main/main.go && chmod 755 bootstrap && zip  build/package/event/save_event.zip bootstrap && rm bootstrap
*/

func saveEvent(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//----------------------------init----------------------------------------//
	dynamoDBClient, err := model.GetDynamoDBClient()
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	schedulerClient, err := scheduler.GetSchedulerClient()
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	scheduleHelper := scheduler.NewEventBridgeSchedulerHelper(schedulerClient)
	eventRepo := event.NewRepository(dynamoDBClient)
	coupleRepo := couple.NewRepository(dynamoDBClient)
	eventSvc := event.NewService(eventRepo)
	coupleSvc := couple.NewService(coupleRepo)

	authContext := shared.NewAuthContext(req)
	saveEventReq := &event.SaveReq{}
	err = json.Unmarshal([]byte(req.Body), &saveEventReq)
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}
	saveEventReq.CoupleId = *authContext.CoupleId
	saveEventReq.CreatedBy = *authContext.SocialId + "-" + *authContext.Provider

	//----------------------------Find Couple----------------------------------------//
	findCouple, svcError := coupleSvc.FindById(*authContext.CoupleId)
	if svcError != nil {
		log.Printf(svcError.Error())
		return shared.LambdaAppErrorResponse(shared.CoupleNotFound{}), nil
	}

	//----------------------------Save Event(DynamoDB)-------------------------------//
	vo, svcError := eventSvc.Save(saveEventReq)
	if svcError != nil {
		return shared.LambdaAppErrorResponse(svcError), nil
	}

	//----------------------------Create EventBridge Schedule------------------------//
	for _, user := range findCouple.Members {
		svcError = scheduleHelper.CreateEventSchedule(*vo, *user)
		if svcError != nil {
			// TODO: 에러 처리
			log.Printf(svcError.Error())
		}
	}

	return shared.LambdaResponseWithData(vo), nil
}

func main() {
	lambda.Start(saveEvent)
}
