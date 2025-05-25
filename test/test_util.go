package test

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/cho8833/duary_lambda/shared"
)

func CreateLocalDynamoDBClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("fakeRegion"),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service string, region string, option ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: "http://localhost:8000"}, nil
				})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "fakeAccessKey", SecretAccessKey: "fakeSecretAccessKey", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB"}}),
	)
	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(cfg)
}

func CreateTestAuthContext(socialId *int64, provider *string, coupleId *string) *shared.AuthContext {
	authContext := shared.GetAuthContext()
	authContext.CoupleId = coupleId
	authContext.Provider = provider
	authContext.SocialId = socialId

	return authContext
}
