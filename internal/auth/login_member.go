package auth

import (
	"github.com/aws/aws-lambda-go/events"
	"strings"
)

type LoginMember struct {
	SocialId string
	Provider string
}

func FromSubject(subject *string) *LoginMember {
	s := strings.Split(*subject, "-")
	socialId := s[0]
	return &LoginMember{Provider: s[1], SocialId: socialId}
}

func FromRequestContext(req events.APIGatewayProxyRequest) *LoginMember {
	lambdaMap := req.RequestContext.Authorizer["lambda"].(map[string]interface{})
	provider := lambdaMap["provider"].(string)
	socialId, _ := lambdaMap["socialId"].(string)
	return &LoginMember{Provider: provider, SocialId: socialId}
}
