package shared

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
)

type DynamoDBWriteTransaction struct {
	transaction *dynamodb.TransactWriteItemsInput
	client      *dynamodb.Client
}

func NewWriteTransaction(client *dynamodb.Client) *DynamoDBWriteTransaction {
	return &DynamoDBWriteTransaction{client: client}
}
func (t *DynamoDBWriteTransaction) BeginTransaction() {
	t.transaction = &dynamodb.TransactWriteItemsInput{}
}

func (t *DynamoDBWriteTransaction) AddTransaction(transaction *types.TransactWriteItem) {
	t.transaction.TransactItems = append(t.transaction.TransactItems, *transaction)
}

func (t *DynamoDBWriteTransaction) Execute() (*dynamodb.TransactWriteItemsOutput, error) {
	output, err := t.client.TransactWriteItems(context.TODO(), t.transaction)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	t.transaction = &dynamodb.TransactWriteItemsInput{}
	return output, nil
}
