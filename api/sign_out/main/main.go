package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/shared"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/sign_out/main/main.go && chmod 755 bootstrap && zip  build/package/sign_out.zip bootstrap && rm bootstrap
*/
func signOut(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	dynamodbClient, err := model.GetDynamoDBClient()
	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	memberRepo := member.NewMemberRepository(dynamodbClient)
	memberSvc := member.NewMemberService(memberRepo)

	authContext := shared.NewAuthContext(request)

	// member 의 fcm token 제거
	_, svcErr := memberSvc.UpdateFcmToken(*authContext.SocialId, *authContext.Provider, nil)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	return shared.LambdaResponseWithData(nil), nil
}

func main() {
	lambda.Start(signOut)
}
