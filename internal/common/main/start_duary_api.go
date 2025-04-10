package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/internal/common"
	"github.com/cho8833/duary_lambda/internal/couple"
	"github.com/cho8833/duary_lambda/internal/member"
	"github.com/cho8833/duary_lambda/internal/util"
	"log"
	"strconv"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap internal/auth/main/start_duary_api.go && chmod 755 bootstrap && zip  build/package/common/start_duary_api.zip bootstrap && rm bootstrap
*/
func startDuaryAPI(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// init
	cacheClient := util.GetCacheClient()
	dynamoDBClient, err := cacheClient.GetDynamoDBClient()
	if err != nil {
		log.Printf(err.Error())
		return util.LambdaAppErrorResponse(util.InternalServerError{}), nil
	}
	transaction := util.NewWriteTransaction(dynamoDBClient)
	coupleRepo := couple.NewCoupleRepository(dynamoDBClient)
	memberRepo := member.NewMemberRepository(dynamoDBClient)
	coupleSvc := couple.NewCoupleService(coupleRepo)
	memberSvc := member.NewMemberService(memberRepo)

	commonSvc := common.NewCommonService(memberSvc, coupleSvc, nil)

	initDuaryReq := &common.InitDuaryInfoReq{}
	err = json.Unmarshal([]byte(req.Body), &initDuaryReq)
	if err != nil {
		log.Printf(err.Error())
		return util.LambdaAppErrorResponse(util.BadRequestError{}), nil
	}

	lambdaMap := req.RequestContext.Authorizer["lambda"].(map[string]interface{})
	initDuaryReq.Provider = lambdaMap["provider"].(string)
	initDuaryReq.SocialId, _ = strconv.ParseInt(lambdaMap["socialId"].(string), 10, 64)

	result, svcError := commonSvc.InitDuaryInfo(initDuaryReq, transaction)
	if svcError != nil {
		return util.LambdaAppErrorResponse(svcError), nil
	}
	return util.LambdaResponseWithData(result), nil
}

func main() {
	lambda.Start(startDuaryAPI)
}
