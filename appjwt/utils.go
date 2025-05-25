package appjwt

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"net/http"
)

func GetLambdaClient() (*lambda.Client, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := lambda.NewFromConfig(cfg)
	return client, nil
}

func GetHttpClient() (*http.Client, error) {
	return &http.Client{}, nil

}

func GetDynamoDBClient() (*dynamodb.Client, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)

	return client, nil
}
