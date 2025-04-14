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
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap internal/event/main/save_event_api.go && chmod 755 bootstrap && zip  build/package/event/save_event_api.zip bootstrap && rm bootstrap
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

	// jwt 의 coupleId 와 생성 요청 event 의 coupleId 가 동일하지 않으면 잘못된 요청
	lambdaMap := req.RequestContext.Authorizer["lambda"].(map[string]interface{})
	coupleId := lambdaMap["coupleId"].(string)
	if saveEventReq.CoupleId != coupleId {
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
