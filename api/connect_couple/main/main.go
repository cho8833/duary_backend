package main

import (
	"connect_couple"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/appjwt"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/event"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/scheduler"
	"github.com/cho8833/duary_lambda/shared"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/connect_couple/main/main.go && chmod 755 bootstrap && zip  build/package/connect_couple.zip bootstrap && rm bootstrap
*/
func handler(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamodbClient, err := model.GetDynamoDBClient()
	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	schedulerClient, err := scheduler.GetSchedulerClient()
	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	schedulerHelper := scheduler.NewEventBridgeSchedulerHelper(schedulerClient)

	transaction := model.NewWriteTransaction(dynamodbClient)
	memberRepo := member.NewRepository(dynamodbClient)
	coupleRepo := couple.NewRepository(dynamodbClient)
	eventRepo := event.NewRepository(dynamodbClient)

	memberSvc := member.NewService(memberRepo)
	coupleSvc := couple.NewService(coupleRepo)
	eventSvc := event.NewService(eventRepo)

	connectCoupleReq := &connect_couple.ConnectCoupleReq{}
	err = json.Unmarshal([]byte(request.Body), &connectCoupleReq)
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	shared.NewAuthContext(request)

	res, svcError := connect_couple.ConnectCouple(connectCoupleReq, transaction, coupleSvc, memberSvc, eventSvc, *schedulerHelper)
	if svcError != nil {
		return shared.LambdaAppErrorResponse(svcError), nil
	}

	return shared.LambdaResponseWithDataAndHeader(res, appjwt.ApplicationJWTToHeader(*res.Token)), nil
}

func main() {
	lambda.Start(handler)
}
