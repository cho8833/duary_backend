package test

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/cho8833/duary_lambda/member"
	"github.com/cho8833/duary_lambda/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_FindById(t *testing.T) {
	client := test.CreateLocalDynamoDBClient()
	repository := member.NewMemberRepository(client)

	member, err := repository.FindBySocialIdAndProvider(1, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("%+v", member)
}

func Test_saveMember(t *testing.T) {
	client := test.CreateLocalDynamoDBClient()
	repository := member.NewMemberRepository(client)
	name := "test"
	now := time.Now()
	gender := "man"
	dummyMember := &member.Member{
		SocialId: 1,
		Name:     &name,
		Birthday: &now,
		Gender:   &gender,
		Provider: "kakao",
		FcmToken: nil,
	}

	_, err := repository.SaveMember(dummyMember)

	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_findBySocialIdAndProvider(t *testing.T) {
	client := test.CreateLocalDynamoDBClient()
	repo := member.NewMemberRepository(client)

	member, err := repo.FindBySocialIdAndProvider(1, "kakao")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("%+v", member)
}

func Test_updateItem(t *testing.T) {
	dynamodbClient := test.CreateLocalDynamoDBClient()
	repo := member.NewMemberRepository(dynamodbClient)

	dummyUpdate := &member.UpdateMemberReq{
		Provider: "kakao",
		SocialId: 1,
		Name:     aws.String("Name"),
		FcmToken: aws.String("FCMTOKEN"),
	}
	result, err := repo.UpdateMember(dummyUpdate)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	assert.Equal(t, "Name", *result.Name)
	assert.Equal(t, "FCMTOKEN", *result.FcmToken)
}
