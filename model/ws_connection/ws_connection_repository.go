package ws_connection

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
)

type Repository interface {
	DeleteBySocialIdAndProvider(socialId string, provider string) error
	FindBySocialIdAndProvider(socialId string, provider string) (*WSConnection, error)
	Create(socialId string, provider string, connectionId string) (*WSConnection, error)
}

type RepositoryDynamoDB struct {
	client    dynamodb.Client
	tableName string
}

func NewRepository(client *dynamodb.Client, stage string) RepositoryDynamoDB {
	var table string
	if stage == "dev" {
		table = "dev_WSConnection"
	} else {
		table = "WSConnection"
	}
	return RepositoryDynamoDB{client: *client, tableName: table}
}

func (repo *RepositoryDynamoDB) FindBySocialIdAndProvider(socialId string, provider string) (*WSConnection, error) {
	result, err := repo.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key:       repo.getKey(socialId, provider),
	})
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, &types.ResourceNotFoundException{
			Message: aws.String(fmt.Sprintf("resource not found for id: %s", socialId)),
		}
	}

	wsConnection := &WSConnection{}
	err = attributevalue.UnmarshalMap(result.Item, wsConnection)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return wsConnection, nil
}

func (repo *RepositoryDynamoDB) DeleteBySocialIdAndProvider(socialId string, provider string) error {
	input := &dynamodb.DeleteItemInput{
		TableName:    aws.String(repo.tableName),
		Key:          repo.getKey(socialId, provider),
		ReturnValues: types.ReturnValueNone,
	}

	_, err := repo.client.DeleteItem(context.TODO(), input)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (repo *RepositoryDynamoDB) Create(socialId string, provider string, connectionId string) (*WSConnection, error) {
	data := &WSConnection{
		SocialId:     socialId,
		Provider:     provider,
		ConnectionId: connectionId,
	}
	item, err := attributevalue.MarshalMap(data)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(repo.tableName),
		Item:      item,
	}

	_, err = repo.client.PutItem(context.TODO(), input)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return data, nil
}

func (repo *RepositoryDynamoDB) getKey(socialId string, provider string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"socialId": &types.AttributeValueMemberS{Value: socialId},
		"provider": &types.AttributeValueMemberS{Value: provider},
	}
}
