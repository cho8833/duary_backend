package appjwt

import (
	"fmt"
	"github.com/aws/smithy-go/ptr"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"strconv"
	"time"
)

type JWTInfo struct {
	Sub      *string
	CoupleId *string
}
type ApplicationJWT struct {
	AccessToken          string `json:"accessToken"`
	ExpireTime           int64  `json:"expireTime"`
	RefreshToken         string `json:"refreshToken"`
	RefreshTokenExpireAt int64  `json:"refreshTokenExpireAt"`
}

// header 는 snake case
func ApplicationJWTToHeader(token ApplicationJWT) map[string]string {
	return map[string]string{
		"access_token":            token.AccessToken,
		"expire_time":             strconv.FormatInt(token.ExpireTime, 10),
		"refresh_token":           token.RefreshToken,
		"refresh_token_expire_at": strconv.FormatInt(token.RefreshTokenExpireAt, 10),
	}
}

type JWTUtil interface {
	NewToken(id string, coupleId *string, key string) *ApplicationJWT
	ValidateApplicationJWT(tokenString string, key string) (*JWTInfo, error)
	GenerateSubject(socialId string, provider string) string
}

type Impl struct {
}

func (util *Impl) NewToken(subject string, coupleId *string, key string) *ApplicationJWT {
	secretKey := []byte(key)
	// accessToken: 1일 후 expire
	expireTime := time.Now().AddDate(0, 0, 1).Unix()
	// refreshToken: 7일 후 expire
	refreshTokenExpireAt := time.Now().AddDate(0, 0, 7).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      subject,
		"coupleId": coupleId,
		"exp":      expireTime,
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": subject,
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

func (util *Impl) ValidateApplicationJWT(tokenString string, key string) (*JWTInfo, error) {
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
	coupleId := token.Claims.(jwt.MapClaims)["coupleId"]
	jwtInfo := JWTInfo{
		Sub: &id,
	}
	if coupleId != nil {
		jwtInfo.CoupleId = ptr.String(coupleId.(string))
	}
	return &jwtInfo, nil
}

func (util *Impl) GenerateSubject(socialId string, provider string) string {
	return fmt.Sprintf("%s-%s", socialId, provider)
}
