package test

import (
	"fmt"
	oidcService "github.com/cho8833/duary_lambda/internal/auth"
	"github.com/cho8833/duary_lambda/internal/shared"
	testUtil "github.com/cho8833/duary_lambda/internal/test"
	"testing"
)

func Test_updatePublicKey(t *testing.T) {
	cacheClient := shared.GetCacheClient()
	httpClient, err := cacheClient.GetHttpClient()
	if err != nil {
		t.Fatalf(err.Error())
	}
	dynamodbClient := testUtil.CreateLocalDynamoDBClient()

	var repo oidcService.OIDCPublicKeyRepository = oidcService.NewOIDCPublicKeyRepository(httpClient, dynamodbClient)

	svc := oidcService.NewOIDCService(&repo)

	certRes, err := svc.GetPublicKey("https://kauth.kakao.com/.well-known/jwks.json", "kakao", "")
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Print(certRes)
}
