package common

import (
	"github.com/cho8833/duary_lambda/internal/couple"
	"github.com/cho8833/duary_lambda/internal/member"
	"github.com/cho8833/duary_lambda/internal/util"
	"log"
)

type MemberService interface {
	UpdateMember(request *UpdateMemberReq, socialId string, provider string, coupleId string) (*UpdateMemberRes, util.ApplicationError)
}

type MemberServiceImpl struct {
	memberSvc member.Service
	coupleSvc couple.Service
}

func NewMemberService(memberSvc member.Service, coupleSvc couple.Service) *MemberServiceImpl {
	return &MemberServiceImpl{memberSvc: memberSvc, coupleSvc: coupleSvc}
}

func (svc *MemberServiceImpl) UpdateMember(request *UpdateMemberReq, socialId string, provider string, coupleId string, transaction *util.DynamoDBWriteTransaction) (*UpdateMemberRes, util.ApplicationError) {

	// begin transaction
	transaction.BeginTransaction()

	foundCouple, svcErr := svc.coupleSvc.FindById(coupleId)
	if svcErr != nil {
		log.Printf(svcErr.Error())
		return nil, svcErr
	}

	// update member
	memberUpdateReq := &member.UpdateMemberReq{
		SocialId: socialId,
		Provider: provider,
		Name:     request.Name,
		Birthday: request.Birthday,
	}
	updatedMember, svcErr := svc.memberSvc.UpdateMember(memberUpdateReq, transaction)
	if svcErr != nil {
		return nil, svcErr
	}

	// update couple
	var updatedMembers []*member.Member // 업데이트 대상이 되는 멤버를 제외한 새로운 리스트
	for _, m := range foundCouple.Members {
		if m.SocialId != socialId || m.Provider != provider {
			updatedMembers = append(updatedMembers, m)
		}
	}
	updatedMembers = append(updatedMembers, updatedMember) // 업데이트된 멤버 추가
	coupleUpdateReq := &couple.UpdateCoupleReq{
		Id:      &coupleId,
		Members: updatedMembers,
	}
	updatedCouple, svcErr := svc.coupleSvc.UpdateCouple(coupleUpdateReq, transaction)
	if svcErr != nil {
		return nil, svcErr
	}

	_, err := transaction.Execute()
	if err != nil {
		log.Printf(err.Error())
		return nil, util.DBUpdateError{}
	}

	return &UpdateMemberRes{
		Member: updatedMember,
		Couple: updatedCouple,
	}, nil

}
