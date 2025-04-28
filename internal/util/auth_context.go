package util

import (
	"github.com/aws/aws-lambda-go/events"
	"strconv"
)

type AuthContext struct {
	CoupleId *string
	SocialId *int64
	Provider *string
}

var singleInstance *AuthContext

func NewAuthContext(request events.APIGatewayProxyRequest) *AuthContext {
	lambdaMap := request.RequestContext.Authorizer["lambda"].(map[string]interface{})

	singleInstance = &AuthContext{}

	coupleId := lambdaMap["coupleId"]
	if coupleId != nil {
		parsedCoupleId := coupleId.(string)
		singleInstance.CoupleId = &parsedCoupleId
	}

	socialId := lambdaMap["socialId"]
	if socialId != nil {
		parsedSocialId, _ := strconv.ParseInt(socialId.(string), 10, 64)
		singleInstance.SocialId = &parsedSocialId
	}

	provider := lambdaMap["provider"]
	if provider != nil {
		parsedProvider := provider.(string)
		singleInstance.Provider = &parsedProvider
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
