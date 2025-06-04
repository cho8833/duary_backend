package event_schedule

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cho8833/duary_lambda/model"
	"log"
)

const tableName = "EventSchedule"

type Repository interface {
	Save(data *EventSchedule) (*EventSchedule, error)
	UpdateTransaction(memberId string, eventId string, req UpdateEventScheduleReq, transaction *model.DynamoDBWriteTransaction) error
}

type RepositoryDynamoDB struct {
	client *dynamodb.Client
}

func NewEventScheduleRepository(client *dynamodb.Client) *RepositoryDynamoDB {
	return &RepositoryDynamoDB{client: client}
}

func (repo *RepositoryDynamoDB) Save(data *EventSchedule) (*EventSchedule, error) {
	item, err := attributevalue.MarshalMap(data)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	res, err := repo.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName:    aws.String(tableName),
		Item:         item,
		ReturnValues: types.ReturnValueAllNew,
	})
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	result := &EventSchedule{}
	err = attributevalue.UnmarshalMap(res.Attributes, result)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	return result, nil
}

func (repo *RepositoryDynamoDB) UpdateTransaction(memberId string, eventId string, req UpdateEventScheduleReq, transaction *model.DynamoDBWriteTransaction) error {
	av, err := attributevalue.MarshalMap(req)
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	builder := expression.UpdateBuilder{}

	for k, v := range av {
		builder = builder.Set(expression.Name(k), expression.Value(v))
	}

	expr, err := expression.NewBuilder().WithUpdate(builder).Build()
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	transactionItem := &types.TransactWriteItem{
		Update: &types.Update{
			TableName:                 aws.String(tableName),
			Key:                       repo.getKey(memberId, eventId),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
		},
	}
	transaction.AddTransaction(transactionItem)

	return nil
}

func (repo *RepositoryDynamoDB) getKey(memberId string, eventId string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"memberId": &types.AttributeValueMemberS{Value: memberId},
		"eventId":  &types.AttributeValueMemberS{Value: eventId},
	}
}
