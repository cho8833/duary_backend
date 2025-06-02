package auth

import (
	"errors"
	"fmt"
	"github.com/cho8833/duary_lambda/appjwt"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/shared"
	"log"
	"os"
)

type Service interface {
	SignIn(token *KakaoOAuthToken) (SignInRes, shared.ApplicationError)
}

type ServiceImpl struct {
	memberSvc    member.Service
	coupleSvc    couple.Service
	jwtValidator appjwt.JWTValidator
	jwtUtil      appjwt.JWTUtil
}

func NewAuthService(jwtValidator appjwt.JWTValidator,
	jwtUtil appjwt.JWTUtil,
	memberSvc member.Service, coupleSvc couple.Service) *ServiceImpl {
	return &ServiceImpl{jwtValidator: jwtValidator, jwtUtil: jwtUtil, memberSvc: memberSvc, coupleSvc: coupleSvc}
}

func (svc *ServiceImpl) KakaoSignIn(kakaoToken *KakaoOAuthToken, fcmToken *string) (*SignInRes, shared.ApplicationError) {

	aud := os.Getenv("aud")
	nonce := os.Getenv("nonce")
	// verify id token
	validateValue := &appjwt.ValidatingValue{
		Url:      "https://kauth.kakao.com/.well-known/jwks.json",
		Aud:      aud,
		Nonce:    nonce,
		Iss:      "https://kauth.kakao.com",
		Provider: "kakao",
	}
	payload, err := svc.jwtValidator.VerifyRSA256(*kakaoToken.IdToken, validateValue)
	if err != nil {
		log.Printf("failed to verify token. idToken: %s, error: %s", kakaoToken.IdToken, err.Error())
		return nil, shared.BadRequestError{}
	}

	res, svcError := svc.onSignInSuccess(payload, "kakao", fcmToken)
	if svcError != nil {
		return nil, svcError
	}

	return res, nil

}

const appId = "com.ivis.duary"

func (svc *ServiceImpl) AppleSignIn(token *AppleOAuthToken, fcmToken *string) (*SignInRes, shared.ApplicationError) {
	nonce := os.Getenv("nonce")

	validateValue := &appjwt.ValidatingValue{
		Url:      "https://appleid.apple.com/auth/oauth2/v2/keys",
		Nonce:    nonce,
		Iss:      "https://appleid.apple.com",
		Provider: "apple",
		Aud:      appId,
	}
	payload, err := svc.jwtValidator.VerifyRSA256(*token.IdentityToken, validateValue)
	if err != nil {
		log.Printf("failed to verify token. idToken: %s, error: %s", *token.IdentityToken, err.Error())
		return nil, shared.BadRequestError{}
	}

	res, svcError := svc.onSignInSuccess(payload, "apple", fcmToken)
	if svcError != nil {
		return nil, svcError
	}
	return res, nil

}

func (svc *ServiceImpl) onSignInSuccess(payload *appjwt.DecodedPayload, provider string, fcmToken *string) (*SignInRes, shared.ApplicationError) {
	// 회원 ID 와 ServiceProvider 로 Member 검색
	// Member 가 없을 경우 ResourceNotFoundException 발생, 해당 Exception 은 오류가 아님
	findMember, svcErr := svc.memberSvc.FindById(payload.SocialId, provider)
	if temp := new(shared.UserNotFound); !errors.As(svcErr, &temp) && svcErr != nil {
		id := fmt.Sprintf("%d-%s", payload.SocialId, provider)
		log.Printf("failed to find findMember\nid:%s\nerror:%s", id, svcErr.Error())
		return nil, svcErr
	}

	// findMember 가 존재하는 경우 DB 필드를 업데이트하고 이미 회원가입된 Member return
	if findMember != nil {

		var myCouple *couple.Couple
		if findMember.CoupleId != nil {
			myCouple, svcErr = svc.coupleSvc.FindById(*findMember.CoupleId)
			if svcErr != nil {
				return nil, svcErr
			}
		}

		// generate application token
		memberId := svc.jwtUtil.GenerateSubject(findMember.SocialId, findMember.Provider)
		key := os.Getenv("secretKey")
		newToken := svc.jwtUtil.NewToken(memberId, findMember.CoupleId, key)

		var memberInfo *member.Member
		if fcmToken == nil {
			memberInfo = findMember
		} else {
			// member 의 fcm_token update
			updateMemberReq := member.UpdateMemberReq{
				FcmToken: fcmToken,
				SocialId: payload.SocialId,
				Provider: provider,
			}
			memberInfo, svcErr = svc.memberSvc.UpdateMember(&updateMemberReq)
			if svcErr != nil {
				log.Printf("failed to save findMember\nfindMember: %+v\nerror: %s", findMember, svcErr)
				return nil, shared.DBSaveError{}
			}
		}
		result := &SignInRes{
			Member:     memberInfo,
			IsRegister: false,
			Token:      newToken,
			Couple:     myCouple,
		}
		return result, nil
	} else {
		// findMember 가 존재하지 않는 경우 Member 생성, 최초 회원가입
		newMemberReq := &member.SaveMemberReq{
			Name:       payload.NickName,
			Birthday:   nil,
			Provider:   provider,
			SocialId:   payload.SocialId,
			FcmToken:   fcmToken,
			LoverAlarm: member.DefaultAlarmOffset,
			MyAlarm:    member.DefaultAlarmOffset,
		}
		newMember, err := svc.memberSvc.SaveMember(newMemberReq)
		if err != nil {
			log.Printf("failed to save findMember\nnew findMember: %+v\nerror: %s", newMember, err.Error())
			return nil, shared.DBSaveError{}
		}

		// generate application token
		memberId := svc.jwtUtil.GenerateSubject(newMember.SocialId, newMember.Provider)
		key := os.Getenv("secretKey")
		newToken := svc.jwtUtil.NewToken(memberId, newMember.CoupleId, key)

		result := &SignInRes{
			Member:     newMember,
			IsRegister: true,
			Token:      newToken,
		}
		return result, nil
	}
}
