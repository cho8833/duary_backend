package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_backend/appjwt"
	"github.com/cho8833/duary_backend/auth"
	"github.com/cho8833/duary_backend/model"
	"github.com/cho8833/duary_backend/model/couple"
	"github.com/cho8833/duary_backend/model/member"
	"github.com/cho8833/duary_backend/shared"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap main.go && chmod 755 bootstrap && zip  ../../build/package/apple_sign_in.zip bootstrap && rm bootstrap
*/
func appleSignIn(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamoDBClient, err := model.GetDynamoDBClient()

	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	stage := request.StageVariables["stage"]

	memberRepository := member.NewRepository(dynamoDBClient, stage)
	coupleRepository := couple.NewRepository(dynamoDBClient, stage)

	memberSvc := member.NewService(memberRepository)
	coupleSvc := couple.NewService(coupleRepository)

	svc := auth.NewAuthService(&appjwt.JWTValidatorImpl{}, &appjwt.Impl{}, memberSvc, coupleSvc)

	signInReq := &auth.SignInReq{}
	err = json.Unmarshal([]byte(request.Body), &signInReq)
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	result, svcError := svc.AppleSignIn(signInReq.AppleOAuthToken, signInReq.FcmToken)
	if svcError != nil {
		return shared.LambdaAppErrorResponse(svcError), nil
	}

	return shared.LambdaResponseWithDataAndHeader(result, appjwt.ApplicationJWTToHeader(*result.Token)), nil

}

func main() {
	lambda.Start(appleSignIn)
}
