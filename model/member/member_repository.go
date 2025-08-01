package member

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cho8833/duary_lambda/model"
	"log"
)

const tableName = "Member"

type Repository interface {
	FindBySocialIdAndProvider(socialId string, provider string) (*Member, error)
	Save(member *Member) (*Member, error)
	UpdateNonNil(member *UpdateMemberReq) (*Member, error)
	UpdateNonNilTransaction(member *UpdateMemberReq, transaction *model.DynamoDBWriteTransaction) (*Member, error)
	UpdateFcmToken(socialId string, provider string, fcmToken *string) (*Member, error)
	RemoveCoupleIdTransaction(socialId string, provider string, transaction *model.DynamoDBWriteTransaction) (*Member, error)
	DeleteTransaction(socialId string, provider string, transaction *model.DynamoDBWriteTransaction) error
}

type RepositoryDynamoDB struct {
	client dynamodb.Client
}

func NewRepository(client *dynamodb.Client) *RepositoryDynamoDB {
	return &RepositoryDynamoDB{client: *client}
}

func (repo *RepositoryDynamoDB) FindBySocialIdAndProvider(socialId string, provider string) (*Member, error) {

	result, err := repo.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
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
	member := &Member{}
	err = attributevalue.UnmarshalMap(result.Item, member)
	if err != nil {
		return nil, err
	}

	return member, nil
}
func (repo *RepositoryDynamoDB) Save(member *Member) (*Member, error) {
	item, err := attributevalue.MarshalMap(member)
	if err != nil {
		return nil, err
	}
	_, err = repo.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (repo *RepositoryDynamoDB) UpdateNonNil(member *UpdateMemberReq) (*Member, error) {
	expr, err := repo.updateMemberExpression(member)
	if err != nil {
		return nil, err
	}

	response, err := repo.client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       repo.getKey(member.SocialId, member.Provider),
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueAllNew,
	})
	if err != nil {
		log.Printf("failed to update item. Member: %+v, Error: %s", member, err.Error())
		return nil, err
	}

	result := &Member{}
	_ = attributevalue.UnmarshalMap(response.Attributes, result)
	return result, nil
}

func (repo *RepositoryDynamoDB) UpdateNonNilTransaction(req *UpdateMemberReq, transaction *model.DynamoDBWriteTransaction) (*Member, error) {
	// Application 단에서 업데이트 값 예측
	cacheMember, err := repo.FindBySocialIdAndProvider(req.SocialId, req.Provider)
	if err != nil {
		return nil, err
	}
	cacheMember.ApplyFrom(*req)

	// transaction
	expr, err := repo.updateMemberExpression(req)
	if err != nil {
		return nil, err
	}
	transactionItem := &types.TransactWriteItem{Update: &types.Update{
		TableName:                 aws.String(tableName),
		Key:                       repo.getKey(req.SocialId, req.Provider),
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
		UpdateExpression:          expr.Update(),
	}}
	transaction.AddTransaction(transactionItem)

	return cacheMember, nil
}

func (repo *RepositoryDynamoDB) UpdateFcmToken(socialId string, provider string, fcmToken *string) (*Member, error) {
	var input *dynamodb.UpdateItemInput
	if fcmToken != nil {
		input = &dynamodb.UpdateItemInput{
			TableName:        aws.String(tableName),
			Key:              repo.getKey(socialId, provider),
			UpdateExpression: aws.String("SET fcmToken = :fcmToken"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":fcmToken": &types.AttributeValueMemberS{Value: *fcmToken},
			},
			ReturnValues: types.ReturnValueAllNew,
		}
	} else {
		input = &dynamodb.UpdateItemInput{
			TableName:        aws.String(tableName),
			Key:              repo.getKey(socialId, provider),
			UpdateExpression: aws.String("REMOVE fcmToken"),
			ReturnValues:     types.ReturnValueAllNew,
		}
	}
	output, err := repo.client.UpdateItem(context.TODO(), input)
	if err != nil {
		log.Printf("failed to update fcm token. Error: %s", err.Error())
		return nil, err
	}
	result := &Member{}
	err = attributevalue.UnmarshalMap(output.Attributes, result)
	if err != nil {
		log.Printf("failed unmarshal member. Error: %s", err.Error())
		return nil, err
	}
	return result, nil
}

func (repo *RepositoryDynamoDB) RemoveCoupleIdTransaction(socialId string, provider string, transaction *model.DynamoDBWriteTransaction) (*Member, error) {

	cacheMember, err := repo.FindBySocialIdAndProvider(socialId, provider)
	if err != nil {
		return nil, err
	}
	transactionItem := &types.TransactWriteItem{
		Update: &types.Update{
			TableName:        aws.String(tableName),
			Key:              repo.getKey(socialId, provider),
			UpdateExpression: aws.String("REMOVE coupleId"),
		},
	}
	transaction.AddTransaction(transactionItem)

	cacheMember.CoupleId = nil
	return cacheMember, nil
}

func (repo *RepositoryDynamoDB) DeleteTransaction(socialId string, provider string, transaction *model.DynamoDBWriteTransaction) error {
	transactionItem := &types.TransactWriteItem{
		Delete: &types.Delete{
			TableName: aws.String(tableName),
			Key:       repo.getKey(socialId, provider),
		},
	}
	transaction.AddTransaction(transactionItem)

	return nil
}

func (repo *RepositoryDynamoDB) updateMemberExpression(req *UpdateMemberReq) (*expression.Expression, error) {
	update := expression.UpdateBuilder{}
	if req.Name != nil {
		update = update.Set(expression.Name("name"), expression.Value(*req.Name))
	}
	if req.AccessToken != nil {
		update = update.Set(expression.Name("accessToken"), expression.Value(req.AccessToken))
	}
	if req.Birthday != nil {
		update = update.Set(expression.Name("birthday"), expression.Value(req.Birthday))
	}
	if req.FcmToken != nil {
		update = update.Set(expression.Name("fcmToken"), expression.Value(req.FcmToken))
	}
	if req.CoupleId != nil {
		update = update.Set(expression.Name("coupleId"), expression.Value(req.CoupleId))
	}
	if req.Character != nil {
		update = update.Set(expression.Name("character"), expression.Value(req.Character))
	}
	if req.LoverAlarm != nil {
		update = update.Set(expression.Name("loverAlarm"), expression.Value(req.LoverAlarm))
	}
	if req.MyAlarm != nil {
		update = update.Set(expression.Name("myAlarm"), expression.Value(req.MyAlarm))
	}

	expr, err := expression.NewBuilder().WithUpdate(update).Build()

	if err != nil {
		log.Printf("failed to build update expression : %s", err.Error())
		return nil, err
	}
	return &expr, nil
}

func (repo *RepositoryDynamoDB) getKey(socialId string, provider string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"socialId": &types.AttributeValueMemberS{Value: socialId},
		"provider": &types.AttributeValueMemberS{Value: provider},
	}
}
