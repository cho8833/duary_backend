package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/event"
	"github.com/cho8833/duary_lambda/shared"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/save_event/main/main.go && chmod 755 bootstrap && zip  build/package/event/save_event_api.zip bootstrap && rm bootstrap
*/

func saveEvent(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamoDBClient, err := model.GetDynamoDBClient()
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	eventRepo := event.NewRepository(dynamoDBClient)

	eventSvc := event.NewService(eventRepo)

	authContext := shared.NewAuthContext(req)

	saveEventReq := &event.SaveReq{}

	err = json.Unmarshal([]byte(req.Body), &saveEventReq)
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}
	saveEventReq.CoupleId = *authContext.CoupleId

	saveEventReq.CreatedBy = *authContext.SocialId + "-" + *authContext.Provider

	vo, svcError := eventSvc.Save(saveEventReq)
	if svcError != nil {
		return shared.LambdaAppErrorResponse(svcError), nil
	}

	return shared.LambdaResponseWithData(vo), nil
}

func main() {
	lambda.Start(saveEvent)
}
