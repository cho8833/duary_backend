package main

import (
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
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/delete_event/main/main.go && chmod 755 bootstrap && zip  build/package/event/delete_event.zip bootstrap && rm bootstrap
*/
func deleteEvent(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	dynamodbClient, err := model.GetDynamoDBClient()
	if err != nil {
		return shared.LambdaAppErrorResponse(shared.DBError{}), nil
	}

	scheduleClient, err := scheduler.GetSchedulerClient()
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	scheduleHelper := scheduler.NewEventBridgeSchedulerHelper(scheduleClient)

	eventRepo := event.NewRepository(dynamodbClient)
	eventSvc := event.NewService(eventRepo)
	coupleRepo := couple.NewRepository(dynamodbClient)
	coupleSvc := couple.NewService(coupleRepo)

	coupleId := request.QueryStringParameters["coupleId"]
	eventId := request.QueryStringParameters["id"]

	if coupleId == "" || eventId == "" {
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	shared.NewAuthContext(request)

	//----------------------------Find Couple----------------------------------------//
	findCouple, svcError := coupleSvc.FindById(coupleId)
	if svcError != nil {
		log.Printf(svcError.Error())
		return shared.LambdaAppErrorResponse(shared.CoupleNotFound{}), nil
	}

	//----------------------------Delete Event(DynamoDB)-------------------------------//
	svcErr := eventSvc.DeleteEvent(coupleId, eventId)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	//----------------------------Delete EventBridge Schedule------------------------//
	for _, m := range findCouple.Members {
		scheduleName := scheduleHelper.GetScheduleName(eventId, m.GetId())
		svcError = scheduleHelper.DeleteEventSchedule(scheduleName)
		if svcError != nil {
			// 삭제 실패해도 상관없을듯
			log.Printf(svcError.Error())
		}
	}

	return shared.LambdaResponseWithData(nil), nil
}

func main() {
	lambda.Start(deleteEvent)
}
