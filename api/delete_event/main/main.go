package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/event"
	"github.com/cho8833/duary_lambda/shared"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/delete_event/main/main.go && chmod 755 bootstrap && zip  build/package/event/delete_event.zip bootstrap && rm bootstrap
*/
func deleteEvent(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	dynamodbClient, err := model.GetDynamoDBClient()

	if err != nil {
		return shared.LambdaAppErrorResponse(shared.DBError{}), nil
	}

	eventRepo := event.NewRepository(dynamodbClient)
	eventSvc := event.NewService(eventRepo)

	coupleId := request.QueryStringParameters["coupleId"]
	eventId := request.QueryStringParameters["id"]

	if coupleId == "" || eventId == "" {
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	shared.NewAuthContext(request)

	svcErr := eventSvc.DeleteEvent(coupleId, eventId)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	return shared.LambdaResponseWithData(nil), nil
}

func main() {
	lambda.Start(deleteEvent)
}
