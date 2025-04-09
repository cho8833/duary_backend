package event

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"time"
)

const tableName = "Event"

type Repository interface {
	FindByCoupleIdAndStartDateBefore(coupleId string, startDate time.Time) ([]VO, error)
	SaveEvent(event *Event) (*Event, error)
}

type RepositoryDynamoDB struct {
	client *dynamodb.Client
}

func NewEventRepository(client *dynamodb.Client) *RepositoryDynamoDB {
	return &RepositoryDynamoDB{client: client}
}

func (repo *RepositoryDynamoDB) SaveEvent(event *Event) (*Event, error) {

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

func (repo *RepositoryDynamoDB) FindByCoupleIdAndStartDateBefore(coupleId string, startDate time.Time) ([]VO, error) {
	rangeStart := startDate.Format(time.RFC3339) + "#"
	keyEx := expression.Key("coupleId").Equal(expression.Value(coupleId)).And(expression.Key("startDateTime").LessThan(expression.Value(rangeStart)))
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

	var result []VO
	for _, e := range events {
		result = append(result, FromEvent(e))
	}

	return result, nil
}

//func (repo *RepositoryDynamoDB) UpdateEvent(req *UpdateReq) (*Event, error) {
//	av, err := attributevalue.MarshalMap(req)
//	if err != nil {
//		log.Printf(err.Error())
//		return nil, err
//	}
//	builder := expression.UpdateBuilder{}
//
//	for k, v := range av {
//		builder = builder.Set(expression.Name(k), expression.Value(v))
//	}
//
//	expr, err := expression.NewBuilder().WithUpdate(builder).Build()
//
//	if err != nil {
//		log.Printf(err.Error())
//		return nil, err
//	}
//
//	response, err := repo.client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
//		TableName: aws.String(tableName),
//		Key:
//	})
//}
