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
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap internal/event/main/get_event.go && chmod 755 bootstrap && zip  build/package/event/get_event.zip bootstrap && rm bootstrap
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
	startDate, err := time.Parse(time.RFC3339, request.QueryStringParameters["startDate"])
	if err != nil {
		return util.LambdaErrorResponse(util.NewCustomApplicationError("startDate must be provided"), 400), nil
	}
	endDate, err := time.Parse(time.RFC3339, request.QueryStringParameters["endDate"])
	if err != nil {
		return util.LambdaErrorResponse(util.NewCustomApplicationError("endDate must be provided"), 400), nil
	}

	res, svcErr := eventSvc.GetEventBetweenStartAndEndDate(coupleId, startDate, endDate)
	if svcErr != nil {
		return util.LambdaAppErrorResponse(svcErr), nil
	}

	return util.LambdaResponseWithData(res), nil
}

func main() {
	lambda.Start(getEvent)
}
