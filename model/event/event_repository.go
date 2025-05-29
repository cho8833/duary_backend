package event

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go/ptr"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/shared"
	"log"
	"strings"
	"time"
)

const tableName = "Event"

type Repository interface {
	FindByCoupleIdAndStartDateBefore(coupleId string, startDate time.Time) ([]VO, error)
	SaveEvent(event *Event) (*VO, error)
	UpdateEvent(req *UpdateReq) (*VO, error)
	DeleteEvent(coupleId string, id string) error
	SaveTransaction(event *Event, transaction *model.DynamoDBWriteTransaction) (*VO, error)
}

type RepositoryDynamoDB struct {
	client *dynamodb.Client
}

func NewEventRepository(client *dynamodb.Client) *RepositoryDynamoDB {
	return &RepositoryDynamoDB{client: client}
}

func (repo *RepositoryDynamoDB) SaveEvent(event *Event) (*VO, error) {

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

func (repo *RepositoryDynamoDB) UpdateEvent(req *UpdateReq) (*VO, error) {
	// startDateTime 이 바뀔 경우 ID 가 바뀜
	// update 할 event 식별을 위한 sortKey 저장
	sortKey := req.Id

	// MarshalMap 호출 시 partitionKey 를 제외하기 위해 req.CoupleId 를 nil 로 설정
	// partitionKey 임시 저장
	partitionKey := req.CoupleId
	req.CoupleId = nil

	// StartDateTime 이 바뀐다면 sortKey 를 업데이트 해줘야 함
	if req.StartDateTime != nil {
		id := strings.Split(*req.Id, "#")[1]
		req.Id = ptr.String(req.StartDateTime.Format(time.RFC3339) + "#" + id)
	}

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
		Key:                       repo.getKey(*partitionKey, *sortKey),
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

func (repo *RepositoryDynamoDB) DeleteEvent(coupleId string, id string) error {

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

func (repo *RepositoryDynamoDB) getKey(coupleId string, id string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"coupleId":      &types.AttributeValueMemberS{Value: coupleId},
		"startDateTime": &types.AttributeValueMemberS{Value: id},
	}
}
