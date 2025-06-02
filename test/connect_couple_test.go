package test

import (
	"connect_couple"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/event"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/shared"
	"testing"
)

func Test_ConnectCouple(t *testing.T) {
	dynamodbClient := CreateLocalDynamoDBClient()

	transaction := model.NewWriteTransaction(dynamodbClient)
	memberRepo := member.NewMemberRepository(dynamodbClient)
	coupleRepo := couple.NewCoupleRepository(dynamodbClient)
	eventRepo := event.NewEventRepository(dynamodbClient)

	memberSvc := member.NewMemberService(memberRepo)
	coupleSvc := couple.NewCoupleService(coupleRepo)
	eventSvc := event.NewEventService(eventRepo)

	req := &connect_couple.ConnectCoupleReq{
		CoupleCode: ptr("0unz465jw"),
	}
	CreateTestAuthContext(ptr("124"), ptr("kakao"), ptr("e2ca39c3-0dbe-4aa0-b882-d3fe6e0b04e1"))

	res, svcErr := connect_couple.ConnectCouple(req, transaction, coupleSvc, memberSvc, eventSvc)
	if svcErr != nil {
		t.Fatalf(svcErr.Error())
	} else {
		print(shared.ToString(res))
	}
}
