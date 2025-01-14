package jwtutil

import (
	"fmt"
	"github.com/cho8833/duary_lambda/internal/member"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

type ApplicationJWT struct {
	AccessToken          string `json:"accessToken"`
	ExpireTime           int64  `json:"expireTime"`
	RefreshToken         string `json:"refreshToken"`
	RefreshTokenExpireAt int64  `json:"refreshTokenExpireAt"`
}

type JWTUtil interface {
	NewToken(id string, key string) *ApplicationJWT
	ValidateApplicationJWT(tokenString string, key string) (*string, error)
	GenerateSubject(member *member.Member) string
}

type Impl struct {
}

func (util *Impl) NewToken(id string, key string) *ApplicationJWT {
	secretKey := []byte(key)
	// accessToken: 1일 후 expire
	expireTime := time.Now().AddDate(0, 0, 1).Unix()
	// refreshToken: 7일 후 expire
	refreshTokenExpireAt := time.Now().AddDate(0, 0, 7).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": expireTime,
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": refreshTokenExpireAt,
	})

	accessTokenString, _ := accessToken.SignedString(secretKey)
	refreshTokenString, _ := refreshToken.SignedString(secretKey)

	result := &ApplicationJWT{
		AccessToken:          accessTokenString,
		RefreshToken:         refreshTokenString,
		ExpireTime:           expireTime,
		RefreshTokenExpireAt: refreshTokenExpireAt,
	}
	return result
}

func (util *Impl) ValidateApplicationJWT(tokenString string, key string) (*string, error) {
	secretKey := []byte(key)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		log.Printf("%s\ntoken:%s", err.Error(), tokenString)
		return nil, err
	}

	if !token.Valid {
		log.Printf("invalid token string: %s \nsecretKey: %s", tokenString, secretKey)
		return nil, fmt.Errorf("토큰이 유효하지 않습니다")
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	return &id, nil
}

func (util *Impl) GenerateSubject(member *member.Member) string {
	return fmt.Sprintf("%d-%s", member.SocialId, member.Provider)
}
