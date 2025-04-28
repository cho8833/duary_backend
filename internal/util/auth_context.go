package util

import "github.com/aws/aws-lambda-go/events"

type AuthContext struct {
	CoupleId *string
	SocialId *int64
	Provider *string
}

var singleInstance *AuthContext

func NewAuthContext(request events.APIGatewayProxyRequest) *AuthContext {
	lambdaMap := request.RequestContext.Authorizer["lambda"].(map[string]interface{})

	singleInstance = &AuthContext{
		CoupleId: lambdaMap["coupleId"].(*string),
		SocialId: lambdaMap["socialId"].(*int64),
		Provider: lambdaMap["provider"].(*string),
	}

	return singleInstance
}

func GetAuthContext() *AuthContext {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			singleInstance = &AuthContext{}
		}
	}
	return singleInstance
}
