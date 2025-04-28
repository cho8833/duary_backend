package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/internal/event"
	"github.com/cho8833/duary_lambda/internal/util"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap internal/event/main/delete_event_api.go && chmod 755 bootstrap && zip  build/package/event/delete_event.zip bootstrap && rm bootstrap
*/
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

	util.NewAuthContext(request)

	svcErr := eventSvc.DeleteEvent(coupleId, eventId)
	if svcErr != nil {
		return util.LambdaAppErrorResponse(svcErr), nil
	}

	return util.LambdaResponseWithData(nil), nil
}

func main() {
	lambda.Start(deleteEvent)
}
