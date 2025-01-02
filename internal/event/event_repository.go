package event

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"time"
)

const tableName = "Event"

type Repository interface {
	FindByCoupleIdAndEndTime(coupleId *string, date *time.Time) ([]Event, error)
}

type RepositoryDynamoDB struct {
	client *dynamodb.Client
}

func NewEventRepository(client *dynamodb.Client) *RepositoryDynamoDB {
	return &RepositoryDynamoDB{client: client}
}

func (repo *RepositoryDynamoDB) SaveEvent(event *Event) (*Event, error) {
	
}

func (repo *RepositoryDynamoDB) FindByCoupleIdAndDate(coupleId *string, date *time.Time) ([]Event, error) {
	keyEx := expression.Key("coupleId").Equal(expression.Value(*coupleId))
	formattedDate := fmt.Sprintf("%d-%02d-%d", date.Year(), date.Month(), date.Day())
	filterEx := expression.Name("endDateTime").BeginsWith(formattedDate)
	expr, err := expression.NewBuilder().WithFilter(filterEx).WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, err
	}
	output, err := repo.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		FilterExpression:          expr.Filter(),
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

func (repo *RepositoryDynamoDB) getKey(event *Event) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"coupleId":      &types.AttributeValueMemberS{Value: *event.CoupleId},
		"startDateTime": &types.AttributeValueMemberS{Value: event.EndTime.String()},
	}
}
