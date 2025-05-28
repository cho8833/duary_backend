package dummy_sign_in

import (
	"errors"
	"github.com/cho8833/duary_lambda/appjwt"
	"github.com/cho8833/duary_lambda/auth"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/shared"
	"log"
	"os"
)

type DummySignInReq struct {
	Username string  `json:"username"`
	FcmToken *string `json:"fcm_token"`
}

func SignIn(req *DummySignInReq, memberSvc member.Service, coupleSvc couple.Service, jwtUtil appjwt.JWTUtil) (*auth.SignInRes, shared.ApplicationError) {

	findMember, svcErr := memberSvc.FindById(req.Username, "kakao")
	if temp := new(shared.UserNotFound); !errors.As(svcErr, &temp) && svcErr != nil {
		log.Printf("failed to find findMember\nid:%s\nerror:%s", req.Username, svcErr.Error())
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

			memberInfo, svcErr = memberSvc.UpdateMember(updateMemberReq)
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
		}
		newMember, err := memberSvc.SaveMember(newMemberReq)
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
