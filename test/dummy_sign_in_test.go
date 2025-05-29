package test

import (
	"dummy_sign_in"
	"fmt"
	"github.com/cho8833/duary_lambda/appjwt"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/shared"
	"log"
	"testing"
)

func Test_dummySignIn_Local(t *testing.T) {
	client := CreateLocalDynamoDBClient()

	memberRepository := member.NewMemberRepository(client)
	coupleRepo := couple.NewCoupleRepository(client)

	coupleSvc := couple.NewCoupleService(coupleRepo)
	memberSvc := member.NewMemberService(memberRepository)

	req := &dummy_sign_in.DummySignInReq{
		Username: "124",
		FcmToken: ptr("test"),
	}
	jwtUtil := &appjwt.Impl{}

	res, svcError := dummy_sign_in.SignIn(req, memberSvc, coupleSvc, jwtUtil)
	if svcError != nil {
		log.Println(svcError.Error())
		t.Fail()
	} else {
		fmt.Printf(shared.ToString(res))
	}

}
