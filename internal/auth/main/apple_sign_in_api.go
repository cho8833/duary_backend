package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/internal/auth"
	"github.com/cho8833/duary_lambda/internal/auth/jwtutil"
	"github.com/cho8833/duary_lambda/internal/member"
	"github.com/cho8833/duary_lambda/internal/util"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap internal/auth/main/apple_sign_in_api.go && chmod 755 bootstrap && zip  build/package/auth/apple_sign_in_api.zip bootstrap && rm bootstrap
*/
func appleSignIn(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cacheClient := util.CacheClient{}
	dynamoDBClient, err := cacheClient.GetDynamoDBClient()

	if err != nil {
		log.Printf(err.Error())
		return util.LambdaAppErrorResponse(util.InternalServerError{}), nil
	}

	memberRepository := member.NewMemberRepository(dynamoDBClient)
	svc := auth.NewAuthService(&jwtutil.JWTValidatorImpl{}, &jwtutil.Impl{}, memberRepository)

	appleToken := &auth.AppleOAuthToken{}
	err = json.Unmarshal([]byte(request.Body), &appleToken)
	if err != nil {
		log.Printf(err.Error())
		return util.LambdaAppErrorResponse(util.BadRequestError{}), nil
	}

	result, svcError := svc.AppleSignIn(appleToken)
	if svcError != nil {
		return util.LambdaAppErrorResponse(svcError), nil
	}

	return util.LambdaResponseWithDataAndHeader(result, jwtutil.ApplicationJWTToHeader(*result.Token)), nil

}

func main() {
	lambda.Start(appleSignIn)
}
