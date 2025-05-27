package couple

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

const tableName = "Couple"

type Repository interface {
	SaveCouple(couple *Couple) (*Couple, error)
	SaveCoupleTransaction(couple *Couple) (*types.TransactWriteItem, error)
	FindById(id *string) (*Couple, error)
	FindByCoupleCode(coupleCode *string) ([]Couple, error)
	UpdateCouple(req *UpdateCoupleReq) (*Couple, error)
	UpdateCoupleTransaction(req *UpdateCoupleReq, transaction *model.DynamoDBWriteTransaction) (*Couple, error)
	DeleteCoupleTransaction(id string) (*types.TransactWriteItem, error)
}

type RepositoryDynamoDB struct {
	*model.DynamoDBRepository[Couple]
	client *dynamodb.Client
}

func NewCoupleRepository(client *dynamodb.Client) *RepositoryDynamoDB {
	return &RepositoryDynamoDB{client: client, DynamoDBRepository: model.NewBaseDynamoRepository[Couple](client, tableName)}
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

	// caching
	repository.SetToCache(*couple.Id, couple)

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

	// caching
	for _, item := range couples {
		repository.SetToCache(*item.Id, &item)
	}
	return couples, nil
}

func (repository *RepositoryDynamoDB) FindById(id *string) (*Couple, error) {
	result, err := repository.FindByIdCaching(context.TODO(), *id)
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

	repository.SetToCache(*couple.Id, couple)
	return transaction, nil
}

func (repository *RepositoryDynamoDB) DeleteCoupleTransaction(id string) (*types.TransactWriteItem, error) {
	transaction := &types.TransactWriteItem{
		Delete: &types.Delete{
			TableName: aws.String(tableName),
			Key:       repository.getKey(id),
		}}

	return transaction, nil
}

func (repository *RepositoryDynamoDB) UpdateCouple(couple *UpdateCoupleReq) (*Couple, error) {
	// MarshalMap 호출 시 partitionKey 를 제외하기 위해 req.CoupleId 를 nil 로 설정
	// partitionKey 임시 저장
	partitionKey := *couple.Id
	couple.Id = nil

	av, err := attributevalue.MarshalMap(couple)
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

	response, err := repository.client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       repository.getKey(partitionKey),
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueAllNew,
	})
	if err != nil {
		log.Printf("failed to update item. Couple: %+v, Error: %s", couple, err.Error())
		return nil, err
	}

	result := &Couple{}
	_ = attributevalue.UnmarshalMap(response.Attributes, result)
	return result, nil
}

func (repository *RepositoryDynamoDB) UpdateCoupleTransaction(req *UpdateCoupleReq, transaction *model.DynamoDBWriteTransaction) (*Couple, error) {
	// Application 단에서 업데이트 값 예측
	cacheCouple, err := repository.FindByIdCaching(context.TODO(), *req.Id)
	if err != nil {
		return nil, err
	}
	cacheCouple.ApplyFrom(*req)

	// MarshalMap 호출 시 partitionKey 를 제외하기 위해 req.CoupleId 를 nil 로 설정
	// partitionKey 임시 저장
	partitionKey := *req.Id
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

	transactionItem := &types.TransactWriteItem{
		Update: &types.Update{
			TableName:                 aws.String(tableName),
			Key:                       repository.getKey(partitionKey),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
		},
	}
	transaction.AddTransaction(transactionItem)

	repository.SetToCache(*cacheCouple.Id, cacheCouple)

	return cacheCouple, nil
}

func (repository *RepositoryDynamoDB) getKey(id string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: id},
	}
}
