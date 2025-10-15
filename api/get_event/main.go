package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_backend/model"
	"github.com/cho8833/duary_backend/model/event"
	"github.com/cho8833/duary_backend/shared"
	"log"
	"time"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/get_event/main.go && chmod 755 bootstrap && zip  build/package/get_event.zip bootstrap && rm bootstrap
*/
func getEvent(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	dynamodbClient, err := model.GetDynamoDBClient()
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.DBError{}), nil
	}

	eventRepo := event.NewRepository(dynamodbClient)

	eventSvc := event.NewService(eventRepo)

	coupleId := request.QueryStringParameters["coupleId"]
	startDate, err := time.Parse(time.RFC3339, request.QueryStringParameters["startDate"])
	if err != nil {
		return shared.LambdaErrorResponse(shared.ValidateError{Message: "startDate must be provided"}, 400), nil
	}
	endDate, err := time.Parse(time.RFC3339, request.QueryStringParameters["endDate"])
	if err != nil {
		return shared.LambdaErrorResponse(shared.ValidateError{Message: "endDate must be provided"}, 400), nil
	}

	res, svcErr := eventSvc.GetBetweenStartAndEndDate(coupleId, startDate, endDate)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	return shared.LambdaResponseWithData(res), nil
}

func main() {
	lambda.Start(getEvent)
}
