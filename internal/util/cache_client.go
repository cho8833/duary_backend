package util

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"net/http"
	"sync"
)

var lock = &sync.Mutex{}

type CacheClient struct {
	httpClient     *http.Client
	dynamoDBClient *dynamodb.Client
	lambdaClient   *lambda.Client
	s3Client       *s3.Client
}

var httpClientInstance *CacheClient

func GetCacheClient() *CacheClient {
	if httpClientInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if httpClientInstance == nil {
			httpClientInstance = &CacheClient{}
		}
	}
	return httpClientInstance
}
func (cacheClient *CacheClient) GetHttpClient() (*http.Client, error) {
	if cacheClient.httpClient == nil {
		cacheClient.httpClient = &http.Client{}
	}
	return cacheClient.httpClient, nil
}

func (cacheClient *CacheClient) GetDynamoDBClient() (*dynamodb.Client, error) {
	if cacheClient.dynamoDBClient == nil {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}
		cacheClient.dynamoDBClient = dynamodb.NewFromConfig(cfg)
	}
	return cacheClient.dynamoDBClient, nil
}

func (cacheClient *CacheClient) GetLambdaClient() (*lambda.Client, error) {
	if cacheClient.lambdaClient == nil {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}
		cacheClient.lambdaClient = lambda.NewFromConfig(cfg)
	}
	return cacheClient.lambdaClient, nil
}

func (cacheClient *CacheClient) GetS3Client() (*s3.Client, error) {
	if cacheClient.s3Client == nil {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}
		cacheClient.s3Client = s3.NewFromConfig(cfg)
	}
	return cacheClient.s3Client, nil
}
