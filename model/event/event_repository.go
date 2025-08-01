package event

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/shared"
	"log"
	"strings"
	"time"
)

const tableName = "Event"
const typeIndexName = "eventType-index"

type Repository interface {
	FindByCoupleIdAndStartDateBefore(coupleId string, startDate time.Time) ([]VO, error)
	Save(event *Event) (*VO, error)
	Update(coupleId string, id string, req EditReq) (*VO, error)
	Delete(coupleId string, id string) error
	SaveTransaction(event *Event, transaction *model.DynamoDBWriteTransaction) (*VO, error)
	DeleteTransaction(coupleId string, id string, transaction *model.DynamoDBWriteTransaction) error
	QueryByCoupleIdAndType(coupleId string, eventType EventType, include bool) ([]VO, error)
}

type RepositoryDynamoDB struct {
	client *dynamodb.Client
}

func NewRepository(client *dynamodb.Client) *RepositoryDynamoDB {
	return &RepositoryDynamoDB{client: client}
}

func (repo *RepositoryDynamoDB) Save(event *Event) (*VO, error) {

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

	vo := FromEvent(*event)

	return &vo, nil
}

func (repo *RepositoryDynamoDB) SaveTransaction(event *Event, transaction *model.DynamoDBWriteTransaction) (*VO, error) {
	item, err := attributevalue.MarshalMap(event)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	transactionItem := &types.TransactWriteItem{
		Put: &types.Put{
			TableName: aws.String(tableName),
			Item:      item,
		}}
	transaction.AddTransaction(transactionItem)

	vo := FromEvent(*event)

	return &vo, nil
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

func (repo *RepositoryDynamoDB) findById(coupleId string, id string) (*Event, error) {
	output, err := repo.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       repo.getKey(coupleId, id),
	})
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	if output.Item == nil {
		return nil, &types.ResourceNotFoundException{
			Message: aws.String(fmt.Sprintf("resource not found for id: %s", id)),
		}
	}

	var item Event
	err = attributevalue.UnmarshalMap(output.Item, &item)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	return &item, nil
}

func (repo *RepositoryDynamoDB) Update(coupleId string, id string, req EditReq) (*VO, error) {

	foundEvent, err := repo.findById(coupleId, id)
	if err != nil {
		return nil, err
	}

	sortKeySplit := strings.Split(foundEvent.StartDateTime, "#")
	originalStartDateTime, err := time.Parse(time.RFC3339, sortKeySplit[0])
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	didUpdateStartDateTime := !originalStartDateTime.Equal(*req.StartDateTime)

	// StartDateTime 이 업데이트 되는 경우 sort key 가 바뀌기 때문에 기존 item 을 삭제하고 새로 생성해야 함
	if didUpdateStartDateTime {
		transaction := model.NewWriteTransaction(repo.client)
		transaction.BeginTransaction()

		err := repo.DeleteTransaction(coupleId, id, transaction)
		if err != nil {
			log.Printf(err.Error())
			return nil, err
		}

		saveEvent := foundEvent.copyWith(req)
		updated, err := repo.SaveTransaction(&saveEvent, transaction)
		if err != nil {
			log.Printf(err.Error())
			return nil, err
		}

		_, err = transaction.Execute()
		if err != nil {
			log.Printf(err.Error())
			return nil, err
		}

		return updated, nil

		// StartDateTime 이 바뀌지 않은 경우 Update
	} else {
		av, err := attributevalue.MarshalMap(req)
		if err != nil {
			log.Printf(err.Error())
			return nil, err
		}

		builder := expression.UpdateBuilder{}

		for k, v := range av {
			builder = builder.Set(expression.Name(k), expression.Value(v))
		}

		expr, err := expression.NewBuilder().WithUpdate(builder).Build()
		if err != nil {
			log.Printf(err.Error())
			return nil, err
		}

		response, err := repo.client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
			TableName:                 aws.String(tableName),
			Key:                       repo.getKey(coupleId, id),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueAllNew,
		})

		if err != nil {
			log.Printf("failed to update item. Req: %+v, error: %s", req, err.Error())
			return nil, err
		}

		updated := &Event{}
		_ = attributevalue.UnmarshalMap(response.Attributes, updated)

		result := FromEvent(*updated)

		return &result, nil
	}
}

func (repo *RepositoryDynamoDB) Delete(coupleId string, id string) error {

	deleteInput := &dynamodb.DeleteItemInput{
		TableName:    aws.String(tableName),
		Key:          repo.getKey(coupleId, id),
		ReturnValues: types.ReturnValueAllOld,
	}
	response, err := repo.client.DeleteItem(context.TODO(), deleteInput)

	if err != nil {
		log.Printf("failed to delete item. coupleId: %s, id: %s, error: %s", coupleId, id, err.Error())
		return err
	}
	if response.Attributes == nil {
		log.Printf("failed to delete item. item doesn't exist. coupleId: %s, id: %s", coupleId, id)
		return shared.DBDeleteError{}
	}
	log.Printf("deleted item. coupleId: %s, id: %s, response: %+v", coupleId, id, response)

	return nil
}

func (repo *RepositoryDynamoDB) DeleteTransaction(coupleId string, id string, transaction *model.DynamoDBWriteTransaction) error {
	transactionItem := &types.TransactWriteItem{
		Delete: &types.Delete{
			TableName: aws.String(tableName),
			Key:       repo.getKey(coupleId, id),
		},
	}
	transaction.AddTransaction(transactionItem)

	return nil
}

func (repo *RepositoryDynamoDB) QueryByCoupleIdAndType(coupleId string, eventType EventType, include bool) ([]VO, error) {
	keyEx := expression.Key("coupleId").Equal(expression.Value(coupleId)).And(expression.Key("eventType").Equal(expression.Value(eventType)))
	var expr expression.Expression
	var err error
	if include {
		proj := expression.NamesList(
			expression.Name("createdBy"),
			expression.Name("coupleId"),
			expression.Name("startDateTime"))
		expr, err = expression.NewBuilder().WithKeyCondition(keyEx).WithProjection(proj).Build()
	} else {
		expr, err = expression.NewBuilder().WithKeyCondition(keyEx).Build()
	}
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String(typeIndexName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
	}
	output, err := repo.client.Query(context.TODO(), input)
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

func (repo *RepositoryDynamoDB) getKey(coupleId string, id string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"coupleId":      &types.AttributeValueMemberS{Value: coupleId},
		"startDateTime": &types.AttributeValueMemberS{Value: id},
	}
}
