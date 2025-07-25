package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/shared"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/google_sign_in/main/main.go && chmod 755 bootstrap && zip  build/package/google_sign_in.zip bootstrap && rm bootstrap
*/

func googleSignIn(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamoDBClient, err := model.GetDynamoDBClient()

	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	memberRepository := member.NewRepository(dynamoDBClient)
	coupleRepository := couple.NewRepository(dynamoDBClient)

	memberSvc := member.NewService(memberRepository)
	coupleSvc := couple.NewService(coupleRepository)
}
