package test

import (
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/member"
	"start_duary"
	"testing"
	"time"
)

func Test_startDuary(t *testing.T) {
	dynamodbClient := CreateLocalDynamoDBClient()

	transaction := model.NewWriteTransaction(dynamodbClient)

	coupleRepo := couple.NewRepository(dynamodbClient)
	memberRepo := member.NewRepository(dynamodbClient)
	coupleSvc := couple.NewService(coupleRepo)
	memberSvc := member.NewService(memberRepo)

	dateTemp := time.Now()
	dateTemp = dateTemp.AddDate(0, -1, 0)

	CreateTestAuthContext(ptr("asdf"), ptr("kakao"), nil)

	req := &start_duary.StartDuaryReq{
		RelationDate: &dateTemp,
		Birthday:     &dateTemp,
		MyCharacter:  ptr("BLUE"),
		Name:         ptr("test"),
	}

	res, svcErr := start_duary.StartDuary(req, transaction, coupleSvc, memberSvc)
	if svcErr != nil {
		t.Error(svcErr)
	}
	print(res)

}
