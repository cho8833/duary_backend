package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/internal/event"
	"github.com/cho8833/duary_lambda/internal/util"
	"time"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap cmd/auth/main/get_event.go && chmod 755 bootstrap && zip  build/package/get_event.zip bootstrap && rm bootstrap
*/
func getEvent(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	cacheClient := util.GetCacheClient()
	dynamodbClient, err := cacheClient.GetDynamoDBClient()

	if err != nil {
		return util.LambdaAppErrorResponse(util.DBError{}), nil
	}

	eventRepo := event.NewEventRepository(dynamodbClient)
	eventSvc := event.NewEventService(eventRepo)

	coupleId := request.QueryStringParameters["coupleId"]
	startDate, _ := time.Parse(time.RFC3339, request.QueryStringParameters["startDate"])
	endDate, _ := time.Parse(time.RFC3339, request.QueryStringParameters["endDate"])

	result, svcErr := eventSvc.GetEventBetweenStartAndEndDate(coupleId, startDate, endDate)
	if svcErr != nil {
		return util.LambdaAppErrorResponse(svcErr), nil
	}

	return util.LambdaResponseWithData(result), nil
}

func main() {
	lambda.Start(getEvent)
}
