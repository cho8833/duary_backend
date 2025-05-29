package test

import (
	"fmt"
	"github.com/cho8833/duary_lambda/model/couple"
	"testing"
)

func Test_FindByCoupleCode(t *testing.T) {
	dynamodbClient := CreateLocalDynamoDBClient()

	repository := couple.NewCoupleRepository(dynamodbClient)

	coupleCode := "72itgwj0t"
	couples, err := repository.FindByCoupleCode(&coupleCode)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("%+v", couples)
}
