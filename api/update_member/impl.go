package update_member

import (
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/event"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/model/ws_connection"
	"github.com/cho8833/duary_lambda/scheduler"
	"github.com/cho8833/duary_lambda/shared"
	"github.com/cho8833/duary_lambda/ws"
	"log"
	"strings"
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

func UpdateMember(req *UpdateMemberReq, transaction *model.DynamoDBWriteTransaction, coupleSvc couple.Service, memberSvc member.Service, eventSvc event.Service, wsRepo ws_connection.Repository, schedulerHelper scheduler.BridgeSchedulerHelper) (*UpdateMemberRes, shared.ApplicationError) {

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
	var updatedMembers []member.Member // 업데이트 대상이 되는 멤버를 제외한 새로운 리스트
	for _, m := range foundCouple.Members {
		if m.SocialId != *socialId || m.Provider != *provider {
			updatedMembers = append(updatedMembers, m)
			break
		}
	}
	updatedMembers = append(updatedMembers, *updatedMember) // 업데이트된 멤버 추가
	coupleUpdateReq := &couple.UpdateCoupleReq{
		Id:      coupleId,
		Members: updatedMembers,
	}
	updatedCouple, svcErr := coupleSvc.Update(coupleUpdateReq, transaction)
	if svcErr != nil {
		return nil, svcErr
	}

	// update birthday event
	var birthday *event.VO
	if req.Birthday != nil {
		birthday, svcErr = UpdateBirthday(*coupleId, *updatedMember, eventSvc, transaction)
		if svcErr != nil {
			return nil, svcErr
		}
	}

	_, err := transaction.Execute()
	if err != nil {
		log.Printf(err.Error())
		return nil, shared.DBUpdateError{}
	}

	// update birthday scheduler
	if birthday != nil {
		err = schedulerHelper.UpdateAnniversarySchedule(*birthday, []member.Member{*updatedMember})
		if err != nil {
			log.Printf(err.Error())
		}
	}

	res := &UpdateMemberRes{
		Member: updatedMember,
		Couple: updatedCouple,
	}

	// send updated lover ws message to lover
	wsSvc, err := ws.NewService(wsRepo)
	if err == nil {
		loginMemberId := *authContext.SocialId + "-" + *authContext.Provider
		foundCouple, svcErr := coupleSvc.FindById(*coupleId)
		if svcErr == nil {
			if len(foundCouple.Members) > 1 {
				var otherMemberId string
				for _, coupleMember := range foundCouple.Members {
					memberId := coupleMember.GetId()
					if memberId != loginMemberId {
						otherMemberId = memberId
					}
				}
				idSplit := strings.Split(otherMemberId, "-")
				err = wsSvc.Send(idSplit[0], idSplit[1], ws.LoverUpdated, res)
				if err != nil {
					log.Println("lover is not connected to ws", err)
				}
			} else {
				log.Printf("couple is not connected, stop sending WS Message")
			}
		}
	} else {
		log.Printf(err.Error())
	}

	return res, nil
}

func UpdateBirthday(coupleId string, targetMember member.Member, eventSvc event.Service, transaction *model.DynamoDBWriteTransaction) (*event.VO, shared.ApplicationError) {
	err := eventSvc.DeleteBirthdayTransaction(coupleId, targetMember.GetId(), transaction)
	if err != nil {
		log.Printf(err.Error())
		return nil, shared.DBUpdateError{}
	}

	birthdayReq := eventSvc.GenerateBirthday(coupleId, targetMember.GetId(), *targetMember.Birthday)
	birthday, svcErr := eventSvc.SaveTransaction(birthdayReq, transaction)
	if svcErr != nil {
		log.Printf(svcErr.Error())
		return nil, svcErr
	}

	return birthday, nil

}
