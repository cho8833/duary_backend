package test

import (
	"fmt"
	"github.com/cho8833/duary_lambda/appjwt"
	"github.com/cho8833/duary_lambda/auth"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/member"
	"testing"
)

func Test_KakaoSignIn(t *testing.T) {
	idToken := "eyJraWQiOiI5ZjI1MmRhZGQ1ZjIzM2Y5M2QyZmE1MjhkMTJmZWEiLCJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiJjMzkzNThiYmQxNmYwNDQ0MjA4ZWI2NThjMzdjZDY5ZSIsInN1YiI6IjM0Mjg4MzU4MDkiLCJhdXRoX3RpbWUiOjE3MzU4MTczMTUsImlzcyI6Imh0dHBzOi8va2F1dGgua2FrYW8uY29tIiwibmlja25hbWUiOiLsobDtmITruYgiLCJleHAiOjE3MzU4NjA1MTUsImlhdCI6MTczNTgxNzMxNSwibm9uY2UiOiJnT2xTclJhOWwyeG5wa2VHRnVLSFZzNnlNV2ZvdDZlT0RJREtyTEdDM2ZNQ1VWWkRiVyIsImVtYWlsIjoidzg4MzNAbmF2ZXIuY29tIn0.Mh9gJP1LPjIAf1RsIz7MT0dMsZIJc6So4xbQkDwwhUZeyqItZkKtlqYErZfIKKK_NXBecqknPsp4o76CdTUF98Invo6TDdg-_hR2f4MxwBBsHczjNupnYKqB6tSY115ipJGL0xXozk6K-UAI9icIQ8eiG76dXAiBpW4OxjR4O004B_me84kfgDKILlFsMU8N2E-5dpiIsDF-1LkVKUrVJbaapoJ6WI6y2YOohVPXpEe1AKsY1TMZYf9Sh9FBRITMtc7o9baBxRxTj1EO26e_I0x2VYyDN8XN46maQZtYKy2tciXJ1iQcB9ml--EVQzzrygHjsOQyOLxAeUEQwUoOTw"
	accessToken := "l3uqVtLQCvnpC6P5wiqHbqsFA7gQUM6CAAAAAQo8IpsAAAGUJsZDKCHmgQBvj-MV"

	reqToken := &auth.KakaoOAuthToken{
		IdToken:     &idToken,
		AccessToken: &accessToken,
	}

	dynamodbClient := CreateLocalDynamoDBClient()
	memberRepository := member.NewRepository(dynamodbClient)
	coupleRepository := couple.NewRepository(dynamodbClient)

	memberSvc := member.NewService(memberRepository)
	coupleSvc := couple.NewService(coupleRepository)

	svc := auth.NewAuthService(&appjwt.JWTValidatorImpl{}, &appjwt.Impl{}, memberSvc, coupleSvc)

	result, svcError := svc.KakaoSignIn(reqToken, ptr("asdf"))

	if svcError != nil {
		t.Fatalf(svcError.Error())
	}
	fmt.Print(result)
}
