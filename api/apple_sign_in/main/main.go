package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/appjwt"
	"github.com/cho8833/duary_lambda/auth"
	"github.com/cho8833/duary_lambda/member"
	"github.com/cho8833/duary_lambda/shared"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/apple_sign_in/main/main.go && chmod 755 bootstrap && zip  build/package/auth/apple_sign_in_api.zip bootstrap && rm bootstrap
*/
func appleSignIn(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamoDBClient, err := auth.GetDynamoDBClient()

	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	memberRepository := member.NewMemberRepository(dynamoDBClient)
	svc := auth.NewAuthService(&appjwt.JWTValidatorImpl{}, &appjwt.Impl{}, memberRepository)

	appleToken := &auth.AppleOAuthToken{}
	err = json.Unmarshal([]byte(request.Body), &appleToken)
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	result, svcError := svc.AppleSignIn(appleToken)
	if svcError != nil {
		return shared.LambdaAppErrorResponse(svcError), nil
	}

	return shared.LambdaResponseWithDataAndHeader(result, appjwt.ApplicationJWTToHeader(*result.Token)), nil

}

func main() {
	lambda.Start(appleSignIn)
}
