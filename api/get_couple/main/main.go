package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/shared"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/get_couple/main/main.go && chmod 755 bootstrap && zip  build/package/get_couple_api.zip bootstrap && rm bootstrap
*/
func getCouple(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// jwt Token 에 coupleId 가 없으면 couple 이 아직 생성되지 않았다고 봄
	authContext := shared.NewAuthContext(request)
	if authContext.CoupleId == nil {
		return shared.LambdaAppErrorResponse(shared.CoupleNotFound{}), nil
	}

	dynamodbClient, err := model.GetDynamoDBClient()
	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	coupleRepo := couple.NewRepository(dynamodbClient)
	coupleSvc := couple.NewService(coupleRepo)

	data, svcErr := coupleSvc.FindById(*authContext.CoupleId)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	return shared.LambdaResponseWithData(data), nil
}

func main() {
	lambda.Start(getCouple)
}
