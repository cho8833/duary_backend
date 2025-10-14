package appjwt

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/cho8833/duary_backend/model/cert"
	"github.com/cho8833/duary_backend/shared"
	"github.com/golang-jwt/jwt/v5"
	"math/big"
	"time"
)

type GetPublicKeyReq struct {
	Url      string `json:"url"`
	Provider string `json:"provider"`
	Kid      string `json:"kid"`
}

type ValidatingValue struct {
	Iss      string
	Aud      []string
	Nonce    string
	Url      string
	Provider string
}

type DecodedPayload struct {
	SocialId string
	Exp      time.Time
	Email    *string
	NickName *string
	Name     *string
}

type JWTValidator interface {
	VerifyRSA256(idToken string, value *ValidatingValue) (*DecodedPayload, error)
}

type JWTValidatorImpl struct {
}

func (validator *JWTValidatorImpl) VerifyRSA256(idToken string, value *ValidatingValue) (*DecodedPayload, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(idToken, claims, func(token *jwt.Token) (interface{}, error) {
		// verify with public key
		kid := token.Header["kid"].(string)
		jwk, err := getPublicKey(kid, value.Url, value.Provider)
		if err != nil {
			return nil, err
		}

		n, err := decode(jwk.N)
		if err != nil {
			return nil, err
		}
		e, err := decode(jwk.E)
		if err != nil {
			return nil, err
		}
		pk := &rsa.PublicKey{
			N: new(big.Int).SetBytes(n),
			E: int(new(big.Int).SetBytes(e).Int64()),
		}
		return pk, nil
	})
	if err != nil {
		return nil, err
	}

	var expireTime time.Time
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// check issuer
		iss, err := claims.GetIssuer()
		if err != nil {
			return nil, err
		}
		if value.Iss != "" {
			if iss != value.Iss {
				return nil, fmt.Errorf("issuer does not match")
			}
		}

		// check audience
		aud, err := claims.GetAudience()
		if err != nil {
			return nil, err
		}
		isAudValid := false
		if len(value.Aud) != 0 {
			for _, a := range value.Aud {
				if aud[0] == a {
					isAudValid = true
				}
			}
		} else {
			isAudValid = true
		}
		if !isAudValid {
			return nil, fmt.Errorf("audience does not match")
		}

		// check expirationTime, must be before now
		expireTime, err := claims.GetExpirationTime()
		if err != nil {
			return nil, err
		}
		if expireTime.Before(time.Now()) {
			return nil, fmt.Errorf("expire time is not valid")
		}

		// check Nonce, pair with frontend
		nonce := claims["nonce"].(string)
		if nonce != value.Nonce {
			return nil, fmt.Errorf("nonce does not match")
		}

	}

	socialId := claims["sub"].(string)

	email := toStringPointer(claims["email"])

	nickName := toStringPointer(claims["nickname"])

	result := &DecodedPayload{
		SocialId: socialId,
		Exp:      expireTime,
		Email:    email,
		NickName: nickName,
	}
	return result, nil
}

func decode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

func getPublicKey(kid string, url string, provider string) (*cert.JWK, error) {
	payload, err := json.Marshal(&GetPublicKeyReq{
		Url:      url,
		Provider: provider,
		Kid:      kid,
	})
	if err != nil {
		return nil, err
	}
	client, err := GetLambdaClient()
	if err != nil {
		return nil, err
	}

	invokeInput := lambda.InvokeInput{
		FunctionName: aws.String("get_oidc_public_key"),
		LogType:      types.LogTypeTail,
		Payload:      payload,
	}
	invokeOutput, err := client.Invoke(context.TODO(), &invokeInput)
	if err != nil {
		return nil, err
	}
	if invokeOutput.StatusCode != 200 {
		return nil, fmt.Errorf("%+v", invokeOutput.Payload)
	}
	jwkResponse := &shared.ServerResponse[cert.JWK]{}
	err = json.Unmarshal(invokeOutput.Payload, jwkResponse)
	if err != nil {
		return nil, err
	}
	if jwkResponse.Status != 200 {
		return nil, fmt.Errorf(*jwkResponse.Message)
	}

	return &jwkResponse.Data, nil
}

func toStringPointer(val interface{}) *string {
	if val == nil {
		return nil
	}

	// Check for string
	if str, ok := val.(string); ok {
		return &str
	}

	// Check for *string
	if strPtr, ok := val.(*string); ok {
		return strPtr
	}

	// If it's something else (like a []byte or fmt.Stringer), handle as needed
	// Example: convert using fmt.Sprintf or reflect if required

	return nil // or panic/log if you want to catch type mismatches
}
