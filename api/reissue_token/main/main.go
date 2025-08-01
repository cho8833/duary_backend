package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/appjwt"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/shared"
	"os"
	"strings"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/reissue_token/main/main.go && chmod 755 bootstrap && zip  build/package/reissue_token.zip bootstrap && rm bootstrap
*/

type ReissueTokenRequest struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func reissueTokenAPI(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	key := os.Getenv("secretKey")
	jwtUtil := appjwt.Impl{}

	req := &ReissueTokenRequest{}
	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}
	jwtInfo, err := jwtUtil.ValidateApplicationJWT(req.RefreshToken, key)
	if err != nil {
		return shared.LambdaErrorResponse(fmt.Errorf("세션이 만료되었습니다"), 400), nil
	}

	subSplit := strings.Split(*jwtInfo.Sub, "-")

	dbClient, err := model.GetDynamoDBClient()
	if err != nil {
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	memberRepo := member.NewRepository(dbClient)
	memberSvc := member.NewService(memberRepo)
	loggedInMember, svcErr := memberSvc.FindById(subSplit[0], subSplit[1])
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	applicationJWT := jwtUtil.NewToken(loggedInMember.GetId(), jwtInfo.CoupleId, key)
	return shared.LambdaResponseWithDataAndHeader(applicationJWT, appjwt.ApplicationJWTToHeader(*applicationJWT)), nil
}

func main() {
	lambda.Start(reissueTokenAPI)
}
