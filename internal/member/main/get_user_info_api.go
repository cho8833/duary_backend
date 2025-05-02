package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/internal/member"
	"github.com/cho8833/duary_lambda/internal/util"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap internal/member/main/get_user_info_api.go && chmod 755 bootstrap && zip  build/package/member/get_user_info.zip bootstrap && rm bootstrap
*/
func getUserInfo(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	cacheClient := util.GetCacheClient()
	dynamodbClient, _ := cacheClient.GetDynamoDBClient()
	memberRepo := member.NewMemberRepository(dynamodbClient)

	memberSvc := member.NewMemberService(memberRepo)

	lambdaMap := req.RequestContext.Authorizer["lambda"].(map[string]interface{})
	provider := lambdaMap["provider"].(string)
	socialId, _ := lambdaMap["socialId"].(string)

	memberInfo, err := memberSvc.GetMember(socialId, provider)

	if err != nil {
		return util.LambdaAppErrorResponse(err), nil
	}

	return util.LambdaResponseWithData(memberInfo), nil

}

func main() {
	lambda.Start(getUserInfo)
}
