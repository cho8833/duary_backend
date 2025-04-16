package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/internal/auth"
	"github.com/cho8833/duary_lambda/internal/auth/jwtutil"
	"github.com/cho8833/duary_lambda/internal/common"
	"github.com/cho8833/duary_lambda/internal/couple"
	"github.com/cho8833/duary_lambda/internal/member"
	"github.com/cho8833/duary_lambda/internal/util"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap internal/common/main/check_connected_api.go && chmod 755 bootstrap && zip  build/package/common/check_connected_api.zip bootstrap && rm bootstrap
*/
func checkConnected(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cacheClient := util.GetCacheClient()
	dynamodbClient, _ := cacheClient.GetDynamoDBClient()

	memberRepo := member.NewMemberRepository(dynamodbClient)
	coupleRepo := couple.NewCoupleRepository(dynamodbClient)

	memberSvc := member.NewMemberService(memberRepo)
	coupleSvc := couple.NewCoupleService(coupleRepo)

	commonSvc := common.NewCommonService(memberSvc, coupleSvc)

	loginMember := auth.FromRequestContext(request)

	res, svcError := commonSvc.CheckConnected(loginMember)
	if svcError != nil {
		return util.LambdaAppErrorResponse(svcError), nil
	}

	if res.Couple == nil {
		return util.LambdaResponseWithData(nil), nil
	}

	return util.LambdaResponseWithDataAndHeader(res, jwtutil.ApplicationJWTToHeader(*res.Token)), nil
}

func main() {
	lambda.Start(checkConnected)
}
