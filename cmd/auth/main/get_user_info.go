package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/internal/member"
	"github.com/cho8833/duary_lambda/internal/util"
	"strconv"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap cmd/auth/main/get_user_info.go && chmod 755 bootstrap && zip  build/package/get_user_info.zip bootstrap && rm bootstrap
*/
func getUserInfo(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	cacheClient := util.GetCacheClient()
	dynamodbClient, _ := cacheClient.GetDynamoDBClient()
	memberRepo := member.NewMemberRepository(dynamodbClient)

	memberSvc := member.NewMemberService(memberRepo)

	lambdaMap := req.RequestContext.Authorizer["lambda"].(map[string]interface{})
	provider := lambdaMap["provider"].(string)
	socialId, _ := strconv.ParseInt(lambdaMap["socialId"].(string), 10, 64)

	nenberInfo, err := memberSvc.GetMember(socialId, provider)

	if err != nil {
		return util.LambdaAppErrorResponse(err), nil
	}

	return util.LambdaResponseWithData(nenberInfo), nil

}

func main() {
	lambda.Start(getUserInfo)
}
