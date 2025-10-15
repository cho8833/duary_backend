package cert

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"net/http"
)

type OIDCPublicKeyRepository interface {
	FindPublicKeyInDB(provider string) (*CertResponse, error)
	GetPublicJWK(url string) (*CertResponse, error)
	SaveJWK(provider string, jwks []JWK) error
}

type OIDCPublicKey struct {
	provider *string `json:"provider"`
	Keys     []JWK   `json:"keys"`
}

type OIDCPublicKeyRepositoryImpl struct {
	httpClient     *http.Client
	dynamoDBClient *dynamodb.Client
	tableName      string
}

func NewOIDCPublicKeyRepository(httpClient *http.Client, dynamodbClient *dynamodb.Client) *OIDCPublicKeyRepositoryImpl {
	return &OIDCPublicKeyRepositoryImpl{httpClient: httpClient, dynamoDBClient: dynamodbClient, tableName: "Cert"}
}

func (repository *OIDCPublicKeyRepositoryImpl) FindPublicKeyInDB(provider string) (*CertResponse, error) {
	result, err := repository.dynamoDBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(repository.tableName),
		Key: map[string]types.AttributeValue{
			"provider": &types.AttributeValueMemberS{Value: provider},
		},
	})
	if err != nil {
		return nil, err
	}

	res := &OIDCPublicKey{}
	err = attributevalue.UnmarshalMap(result.Item, res)
	if err != nil {
		return nil, err
	}

	return &CertResponse{Keys: res.Keys}, nil
}

func (repository *OIDCPublicKeyRepositoryImpl) GetPublicJWK(url string) (*CertResponse, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := repository.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		errorMsg := fmt.Sprintf("failed to retrieve public key: %s", err.Error())
		log.Printf(errorMsg)
		return nil, fmt.Errorf(errorMsg)
	}

	certRes := &CertResponse{}
	if err := json.NewDecoder(res.Body).Decode(certRes); err != nil {
		return nil, err
	}
	log.Printf("got public jwk : %+v\n", certRes)
	return certRes, nil
}

func (repository *OIDCPublicKeyRepositoryImpl) SaveJWK(provider string, jwks []JWK) error {

	var keys []types.AttributeValue

	for _, value := range jwks {
		key, err := attributevalue.MarshalMap(value)
		if err != nil {
			return err
		}
		keys = append(keys, &types.AttributeValueMemberM{Value: key})
	}

	_, err := repository.dynamoDBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(repository.tableName),
		Item: map[string]types.AttributeValue{
			"provider": &types.AttributeValueMemberS{Value: provider},
			"keys": &types.AttributeValueMemberL{
				Value: keys,
			},
		},
	})

	if err != nil {
		log.Printf("failed to save jwk : %s", err.Error())
		return err
	}

	log.Printf("saved jwk : %+v\n", jwks)
	return nil
}
