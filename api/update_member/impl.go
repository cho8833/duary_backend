package update_member

import (
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/event"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/shared"
	"log"
	"time"
)

type UpdateMemberReq struct {
	Birthday   *time.Time         `json:"birthday"`
	Name       *string            `json:"name"`
	Character  *string            `json:"character"`
	FcmToken   *string            `json:"fcmToken"`
	MyAlarm    *event.AlarmOffset `json:"myAlarm"`
	LoverAlarm *event.AlarmOffset `json:"loverAlarm"`
}

type UpdateMemberRes struct {
	Member *member.Member `json:"member"`
	Couple *couple.Couple `json:"couple"`
}

func UpdateMember(req *UpdateMemberReq, transaction *model.DynamoDBWriteTransaction, coupleSvc couple.Service, memberSvc member.Service) (*UpdateMemberRes, shared.ApplicationError) {

	authContext := shared.GetAuthContext()

	socialId := authContext.SocialId
	provider := authContext.Provider
	coupleId := authContext.CoupleId

	// begin transaction
	transaction.BeginTransaction()

	foundCouple, svcErr := coupleSvc.FindById(*coupleId)
	if svcErr != nil {
		log.Printf(svcErr.Error())
		return nil, svcErr
	}

	// update member
	memberUpdateReq := &member.UpdateMemberReq{
		SocialId:   *socialId,
		Provider:   *provider,
		Name:       req.Name,
		Birthday:   req.Birthday,
		Character:  req.Character,
		FcmToken:   req.FcmToken,
		MyAlarm:    req.MyAlarm,
		LoverAlarm: req.LoverAlarm,
	}
	updatedMember, svcErr := memberSvc.UpdateTransaction(memberUpdateReq, transaction)
	if svcErr != nil {
		return nil, svcErr
	}

	// update couple
	var updatedMembers []*member.Member // 업데이트 대상이 되는 멤버를 제외한 새로운 리스트
	for _, m := range foundCouple.Members {
		if m.SocialId != *socialId || m.Provider != *provider {
			updatedMembers = append(updatedMembers, m)
		}
	}
	updatedMembers = append(updatedMembers, updatedMember) // 업데이트된 멤버 추가
	coupleUpdateReq := &couple.UpdateCoupleReq{
		Id:      coupleId,
		Members: updatedMembers,
	}
	updatedCouple, svcErr := coupleSvc.Update(coupleUpdateReq, transaction)
	if svcErr != nil {
		return nil, svcErr
	}

	_, err := transaction.Execute()
	if err != nil {
		log.Printf(err.Error())
		return nil, shared.DBUpdateError{}
	}

	return &UpdateMemberRes{
		Member: updatedMember,
		Couple: updatedCouple,
	}, nil

}
