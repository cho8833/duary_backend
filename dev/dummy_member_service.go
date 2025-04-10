package dev

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cho8833/duary_lambda/internal/auth"
	"github.com/cho8833/duary_lambda/internal/auth/jwtutil"
	"github.com/cho8833/duary_lambda/internal/member"
	"github.com/cho8833/duary_lambda/internal/util"
	"log"
	"os"
)

type DummyMemberService interface {
	SignIn(req *DummySignInReq) (auth.SignInRes, util.ApplicationError)
}

type DummyMemberServiceImpl struct {
	memberRepository member.Repository
	jwtUtil          jwtutil.JWTUtil
}

func NewDummyMemberService(jwtUtil jwtutil.JWTUtil, memberRepository member.Repository) *DummyMemberServiceImpl {
	return &DummyMemberServiceImpl{jwtUtil: jwtUtil, memberRepository: memberRepository}
}

func (svc *DummyMemberServiceImpl) SignIn(req *DummySignInReq) (*auth.SignInRes, util.ApplicationError) {

	findMember, err := svc.memberRepository.FindBySocialIdAndProvider(req.Username, "kakao")
	if temp := new(types.ResourceNotFoundException); !errors.As(err, &temp) && err != nil {
		log.Printf("failed to find findMember\nid:%s\nerror:%s", req.Username, err.Error())
		return nil, util.DBReadError{}
	}

	if findMember != nil {
		// generate application token
		memberId := svc.jwtUtil.GenerateSubject(findMember)
		key := os.Getenv("secretKey")
		newToken := svc.jwtUtil.NewToken(memberId, findMember.CoupleId, key)
		_, err := svc.memberRepository.SaveMember(findMember)
		if err != nil {
			log.Printf("failed to save findMember\nfindMember: %+v\nerror: %s", findMember, err.Error())
			return nil, util.DBSaveError{}
		}
		result := &auth.SignInRes{
			Member:     findMember,
			IsRegister: false,
			Token:      newToken,
		}
		return result, nil
	} else {
		// findMember 가 존재하지 않는 경우 Member 생성, 최초 회원가입
		newMember := &member.Member{
			Name:        nil,
			Birthday:    nil,
			AccessToken: nil,
			Provider:    "kakao",
			Gender:      nil,
			SocialId:    req.Username,
			FcmToken:    nil,
			Email:       nil,
		}
		_, err := svc.memberRepository.SaveMember(newMember)
		if err != nil {
			log.Printf("failed to save findMember\nnew findMember: %+v\nerror: %s", newMember, err.Error())
			return nil, util.DBSaveError{}
		}

		// generate application token
		memberId := svc.jwtUtil.GenerateSubject(newMember)
		key := os.Getenv("secretKey")
		newToken := svc.jwtUtil.NewToken(memberId, newMember.CoupleId, key)

		result := &auth.SignInRes{
			Member:     newMember,
			IsRegister: true,
			Token:      newToken,
		}
		return result, nil
	}
}
