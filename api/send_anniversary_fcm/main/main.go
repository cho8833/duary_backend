package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/cho8833/duary_lambda/fcm"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/shared"
	"log"
)

func sendNotification(ctx context.Context, event json.RawMessage) (events.APIGatewayProxyResponse, error) {
	fcmClient, err := fcm.GetFCMClient()
	if err != nil {
		log.Printf("failed to get fcmClient: %+v\n", err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	dynamodbClient, err := model.GetDynamoDBClient()
	if err != nil {
		log.Printf("failed to get DynamoDBClient: %+v\n", err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	memberRepo := member.NewRepository(dynamodbClient)
	memberSvc := member.NewService(memberRepo)

	fcmReq := &fcm.SendReq{}
	err = json.Unmarshal(event, fcmReq)
	if err != nil {
		log.Printf("failed to get req body: %+v\n\n", err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	fcmService := fcm.NewService(fcmClient, memberSvc)

	fcmService.Send(*fcmReq)

	return shared.LambdaResponseWithData(nil), nil
}
