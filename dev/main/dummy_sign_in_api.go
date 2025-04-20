package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/dev"
	"github.com/cho8833/duary_lambda/internal/auth/jwtutil"
	"github.com/cho8833/duary_lambda/internal/member"
	"github.com/cho8833/duary_lambda/internal/util"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap dev/main/dummy_sign_in_api.go && chmod 755 bootstrap && zip  build/package/dev/dummy_sign_in_api.zip bootstrap && rm bootstrap
*/
func dummySignIn(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cacheClient := util.GetCacheClient()
	dynamodbClient, _ := cacheClient.GetDynamoDBClient()

	memberRepo := member.NewMemberRepository(dynamodbClient)
	jwtUtil := &jwtutil.Impl{}

	dummySvc := dev.NewDummyMemberService(jwtUtil, memberRepo)

	req := &dev.DummySignInReq{}
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf(err.Error())
		return util.LambdaAppErrorResponse(util.BadRequestError{}), nil
	}

	res, svcError := dummySvc.SignIn(req)
	if svcError != nil {
		return util.LambdaAppErrorResponse(svcError), nil
	}

	return util.LambdaResponseWithDataAndHeader(res, jwtutil.ApplicationJWTToHeader(*res.Token)), nil
}

func main() {
	lambda.Start(dummySignIn)
}
