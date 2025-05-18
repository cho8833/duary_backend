package base

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
)

type DynamoDBRepository[T any] struct {
	client *dynamodb.Client
	table  string
	cache  map[string]*T
}

func NewBaseDynamoRepository[T any](client *dynamodb.Client, table string) *DynamoDBRepository[T] {
	return &DynamoDBRepository[T]{
		client: client,
		table:  table,
		cache:  make(map[string]*T),
	}
}

func (r *DynamoDBRepository[T]) FindByIdCaching(ctx context.Context, id string) (*T, error) {
	// 캐시 확인
	if item, ok := r.cache[id]; ok {
		return item, nil
	}

	// DynamoDB에서 조회
	response, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.table,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	if response.Item == nil {
		return nil, &types.ResourceNotFoundException{
			Message: aws.String(fmt.Sprintf("resource not found for id: %s", id)),
		}
	}

	var result T
	err = attributevalue.UnmarshalMap(response.Item, &result)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	// 캐시에 저장
	r.cache[id] = &result

	return &result, nil
}
func (r *DynamoDBRepository[T]) DeleteByID(ctx context.Context, id string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &r.table,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return err
	}

	// 캐시도 제거
	r.InvalidateCache(id)
	return nil
}

func (r *DynamoDBRepository[T]) InvalidateCache(id string) {
	delete(r.cache, id)
}

func (r *DynamoDBRepository[T]) SetToCache(id string, item *T) {
	r.cache[id] = item
}
