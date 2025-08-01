package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/event"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/scheduler"
	"github.com/cho8833/duary_lambda/shared"
	"log"
	"update_member"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/update_member/main/main.go && chmod 755 bootstrap && zip  build/package/update_member.zip bootstrap && rm bootstrap
*/
func updateMember(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

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

	memberRepo := member.NewRepository(dynamodbClient)
	memberSvc := member.NewService(memberRepo)
	coupleRepo := couple.NewRepository(dynamodbClient)
	coupleSvc := couple.NewService(coupleRepo)
	eventRepo := event.NewRepository(dynamodbClient)
	eventSvc := event.NewService(eventRepo)

	updateReq := &update_member.UpdateMemberReq{}
	err = json.Unmarshal([]byte(request.Body), updateReq)
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	transaction := model.NewWriteTransaction(dynamodbClient)

	shared.NewAuthContext(request)

	result, svcErr := update_member.UpdateMember(updateReq, transaction, coupleSvc, memberSvc, eventSvc, *schedulerHelper)

	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}
	return shared.LambdaResponseWithData(result), nil
}

func main() {
	lambda.Start(updateMember)
}
