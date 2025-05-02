package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cho8833/duary_lambda/internal/util"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type AppleAuthHelper struct{}

type AppleValidateCodeReq struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
}

type AppleValidateCodeRes struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
}

const appleValidateCodeUrl = "https://appleid.apple.com/auth/oauth2/v2/token"

// ValidateCode followed by generate and validate tokens: https://developer.apple.com/documentation/accountorganizationaldatasharing/generate-and-validate-tokens
func (helper *AppleAuthHelper) ValidateCode(req *AppleValidateCodeReq) (*AppleValidateCodeRes, error) {
	httpClient, err := util.GetCacheClient().GetHttpClient()
	if err != nil {
		log.Printf("get http client error: %v", err)
		return nil, err
	}

	pbytes, _ := json.Marshal(req)

	httpReq, err := http.NewRequest(http.MethodPost, appleValidateCodeUrl, bytes.NewBuffer(pbytes))

	res, err := httpClient.Do(httpReq)
	if err != nil {
		log.Printf("failed to request validation. error: %s", err.Error())
		return nil, util.InternalServerError{}
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Printf("failed to request validation. error: %s", res.Status)
		return nil, util.InternalServerError{}
	}

	validationRes := &AppleValidateCodeRes{}
	if err := json.NewDecoder(res.Body).Decode(validationRes); err != nil {
		log.Printf("failed to decode validation. error: %s", err.Error())
		return nil, util.InternalServerError{}
	}

	return validationRes, nil
}

// followed by creating client secret guideline: https://developer.apple.com/documentation/accountorganizationaldatasharing/creating-a-client-secret
func (helper *AppleAuthHelper) createClientSecret(privateKey []byte) (*string, error) {
	iss := os.Getenv("apple_team_id")
	kid := os.Getenv("apple_private_key_identifier")
	iat := time.Now().UTC().Unix()
	exp := iat + 3600
	secret := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": iss,
		"iat": iat,
		"exp": exp,
		"aud": "https://appleid.apple.com",
		"sub": appId,
	}, func(token *jwt.Token) {
		token.Header["kid"] = kid
	})

	token, err := secret.SignedString(privateKey)
	if err != nil {
		log.Printf("failed to sign token: %v", err)
	}
	return &token, nil
}

func (helper *AppleAuthHelper) getPrivateKey() (*[]byte, error) {

	s3Client, err := util.GetCacheClient().GetS3Client()
	if err != nil {
		log.Printf("failed to get S3 client, error: %s", err.Error())
		return nil, err
	}

	fileName := os.Getenv("apple_sign_in_private_key")

	fetch, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("duary"),
		Key:    aws.String(fileName),
	})
	if err != nil {
		log.Printf("failed to find apple sign in private key in duary bucket: %s", fileName)
		return nil, err
	}
	defer fetch.Body.Close()

	result, err := io.ReadAll(fetch.Body)
	if err != nil {
		log.Printf("failed read s3 fetch body: %s", err)
		return nil, err
	}
	return &result, nil
}
