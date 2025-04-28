package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/cho8833/duary_lambda/internal/event"
	"github.com/cho8833/duary_lambda/internal/util"
)

func deleteEvent(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cacheClient := util.GetCacheClient()

	dynamodbClient, err := cacheClient.GetDynamoDBClient()

	if err != nil {
		return util.LambdaAppErrorResponse(util.DBError{}), nil
	}

	eventRepo := event.NewEventRepository(dynamodbClient)
	eventSvc := event.NewEventService(eventRepo)

	coupleId := request.QueryStringParameters["coupleId"]
	eventId := request.QueryStringParameters["id"]

	if coupleId == "" || eventId == "" {
		return util.LambdaAppErrorResponse(util.BadRequestError{}), nil
	}

	svcErr := eventSvc.DeleteEvent(coupleId, eventId)
	if svcErr != nil {
		return util.LambdaAppErrorResponse(svcErr), nil
	}

	return util.LambdaResponseWithData(nil), nil
}
