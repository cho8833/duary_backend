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
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap internal/event/main/update_event_api.go && chmod 755 bootstrap && zip  build/package/event/update_event_api.zip bootstrap && rm bootstrap
*/
func updateEvent(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cacheClient := util.GetCacheClient()
	dynamodbClient, err := cacheClient.GetDynamoDBClient()

	if err != nil {
		return util.LambdaAppErrorResponse(util.DBError{}), nil
	}

	eventRepo := event.NewEventRepository(dynamodbClient)
	eventSvc := event.NewEventService(eventRepo)

	updateEventReq := &event.UpdateReq{}

	err = json.Unmarshal([]byte(request.Body), updateEventReq)

	if err != nil {
		log.Println(err.Error())
		return util.LambdaAppErrorResponse(util.BadRequestError{}), nil
	}

	// jwt 의 coupleId 와 생성 요청 event 의 coupleId 가 동일하지 않으면 잘못된 요청
	lambdaMap := request.RequestContext.Authorizer["lambda"].(map[string]interface{})
	coupleId := lambdaMap["coupleId"].(string)
	if *updateEventReq.CoupleId != coupleId {
		return util.LambdaAppErrorResponse(util.BadRequestError{}), nil
	}

	vo, svcError := eventSvc.UpdateEvent(updateEventReq)
	if svcError != nil {
		return util.LambdaAppErrorResponse(util.BadRequestError{}), nil
	}

	return util.LambdaResponseWithData(vo), nil
}

func main() {
	lambda.Start(updateEvent)
}
