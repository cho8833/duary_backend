package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/internal/connectcouple"
	"github.com/cho8833/duary_lambda/internal/util"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap cmd/auth/main/on_disconnect_couple_ws.go && chmod 755 bootstrap && zip  build/package/on_disconnect_couple_ws.zip bootstrap && rm bootstrap
*/

func onDisconnectCoupleWS(_ context.Context, req *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	cacheClient := util.GetCacheClient()
	dynamodbClient, err := cacheClient.GetDynamoDBClient()
	if err != nil {
		return util.LambdaAppErrorResponse(util.InternalServerError{}), err
	}
	repository := connectcouple.NewConnectCoupleRepository(dynamodbClient)
	handler := connectcouple.NewConnectCoupleService(repository)
	log.Printf("delete session connectionId :%s", req.RequestContext.ConnectionID)
	svcError := handler.DeleteSession(&req.RequestContext.ConnectionID)

	if svcError != nil {
		return util.LambdaAppErrorResponse(svcError), nil
	}
	return util.LambdaResponseWithData(nil), nil
}

func main() {
	lambda.Start(onDisconnectCoupleWS)
}
