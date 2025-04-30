package test

import (
	"fmt"
	"github.com/cho8833/duary_lambda/dev"
	"github.com/cho8833/duary_lambda/internal/auth/jwtutil"
	"github.com/cho8833/duary_lambda/internal/member"
	"github.com/cho8833/duary_lambda/internal/test"
	"log"
	"testing"
)

func Test_dummySignIn_Local(t *testing.T) {
	client := test.CreateLocalDynamoDBClient()

	memberRepository := member.NewMemberRepository(client)
	jwtUtil := &jwtutil.Impl{}
	dummySvc := dev.NewDummyMemberService(jwtUtil, memberRepository)

	req := &dev.DummySignInReq{
		Username: int64(4),
	}

	res, svcError := dummySvc.SignIn(req)
	if svcError != nil {
		log.Println(svcError.Error())
		t.Fail()
	} else {
		fmt.Printf("%+v\n", res)
	}

}
