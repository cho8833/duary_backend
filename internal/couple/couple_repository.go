package couple

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
)

const tableName = "Couple"

type Repository interface {
	SaveCouple(couple *Couple) (*Couple, error)
	SaveCoupleTransaction(couple *Couple) (*types.TransactWriteItem, error)
	FindById(id *string) (*Couple, error)
	FindByCoupleCode(coupleCode *string) ([]Couple, error)
	UpdateCoupleTransaction(req *UpdateCoupleReq) (*types.TransactWriteItem, error)
}

type RepositoryDynamoDB struct {
	client *dynamodb.Client
}

func NewCoupleRepository(client *dynamodb.Client) *RepositoryDynamoDB {
	return &RepositoryDynamoDB{client: client}
}

func (repository *RepositoryDynamoDB) SaveCouple(couple *Couple) (*Couple, error) {
	item, err := attributevalue.MarshalMap(couple)
	if err != nil {
		return nil, err
	}

	_, err = repository.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Couple"),
		Item:      item,
	})

	if err != nil {
		return nil, err
	}

	return couple, nil
}

func (repository *RepositoryDynamoDB) FindByCoupleCode(coupleCode *string) ([]Couple, error) {
	filterEx := expression.Name("code").Equal(expression.Value(*coupleCode))

	expr, err := expression.NewBuilder().WithFilter(filterEx).Build()
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	result, err := repository.client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	})
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	var couples []Couple
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &couples); err != nil {
		return nil, err
	}
	return couples, nil
}

func (repository *RepositoryDynamoDB) FindById(id *string) (*Couple, error) {
	response, err := repository.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       repository.getKey(*id),
	})
	if err != nil {
		return nil, err
	}
	if response.Item == nil {
		return nil, &types.ResourceNotFoundException{
			Message: aws.String(fmt.Sprintf("resource not found for id: %s", *id)),
		}
	}
	result := &Couple{}
	err = attributevalue.UnmarshalMap(response.Item, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repository *RepositoryDynamoDB) SaveCoupleTransaction(couple *Couple) (*types.TransactWriteItem, error) {
	item, err := attributevalue.MarshalMap(couple)
	if err != nil {
		return nil, err
	}
	transaction := &types.TransactWriteItem{Put: &types.Put{
		TableName: aws.String(tableName),
		Item:      item,
	}}

	return transaction, nil
}

func (repository *RepositoryDynamoDB) UpdateCoupleTransaction(req *UpdateCoupleReq) (*types.TransactWriteItem, error) {
	// MarshalMap 호출 시 partitionKey 를 제외하기 위해 req.CoupleId 를 nil 로 설정
	// partitionKey 임시 저장
	partitionKey := req.Id
	req.Id = nil

	av, err := attributevalue.MarshalMap(req)
	if err != nil {
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

	transaction := &types.TransactWriteItem{
		Update: &types.Update{
			TableName:                 aws.String(tableName),
			Key:                       repository.getKey(*partitionKey),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
		},
	}

	return transaction, nil
}

func (repository *RepositoryDynamoDB) getKey(id string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: id},
	}
}
