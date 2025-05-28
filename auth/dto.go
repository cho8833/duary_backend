package auth

import (
	"github.com/cho8833/duary_lambda/appjwt"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/member"
	"time"
)

type KakaoOAuthToken struct {
	AccessToken           *string `json:"access_token"`
	expiresAt             *time.Time
	refreshToken          *string
	refreshTokenExpiresAt *time.Time
	scopes                *[]string
	IdToken               *string `json:"id_token"`
}

type AppleOAuthToken struct {
	UserIdentifier    *string `json:"userIdentifier"`
	GivenName         *string `json:"givenName"`
	FamilyName        *string `json:"familyName"`
	Email             *string `json:"email"`
	AuthorizationCode string  `json:"authorizationCode"`
	IdentityToken     *string `json:"identityToken"`
	State             *string `json:"state"`
}

type SignInReq struct {
	FcmToken        *string          `json:"fcmToken"`
	AppleOAuthToken *AppleOAuthToken `json:"appleOAuthToken"`
	KakaoOAuthToken *KakaoOAuthToken `json:"kakaoOAuthToken"`
}

type SignInRes struct {
	IsRegister bool                   `json:"isRegister"`
	Member     *member.Member         `json:"member"`
	Token      *appjwt.ApplicationJWT `json:"token"`
	Couple     *couple.Couple         `json:"couple"`
}
