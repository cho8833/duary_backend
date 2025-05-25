package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/fcm"
	"github.com/cho8833/duary_lambda/shared"
	"log"
)

/*
GCO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/send_notification/main/main.go && chmod 755 bootstrap && zip  build/package/common/send_notification.zip bootstrap duary-8c5b2-firebase-adminsdk-9d1a5-e86abbedfc.json && rm bootstrap
*/
func sendNotification(ctx context.Context, event json.RawMessage) (events.APIGatewayProxyResponse, error) {
	fcmClient, err := fcm.GetFCMClient()

	if err != nil {
		log.Printf("failed to get fcmClient: %+v\n", err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	fcmReq := &fcm.SendReq{}
	err = json.Unmarshal(event, fcmReq)
	if err != nil {
		log.Printf("failed to get req body: %+v\n\n", err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	fcmService := fcm.NewService(fcmClient)

	fcmService.Send(*fcmReq)

	return shared.LambdaResponseWithData(nil), nil
}

func main() {
	lambda.Start(sendNotification)
}
