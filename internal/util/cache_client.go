package util

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"log"
	"net/http"
	"sync"
)

var lock = &sync.Mutex{}

type CacheClient struct {
	httpClient      *http.Client
	dynamoDBClient  *dynamodb.Client
	lambdaClient    *lambda.Client
	s3Client        *s3.Client
	schedulerClient *scheduler.Client
	fcmClient       *messaging.Client
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

func (cacheClient *CacheClient) GetSchedulerClient() (*scheduler.Client, error) {
	if cacheClient.schedulerClient == nil {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}
		cacheClient.schedulerClient = scheduler.NewFromConfig(cfg)
	}
	return cacheClient.schedulerClient, nil
}

func (cacheClient *CacheClient) GetFCMClient() (*messaging.Client, error) {
	if cacheClient.fcmClient == nil {

		ctx := context.Background()

		app, err := firebase.NewApp(ctx, nil)
		if err != nil {
			log.Fatalf("Failed to create Firebase app: %v", err)
			return nil, err
		}
		log.Printf("Using Firebase Cloud Messaging client: %+v", app)
		client, err := app.Messaging(ctx)
		if err != nil {
			log.Printf("Error creating Messaging client: %+v", err)
			return nil, err
		}
		cacheClient.fcmClient = client
	}
	return cacheClient.fcmClient, nil
}
