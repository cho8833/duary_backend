package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

/*
GCO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/update_couple/main/main.go && chmod 755 bootstrap && zip build/package/common/update_couple.zip bootstrap && rm bootstrap
*/

func updateCouple(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

}

func main() {
	lambda.Start(updateCouple)
}
