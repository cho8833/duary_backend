package appjwt

import (
	"fmt"
	"log"
)

type OIDCService interface {
	GetPublicKey(url string, provider string, kid string) (*JWK, error)
}

type OIDCServiceImpl struct {
	repository OIDCPublicKeyRepository
}

func NewOIDCService(repository *OIDCPublicKeyRepository) *OIDCServiceImpl {
	return &OIDCServiceImpl{repository: *repository}
}

func (svc *OIDCServiceImpl) GetPublicKey(url string, provider string, kid string) (*JWK, error) {
	// retrieve public jwk from DB
	// FindPublicKeyInDB 에 대한 Error 처리는 따로 하지 않고 logging 만 함
	// -> 다시 받아오면 되기 때문
	certRes, err := svc.repository.FindPublicKeyInDB(provider)
	if err != nil {
		log.Printf("failed to find key in db. provider: %s, error: %s", provider, err.Error())
	}

	// kid 가 일치하는 Public Key 찾기
	jwk := findMatchingKey(kid, certRes.Keys)

	// kid 가 일치하는 Public Key 가 DB 에 없으면 다시 받아옴
	if jwk == nil {
		log.Printf("there's no matching key with kid: %s, provider: %s, url: %s. Update Public Key", kid, provider, url)
		newCert, err := svc.repository.GetPublicJWK(url)
		if err != nil {
			log.Printf("failed to read jwk from %s", url)
			return nil, err
		}

		// SaveJWK 에 대한 Error 처리는 따로 하지 않고 logging 만 함
		// -> 다시 받아오면 되기 때문
		err = svc.repository.SaveJWK(provider, newCert.Keys)
		if err != nil {
			log.Printf("failed to save jwk\njwk:%+v\nerror:%s", newCert, err.Error())
		}

		// 새로 받아온 key 들에서 kid 가 일치하는 Public Key 찾기
		jwk = findMatchingKey(kid, newCert.Keys)
	}
	// if there's no matching key in new cert, return err
	if jwk == nil {
		errorString := fmt.Sprintf("no matching key with kid: %s, provider: %s, url: %s", kid, provider, url)
		log.Printf(errorString)
		return nil, fmt.Errorf(errorString)

	}

	return jwk, nil
}

func findMatchingKey(kid string, jwks []JWK) *JWK {
	if jwks == nil {
		return nil
	}
	for _, val := range jwks {
		if kid == val.Kid {
			return &val
		}
	}
	return nil
}
