package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/internal/auth"
	"github.com/cho8833/duary_lambda/internal/auth/jwtutil"
	"log"
	"os"
	"strings"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap internal/auth/main/jwt_authorizer.go && chmod 755 bootstrap && zip  build/package/auth/jwt_authorizer.zip bootstrap && rm bootstrap
*/
func jwtAuthorizer(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayV2CustomAuthorizerIAMPolicyResponse, error) {
	key := os.Getenv("secretKey")
	jwtUtil := jwtutil.Impl{}
	header := request.Headers["authorization"]
	log.Printf("got authorization : %s", header)

	check := strings.Split(header, " ")
	token := check[1]
	jwtInfo, err := jwtUtil.ValidateApplicationJWT(token, key)
	if err != nil {
		log.Printf("Authorization: fail to authorize. token: %s, error: %s", token, err.Error())
		return denyResponse(), nil
	}

	loginMember := auth.FromSubject(jwtInfo.Sub)

	log.Printf("authorized %#v", *jwtInfo)

	lambdaContext := map[string]interface{}{
		"socialId": loginMember.SocialId,
		"provider": loginMember.Provider,
	}

	// values of context map seem not to be nil
	if jwtInfo.CoupleId != nil {
		lambdaContext["coupleId"] = *jwtInfo.CoupleId
	}

	return events.APIGatewayV2CustomAuthorizerIAMPolicyResponse{
		PrincipalID: loginMember.SocialId,
		Context:     lambdaContext,
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action: []string{"execute-api:Invoke"},
					Effect: "Allow",
					Resource: []string{
						"arn:aws:execute-api:ap-northeast-2:922001515124:i3lm91q9v5/*",
					},
				},
			},
		},
	}, nil
}

func denyResponse() events.APIGatewayV2CustomAuthorizerIAMPolicyResponse {
	return events.APIGatewayV2CustomAuthorizerIAMPolicyResponse{
		PrincipalID: "",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Deny",
					Resource: []string{"*"},
				},
			},
		},
	}
}

func main() {
	lambda.Start(jwtAuthorizer)
}
