package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_backend/model"
	"github.com/cho8833/duary_backend/model/ws_connection"
	"github.com/cho8833/duary_backend/shared"
	"log"
)

func main() {
	lambda.Start(connect)
}

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap ws_api/connect/main.go && chmod 755 bootstrap && zip  build/package/ws_connect.zip bootstrap && rm bootstrap
*/
func connect(_ context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	authContext := shared.AuthContextFromWS(req)

	socialId := authContext.SocialId
	provider := authContext.Provider
	if socialId == nil || provider == nil {
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	dbClient, err := model.GetDynamoDBClient()
	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), err
	}
	stage := req.StageVariables["stage"]

	wsRepo := ws_connection.NewRepository(dbClient, stage)

	_, err = wsRepo.Create(*socialId, *provider, req.RequestContext.ConnectionID)
	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), err
	}

	return shared.LambdaResponseWithData(nil), nil
}
