package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/appjwt"
	"github.com/cho8833/duary_lambda/shared"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/get_oidc_public_key/main/main.go && chmod 755 bootstrap && zip  build/package/auth/get_oidc_public_key.zip bootstrap && rm bootstrap
*/
func getOIDCPublicKeyAPI(ctx context.Context, request *appjwt.GetPublicKeyReq) (*shared.ServerResponse[any], error) {
	// check Req
	if request.Url == "" || request.Provider == "" || request.Kid == "" {
		return shared.ResponseFromError(fmt.Errorf("Bad Request"), 400), nil
	}

	// load client
	httpClient, err := appjwt.GetHttpClient()
	if err != nil {
		return shared.ResponseFromError(err, 500), nil
	}
	dynamodbClient, err := appjwt.GetDynamoDBClient()
	if err != nil {
		return shared.ResponseFromError(err, 500), nil
	}

	// init service
	var repo appjwt.OIDCPublicKeyRepository = appjwt.NewOIDCPublicKeyRepository(httpClient, dynamodbClient)
	svc := appjwt.NewOIDCService(&repo)

	res, err := svc.GetPublicKey(request.Url, request.Provider, request.Kid)
	if err != nil {
		return shared.ResponseFromError(err, 400), nil
	}
	return shared.ResponseWithData(res), nil
}

func main() {
	lambda.Start(getOIDCPublicKeyAPI)
}
