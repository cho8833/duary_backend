package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/appjwt"
	"github.com/cho8833/duary_lambda/auth"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/shared"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/apple_sign_in/main/main.go && chmod 755 bootstrap && zip  build/package/apple_sign_in_api.zip bootstrap && rm bootstrap
*/
func appleSignIn(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamoDBClient, err := model.GetDynamoDBClient()

	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	memberRepository := member.NewMemberRepository(dynamoDBClient)
	coupleRepository := couple.NewCoupleRepository(dynamoDBClient)

	memberSvc := member.NewMemberService(memberRepository)
	coupleSvc := couple.NewCoupleService(coupleRepository)

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
