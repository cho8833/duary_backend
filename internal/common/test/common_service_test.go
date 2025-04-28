package test

import (
	"fmt"
	"github.com/cho8833/duary_lambda/internal/common"
	"github.com/cho8833/duary_lambda/internal/couple"
	"github.com/cho8833/duary_lambda/internal/member"
	"github.com/cho8833/duary_lambda/internal/test"
	util2 "github.com/cho8833/duary_lambda/internal/util"
	"testing"
	"time"
)

func Test_InitDuary_Local(t *testing.T) {
	client := test.CreateLocalDynamoDBClient()

	transaction := util2.NewWriteTransaction(client)
	coupleRepo := couple.NewCoupleRepository(client)
	memberRepo := member.NewMemberRepository(client)
	coupleSvc := couple.NewCoupleService(coupleRepo)
	memberSvc := member.NewMemberService(memberRepo)

	commonSvc := common.NewCommonService(memberSvc, coupleSvc)

	socialId := int64(3)
	provider := "kakao"
	test.CreateTestAuthContext(&socialId, &provider, nil)

	name := "test"
	birthday := time.Now()
	relationDate := time.Now()
	myCharacter := "yellow"
	req := &common.StartDuaryReq{
		Birthday:     &birthday,
		RelationDate: &relationDate,
		Name:         &name,
		MyCharacter:  &myCharacter,
	}
	res, err := commonSvc.StartDuary(req, transaction)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("%+v", res)

}
