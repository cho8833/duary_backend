package main

import (
	"context"
	"dummy_sign_in"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/appjwt"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/shared"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/dummy_sign_in/main/main.go && chmod 755 bootstrap && zip  build/package/dummy_sign_in_api.zip bootstrap && rm bootstrap
*/
func dummySignIn(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamodbClient, _ := model.GetDynamoDBClient()

	memberRepo := member.NewRepository(dynamodbClient)
	coupleRepo := couple.NewRepository(dynamodbClient)

	memberSvc := member.NewService(memberRepo)
	coupleSvc := couple.NewService(coupleRepo)

	jwtUtil := &appjwt.Impl{}

	req := &dummy_sign_in.DummySignInReq{}
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	res, svcError := dummy_sign_in.SignIn(req, memberSvc, coupleSvc, jwtUtil)
	if svcError != nil {
		return shared.LambdaAppErrorResponse(svcError), nil
	}

	return shared.LambdaResponseWithDataAndHeader(res, appjwt.ApplicationJWTToHeader(*res.Token)), nil
}

func main() {
	lambda.Start(dummySignIn)
}
