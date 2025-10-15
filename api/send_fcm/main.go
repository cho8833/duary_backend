package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_backend/fcm"
	"github.com/cho8833/duary_backend/model"
	"github.com/cho8833/duary_backend/model/member"
	"github.com/cho8833/duary_backend/shared"
	"log"
)

/*
GCO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/send_fcm/main.go && chmod 755 bootstrap && zip build/package/send_fcm.zip bootstrap duary-8c5b2-firebase-adminsdk-9d1a5-e86abbedfc.json && rm bootstrap
*/
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

func main() {
	lambda.Start(sendNotification)
}
