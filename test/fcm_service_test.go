package test

import (
	"github.com/cho8833/duary_lambda/fcm"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/member"
	"testing"
)

func Test_SendFCM(t *testing.T) {
	fcmClient, err := fcm.GetFCMClient()
	if err != nil {
		t.Fatal(err)
	}
	dbClient, err := model.GetDynamoDBClient()
	if err != nil {
		t.Fatal(err)
	}
	mRepository := member.NewRepository(dbClient)
	mService := member.NewService(mRepository)

	fcmService := fcm.NewService(fcmClient, mService)

	memberId := "3428835809-kakao"

	req := fcm.SendReq{
		TargetMemberId: memberId,
		Title:          "test",
		Body:           "test",
	}
	fcmService.Send(req)

}
