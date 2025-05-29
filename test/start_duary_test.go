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

	coupleRepo := couple.NewCoupleRepository(dynamodbClient)
	memberRepo := member.NewMemberRepository(dynamodbClient)
	coupleSvc := couple.NewCoupleService(coupleRepo)
	memberSvc := member.NewMemberService(memberRepo)

	dateTemp := time.Now()
	dateTemp = dateTemp.AddDate(0, -1, 0)

	CreateTestAuthContext(ptr("123"), ptr("kakao"), nil)

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
