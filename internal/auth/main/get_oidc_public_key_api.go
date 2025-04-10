package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	authService "github.com/cho8833/duary_lambda/internal/auth"
	"github.com/cho8833/duary_lambda/internal/auth/jwtutil"
	"github.com/cho8833/duary_lambda/internal/util"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap internal/auth/main/get_oidc_public_key.go && chmod 755 bootstrap && zip  build/package/auth/get_oidc_public_key.zip bootstrap && rm bootstrap
*/
func getOIDCPublicKeyAPI(ctx context.Context, request *jwtutil.GetPublicKeyReq) (*util.ServerResponse[any], error) {
	// check Req
	if request.Url == "" || request.Provider == "" || request.Kid == "" {
		return util.ResponseFromError(fmt.Errorf("Bad Request"), 400), nil
	}

	// load client
	cacheClient := util.GetCacheClient()
	httpClient, err := cacheClient.GetHttpClient()
	if err != nil {
		return util.ResponseFromError(err, 500), nil
	}
	dynamodbClient, err := cacheClient.GetDynamoDBClient()
	if err != nil {
		return util.ResponseFromError(err, 500), nil
	}

	// init service
	var repo authService.OIDCPublicKeyRepository = authService.NewOIDCPublicKeyRepository(httpClient, dynamodbClient)
	svc := authService.NewOIDCService(&repo)

	res, err := svc.GetPublicKey(request.Url, request.Provider, request.Kid)
	if err != nil {
		return util.ResponseFromError(err, 400), nil
	}
	return util.ResponseWithData(res), nil
}

func main() {
	lambda.Start(getOIDCPublicKeyAPI)
}
