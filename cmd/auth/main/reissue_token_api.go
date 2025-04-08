package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/internal/auth/jwtutil"
	"github.com/cho8833/duary_lambda/internal/util"
	"os"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap cmd/auth/main/reissue_token_api.go && chmod 755 bootstrap && zip  build/package/reissue_token_api.zip bootstrap && rm bootstrap
*/

type ReissueTokenRequest struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func reissueTokenAPI(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	key := os.Getenv("secretKey")
	jwtUtil := jwtutil.Impl{}

	req := &ReissueTokenRequest{}
	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return util.LambdaAppErrorResponse(util.BadRequestError{}), nil
	}
	memberId, err := jwtUtil.ValidateApplicationJWT(req.RefreshToken, key)
	if err != nil {
		return util.LambdaErrorResponse(fmt.Errorf("세션이 만료되었습니다"), 400), nil
	}

	applicationJWT := jwtUtil.NewToken(*memberId, key)
	return util.LambdaResponseWithData(applicationJWT), nil
}

func main() {
	lambda.Start(reissueTokenAPI)
}
