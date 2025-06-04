package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/model/event"
	"github.com/cho8833/duary_lambda/shared"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/update_event/main/main.go && chmod 755 bootstrap && zip  build/package/event/update_event_api.zip bootstrap && rm bootstrap
*/
func updateEvent(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamodbClient, err := event.GetDynamoDBClient()

	if err != nil {
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	eventRepo := event.NewRepository(dynamodbClient)

	eventSvc := event.NewService(eventRepo)

	updateEventReq := &event.UpdateReq{}

	err = json.Unmarshal([]byte(request.Body), updateEventReq)

	if err != nil {
		log.Println(err.Error())
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	// jwt 의 coupleId 와 생성 요청 event 의 coupleId 가 동일하지 않으면 잘못된 요청
	lambdaMap := request.RequestContext.Authorizer["lambda"].(map[string]interface{})
	coupleId := lambdaMap["coupleId"].(string)
	if *updateEventReq.CoupleId != coupleId {
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	vo, svcError := eventSvc.Update(updateEventReq)
	if svcError != nil {
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	return shared.LambdaResponseWithData(vo), nil
}

func main() {
	lambda.Start(updateEvent)
}
