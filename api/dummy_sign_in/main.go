package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_backend/appjwt"
	"github.com/cho8833/duary_backend/auth"
	"github.com/cho8833/duary_backend/model"
	"github.com/cho8833/duary_backend/model/couple"
	"github.com/cho8833/duary_backend/model/member"
	"github.com/cho8833/duary_backend/shared"
	"log"
	"os"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap main.go && chmod 755 bootstrap && zip  ../../build/package/dummy_sign_in_api.zip bootstrap && rm bootstrap
*/

type DummySignInReq struct {
	Username string  `json:"username"`
	FcmToken *string `json:"fcmToken"`
}

func dummySignIn(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamodbClient, _ := model.GetDynamoDBClient()
	stage := request.StageVariables["stage"]

	memberRepo := member.NewRepository(dynamodbClient, stage)
	coupleRepo := couple.NewRepository(dynamodbClient, stage)

	memberSvc := member.NewService(memberRepo)
	coupleSvc := couple.NewService(coupleRepo)

	jwtUtil := &appjwt.Impl{}

	req := &DummySignInReq{}
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	res, svcError := SignIn(req, memberSvc, coupleSvc, jwtUtil)
	if svcError != nil {
		return shared.LambdaAppErrorResponse(svcError), nil
	}

	return shared.LambdaResponseWithDataAndHeader(res, appjwt.ApplicationJWTToHeader(*res.Token)), nil
}

func SignIn(req *DummySignInReq, memberSvc member.Service, coupleSvc couple.Service, jwtUtil appjwt.JWTUtil) (*auth.SignInRes, shared.ApplicationError) {

	findMember, svcErr := memberSvc.FindById(req.Username, "kakao")
	if temp := new(shared.UserNotFound); !errors.As(svcErr, temp) && svcErr != nil {
		log.Printf("failed to find findMember. id:%s\nerror:%s", req.Username, svcErr.Error())
		return nil, shared.DBReadError{}
	}

	if findMember != nil {
		// generate application token
		memberId := jwtUtil.GenerateSubject(findMember.SocialId, findMember.Provider)
		key := os.Getenv("secretKey")
		newToken := jwtUtil.NewToken(memberId, findMember.CoupleId, key)

		var memberCouple *couple.Couple
		if findMember.CoupleId != nil {
			memberCouple, svcErr = coupleSvc.FindById(*findMember.CoupleId)
			if svcErr != nil {
				return nil, shared.DBReadError{}
			}
		}

		var memberInfo *member.Member
		if req.FcmToken == nil {
			memberInfo = findMember
		} else {
			updateMemberReq := &member.UpdateMemberReq{
				FcmToken: req.FcmToken,
				Provider: findMember.Provider,
				SocialId: findMember.SocialId,
			}

			memberInfo, svcErr = memberSvc.Update(updateMemberReq)
			if svcErr != nil {
				return nil, shared.DBUpdateError{}
			}
		}

		result := &auth.SignInRes{
			Member:     memberInfo,
			IsRegister: false,
			Token:      newToken,
			Couple:     memberCouple,
		}
		return result, nil
	} else {
		// findMember 가 존재하지 않는 경우 Member 생성, 최초 회원가입
		newMemberReq := &member.SaveMemberReq{
			Name:        nil,
			Birthday:    nil,
			AccessToken: nil,
			Provider:    "kakao",
			SocialId:    req.Username,
			FcmToken:    req.FcmToken,
			MyAlarm:     member.DefaultAlarmOffset,
			LoverAlarm:  member.DefaultAlarmOffset,
		}
		newMember, svcErr := memberSvc.Save(newMemberReq)
		if svcErr != nil {
			log.Printf("failed to save findMember\nnew findMember: %+v\nerror: %s", newMember, svcErr.Error())
			return nil, shared.DBSaveError{}
		}

		// generate application token
		memberId := jwtUtil.GenerateSubject(newMember.SocialId, newMember.Provider)
		key := os.Getenv("secretKey")
		newToken := jwtUtil.NewToken(memberId, newMember.CoupleId, key)

		result := &auth.SignInRes{
			Member:     newMember,
			IsRegister: true,
			Token:      newToken,
		}
		return result, nil
	}
}

func main() {
	lambda.Start(dummySignIn)
}
