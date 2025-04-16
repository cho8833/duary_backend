package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/internal/auth"
	"github.com/cho8833/duary_lambda/internal/auth/jwtutil"
	"github.com/cho8833/duary_lambda/internal/common"
	"github.com/cho8833/duary_lambda/internal/couple"
	"github.com/cho8833/duary_lambda/internal/member"
	"github.com/cho8833/duary_lambda/internal/util"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap internal/common/main/connect_couple_api.go && chmod 755 bootstrap && zip  build/package/common/connect_couple_api.zip bootstrap && rm bootstrap
*/
func connectCouple(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cacheClient := util.GetCacheClient()
	dynamodbClient, _ := cacheClient.GetDynamoDBClient()

	transaction := util.NewWriteTransaction(dynamodbClient)
	memberRepo := member.NewMemberRepository(dynamodbClient)
	coupleRepo := couple.NewCoupleRepository(dynamodbClient)

	memberSvc := member.NewMemberService(memberRepo)
	coupleSvc := couple.NewCoupleService(coupleRepo)

	commonSvc := common.NewCommonService(memberSvc, coupleSvc)

	loginMember := auth.FromRequestContext(request)

	connectCoupleReq := &common.ConnectCoupleReq{}
	err := json.Unmarshal([]byte(request.Body), &connectCoupleReq)
	if err != nil {
		log.Printf(err.Error())
		return util.LambdaAppErrorResponse(util.BadRequestError{}), nil
	}

	res, svcError := commonSvc.ConnectCouple(loginMember, connectCoupleReq, transaction)
	if svcError != nil {
		return util.LambdaAppErrorResponse(svcError), nil
	}

	return util.LambdaResponseWithDataAndHeader(res, jwtutil.ApplicationJWTToHeader(*res.Token)), nil
}

func main() {
	lambda.Start(connectCouple)
}
