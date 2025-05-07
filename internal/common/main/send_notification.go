package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/internal/util"
	"log"
)

/*
GCO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap internal/common/main/send_notification.go && chmod 755 bootstrap && zip  build/package/common/send_notification.zip bootstrap duary-8c5b2-firebase-adminsdk-9d1a5-e86abbedfc.json && rm bootstrap
*/
func sendNotification(ctx context.Context, event json.RawMessage) (events.APIGatewayProxyResponse, error) {

	cacheClient := util.GetCacheClient()
	fcmClient, err := cacheClient.GetFCMClient()

	if err != nil {
		log.Printf("failed to get fcmClient: %+v\n", err)
		return util.LambdaAppErrorResponse(util.InternalServerError{}), nil
	}

	fcmReq := &util.FCMReq{}
	err = json.Unmarshal(event, fcmReq)
	if err != nil {
		log.Printf("failed to get req body: %+v\n\n", err)
		return util.LambdaAppErrorResponse(util.InternalServerError{}), nil
	}

	fcmService := util.FcmService(fcmClient)

	fcmService.Send(*fcmReq)

	return util.LambdaResponseWithData(nil), nil
}

func main() {
	lambda.Start(sendNotification)
}
