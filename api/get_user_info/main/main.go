package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/cho8833/duary_lambda/appjwt"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/shared"
	"log"
	"os"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/get_user_info/main/main.go && chmod 755 bootstrap && zip  build/package/get_user_info.zip bootstrap && rm bootstrap
*/
func getUserInfo(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	dynamoDBClient := dynamodb.NewFromConfig(cfg)

	memberRepo := member.NewMemberRepository(dynamoDBClient)

	memberSvc := member.NewMemberService(memberRepo)

	lambdaMap := req.RequestContext.Authorizer["lambda"].(map[string]interface{})
	provider := lambdaMap["provider"].(string)
	socialId, _ := lambdaMap["socialId"].(string)

	memberInfo, svcError := memberSvc.GetMember(socialId, provider)
	if svcError != nil {
		return shared.LambdaAppErrorResponse(svcError), nil
	}

	// generate application Token
	jwtUtil := &appjwt.Impl{}

	memberId := jwtUtil.GenerateSubject(memberInfo.SocialId, memberInfo.Provider)
	key := os.Getenv("secretKey")
	newToken := jwtUtil.NewToken(memberId, memberInfo.CoupleId, key)

	return shared.LambdaResponseWithDataAndHeader(memberInfo, appjwt.ApplicationJWTToHeader(*newToken)), nil

}

func main() {
	lambda.Start(getUserInfo)
}
