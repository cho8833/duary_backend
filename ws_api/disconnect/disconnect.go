package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/ws_connection"
	"github.com/cho8833/duary_lambda/shared"
	"log"
)

func main() {
	lambda.Start(disconnect)
}

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap ws_api/disconnect/disconnect.go && chmod 755 bootstrap && zip  build/package/ws_disconnect.zip bootstrap && rm bootstrap
*/
func disconnect(_ context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

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

	wsRepo := ws_connection.NewRepository(dbClient)

	err = wsRepo.DeleteBySocialIdAndProvider(*socialId, *provider)
	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), err
	}

	return shared.LambdaResponseWithData(nil), nil

}
