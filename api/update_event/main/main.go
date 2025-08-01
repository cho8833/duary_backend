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
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/update_event/main/main.go && chmod 755 bootstrap && zip  build/package/update_event.zip bootstrap && rm bootstrap
*/
func updateEvent(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// init
	dynamodbClient, err := model.GetDynamoDBClient()
	if err != nil {
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	eventRepo := event.NewRepository(dynamodbClient)
	eventSvc := event.NewService(eventRepo)

	// get auth
	authCtx := shared.NewAuthContext(request)
	coupleId := authCtx.CoupleId
	if coupleId == nil {
		log.Printf("coupleId is null")
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}
	socialId := authCtx.SocialId
	if socialId == nil {
		log.Printf("socialId is null")
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}
	provider := authCtx.Provider
	if provider == nil {
		log.Printf("provider is null")
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	// parse req
	updateEventReq := &event.EditReq{}
	err = json.Unmarshal([]byte(request.Body), updateEventReq)
	if err != nil {
		log.Println(err.Error())
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}
	eventId := request.QueryStringParameters["id"]

	// update event (dynamodb)
	vo, svcError := eventSvc.Update(*coupleId, eventId, updateEventReq)
	if svcError != nil {
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	return shared.LambdaResponseWithData(vo), nil
}

func main() {
	lambda.Start(updateEvent)
}
