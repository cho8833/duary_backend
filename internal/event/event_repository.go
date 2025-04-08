package event

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	uuid2 "github.com/google/uuid"
	"log"
	"time"
)

const tableName = "Event"

type Repository interface {
	FindByCoupleIdAndStartDateBefore(coupleId string, startDate time.Time) ([]Event, error)
	SaveEvent(event *Event) (*Event, error)
}

type RepositoryDynamoDB struct {
	client *dynamodb.Client
}

func NewEventRepository(client *dynamodb.Client) *RepositoryDynamoDB {
	return &RepositoryDynamoDB{client: client}
}

func (repo *RepositoryDynamoDB) SaveEvent(event *Event) (*Event, error) {
	if event.Id == nil {
		event.Id = repo.generateUID()
	}
	item, err := attributevalue.MarshalMap(event)
	if err != nil {
		return nil, err
	}

	_, err = repo.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Event"),
		Item:      item,
	})

	if err != nil {
		return nil, err
	}

	return event, nil
}

func (repo *RepositoryDynamoDB) FindByCoupleIdAndStartDateBefore(coupleId string, startDate time.Time) ([]Event, error) {
	keyEx := expression.Key("coupleId").Equal(expression.Value(coupleId)).And(expression.Key("startDate").LessThan(expression.Value(startDate)))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, err
	}
	output, err := repo.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
	})
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	var events []Event
	if err := attributevalue.UnmarshalListOfMaps(output.Items, &events); err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	return events, nil
}

func (repo *RepositoryDynamoDB) generateUID() *string {
	uuid := uuid2.New().String()
	return &uuid
}
