package auth

import (
	"github.com/cho8833/duary_lambda/internal/auth/jwtutil"
	"github.com/cho8833/duary_lambda/internal/member"
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

type CertResponse struct {
	Keys []jwtutil.JWK `json:"keys"`
}

type SignInRes struct {
	IsRegister bool                    `json:"isRegister"`
	Member     *member.Member          `json:"member"`
	Token      *jwtutil.ApplicationJWT `json:"token"`
}
