package dummy_sign_in

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cho8833/duary_lambda/appjwt"
	"github.com/cho8833/duary_lambda/auth"
	"github.com/cho8833/duary_lambda/member"
	"github.com/cho8833/duary_lambda/shared"
	"log"
	"os"
)

type DummySignInReq struct {
	Username string `json:"username"`
}

func SignIn(req *DummySignInReq, memberRepository member.Repository, jwtUtil appjwt.JWTUtil) (*auth.SignInRes, shared.ApplicationError) {

	findMember, err := memberRepository.FindBySocialIdAndProvider(req.Username, "kakao")
	if temp := new(types.ResourceNotFoundException); !errors.As(err, &temp) && err != nil {
		log.Printf("failed to find findMember\nid:%s\nerror:%s", req.Username, err.Error())
		return nil, shared.DBReadError{}
	}

	if findMember != nil {
		// generate application token
		memberId := jwtUtil.GenerateSubject(findMember.SocialId, findMember.Provider)
		key := os.Getenv("secretKey")
		newToken := jwtUtil.NewToken(memberId, findMember.CoupleId, key)
		_, err := memberRepository.SaveMember(findMember)
		if err != nil {
			log.Printf("failed to save findMember\nfindMember: %+v\nerror: %s", findMember, err.Error())
			return nil, shared.DBSaveError{}
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
		_, err := memberRepository.SaveMember(newMember)
		if err != nil {
			log.Printf("failed to save findMember\nnew findMember: %+v\nerror: %s", newMember, err.Error())
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
