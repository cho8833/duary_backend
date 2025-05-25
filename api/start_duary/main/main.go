package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/appjwt"
	"github.com/cho8833/duary_lambda/couple"
	"github.com/cho8833/duary_lambda/member"
	"github.com/cho8833/duary_lambda/shared"
	"log"
	"start_duary"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/start_duary/main/main.go && chmod 755 bootstrap && zip  build/package/common/start_duary_api.zip bootstrap && rm bootstrap
*/
func startDuaryAPI(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// init
	dynamoDBClient, err := member.GetDynamoDBClient()
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	transaction := shared.NewWriteTransaction(dynamoDBClient)
	coupleRepo := couple.NewCoupleRepository(dynamoDBClient)
	memberRepo := member.NewMemberRepository(dynamoDBClient)
	coupleSvc := couple.NewCoupleService(coupleRepo)
	memberSvc := member.NewMemberService(memberRepo)

	initDuaryReq := &start_duary.StartDuaryReq{}
	err = json.Unmarshal([]byte(req.Body), &initDuaryReq)
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	shared.NewAuthContext(req)

	result, svcError := start_duary.StartDuary(initDuaryReq, transaction, coupleSvc, memberSvc)

	if svcError != nil {
		return shared.LambdaAppErrorResponse(svcError), nil
	}

	return shared.LambdaResponseWithDataAndHeader(result, appjwt.ApplicationJWTToHeader(*result.Token)), nil
}

func main() {
	lambda.Start(startDuaryAPI)
}
