package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/internal/event"
	"github.com/cho8833/duary_lambda/internal/util"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap internal/event/main/save_event.go && chmod 755 bootstrap && zip  build/package/event/save_event.zip bootstrap && rm bootstrap
*/

func saveEvent(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cacheClient := util.GetCacheClient()
	dynamoDBClient, err := cacheClient.GetDynamoDBClient()
	if err != nil {
		log.Printf(err.Error())
		return util.LambdaAppErrorResponse(util.InternalServerError{}), nil
	}

	eventRepo := event.NewEventRepository(dynamoDBClient)

	eventSvc := event.NewEventService(eventRepo)

	saveEventReq := &event.SaveReq{}
	err = json.Unmarshal([]byte(req.Body), &saveEventReq)
	if err != nil {
		log.Printf(err.Error())
		return util.LambdaAppErrorResponse(util.BadRequestError{}), nil
	}

	vo, svcError := eventSvc.SaveEvent(saveEventReq)
	if svcError != nil {
		return util.LambdaAppErrorResponse(util.BadRequestError{}), nil
	}

	return util.LambdaResponseWithData(vo), nil
}

func main() {
	lambda.Start(saveEvent)
}
