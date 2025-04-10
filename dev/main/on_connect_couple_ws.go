package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/internal/connectcouple"
	"github.com/cho8833/duary_lambda/internal/util"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap cmd/auth/main/on_connect_couple_ws.go && chmod 755 bootstrap && zip  build/package/on_connect_couple_ws.zip bootstrap && rm bootstrap
*/
func onConnectCoupleWS(_ context.Context, req *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	cacheClient := util.GetCacheClient()
	dynamodbClient, err := cacheClient.GetDynamoDBClient()
	if err != nil {
		return util.LambdaAppErrorResponse(util.InternalServerError{}), err
	}

	repository := connectcouple.NewConnectCoupleRepository(dynamodbClient)
	handler := connectcouple.NewConnectCoupleService(repository)

	coupleCode := req.QueryStringParameters["coupleCode"]
	memberId := req.QueryStringParameters["memberId"]
	coupleId := req.QueryStringParameters["coupleId"]
	reqData := &connectcouple.SessionReq{
		CoupleCode:   &coupleCode,
		MemberId:     &memberId,
		CoupleId:     &coupleId,
		ConnectionId: &req.RequestContext.ConnectionID,
	}

	session, err := handler.CreateSession(reqData)

	if err != nil {
		return util.LambdaResponseWithData(err), nil
	}
	return util.LambdaResponseWithData(session), nil
}

func main() {
	lambda.Start(onConnectCoupleWS)
}
