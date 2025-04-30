package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/internal/couple"
	"github.com/cho8833/duary_lambda/internal/util"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap internal/couple/main/get_couple_api.go && chmod 755 bootstrap && zip  build/package/couple/get_couple_api.zip bootstrap && rm bootstrap
*/
func getCouple(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// jwt Token 에 coupleId 가 없으면 couple 이 아직 생성되지 않았다고 봄
	authContext := util.NewAuthContext(request)
	if authContext.CoupleId == nil {
		return util.LambdaAppErrorResponse(util.CoupleNotFound{}), nil
	}

	cacheClient := util.GetCacheClient()
	dynamodbClient, _ := cacheClient.GetDynamoDBClient()

	coupleRepo := couple.NewCoupleRepository(dynamodbClient)
	coupleSvc := couple.NewCoupleService(coupleRepo)

	data, svcErr := coupleSvc.FindById(*authContext.CoupleId)
	if svcErr != nil {
		return util.LambdaAppErrorResponse(svcErr), nil
	}

	return util.LambdaResponseWithData(data), nil
}

func main() {
	lambda.Start(getCouple)
}
