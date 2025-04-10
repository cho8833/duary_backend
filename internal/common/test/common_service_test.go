package test

import (
	"fmt"
	"github.com/cho8833/duary_lambda/internal/common"
	"github.com/cho8833/duary_lambda/internal/couple"
	"github.com/cho8833/duary_lambda/internal/member"
	"github.com/cho8833/duary_lambda/internal/test/util"
	util2 "github.com/cho8833/duary_lambda/internal/util"
	"testing"
	"time"
)

func Test_InitDuary_Local(t *testing.T) {
	client := util.CreateLocalDynamoDBClient()

	transaction := util2.NewWriteTransaction(client)
	coupleRepo := couple.NewCoupleRepository(client)
	memberRepo := member.NewMemberRepository(client)
	coupleSvc := couple.NewCoupleService(coupleRepo)
	memberSvc := member.NewMemberService(memberRepo)

	commonSvc := common.NewCommonService(memberSvc, coupleSvc, nil)

	name := "test"
	otherCharacter := "blue"
	birthday := time.Now()
	relationDate := time.Now()
	myCharacter := "yellow"
	req := &common.StartDuaryReq{
		Birthday:       &birthday,
		RelationDate:   &relationDate,
		Name:           &name,
		OtherCharacter: &otherCharacter,
		MyCharacter:    &myCharacter,
		Provider:       "kakao",
		SocialId:       1,
	}
	res, err := commonSvc.StartDuary(req, transaction)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("%+v", res)

}
