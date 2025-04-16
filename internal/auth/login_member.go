package auth

import (
	"github.com/aws/aws-lambda-go/events"
	"strconv"
	"strings"
)

type LoginMember struct {
	SocialId int64
	Provider string
}

func FromSubject(subject *string) *LoginMember {
	s := strings.Split(*subject, "-")
	socialId, _ := strconv.ParseInt(s[0], 10, 64)
	return &LoginMember{Provider: s[1], SocialId: socialId}
}

func FromRequestContext(req events.APIGatewayProxyRequest) *LoginMember {
	lambdaMap := req.RequestContext.Authorizer["lambda"].(map[string]interface{})
	provider := lambdaMap["provider"].(string)
	socialId, _ := strconv.ParseInt(lambdaMap["socialId"].(string), 10, 64)
	return &LoginMember{Provider: provider, SocialId: socialId}
}
