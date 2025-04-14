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
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap internal/auth/main/kakao_sign_in_api.go && chmod 755 bootstrap && zip  build/package/auth/kakao_sign_in_api.zip bootstrap && rm bootstrap
*/
func kakaoSignInAPI(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// init
	cacheClient := util.GetCacheClient()
	dynamoDBClient, err := cacheClient.GetDynamoDBClient()
	if err != nil {
		log.Printf(err.Error())
		return util.LambdaAppErrorResponse(util.InternalServerError{}), nil
	}
	memberRepository := member.NewMemberRepository(dynamoDBClient)
	svc := auth.NewKakaoAuthService(&jwtutil.JWTValidatorImpl{}, &jwtutil.Impl{}, memberRepository)

	// parse request
	kakaoToken := &auth.KakaoOAuthToken{}
	err = json.Unmarshal([]byte(request.Body), &kakaoToken)
	if err != nil {
		log.Printf(err.Error())
		return util.LambdaAppErrorResponse(util.BadRequestError{}), nil
	}

	// process
	result, svcError := svc.SignIn(kakaoToken)
	if svcError != nil {
		return util.LambdaAppErrorResponse(svcError), nil
	}

	return util.LambdaResponseWithDataAndHeader(result, jwtutil.ApplicationJWTToHeader(*result.Token)), nil
}

func main() {
	lambda.Start(kakaoSignInAPI)
}
