package connect_couple

import (
	"github.com/aws/smithy-go/ptr"
	"github.com/cho8833/duary_lambda/appjwt"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/event"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/shared"
	"log"
	"os"
)

type ConnectCoupleReq struct {
	CoupleCode *string `json:"coupleCode"`
}

type StartDuaryRes struct {
	Member *member.Member         `json:"member"`
	Couple *couple.Couple         `json:"couple"`
	Token  *appjwt.ApplicationJWT `json:"token"`
}

func ConnectCouple(req *ConnectCoupleReq, transaction *model.DynamoDBWriteTransaction, coupleSvc couple.Service, memberSvc member.Service, eventSvc event.Service) (*StartDuaryRes, shared.ApplicationError) {
	if req.CoupleCode == nil || len(*req.CoupleCode) == 0 {
		return nil, shared.BadRequestError{}
	}

	authContext := shared.GetAuthContext()

	// findMember.coupleId == nil 은 startDuary 를 거치지 않았다는 의미 => bad request
	if authContext.CoupleId == nil {
		log.Printf("member.coupleId == nil")
		return nil, shared.BadRequestError{}
	}

	// Find Couple
	targetCouple, svcError := coupleSvc.FindByCoupleCode(req.CoupleCode)
	if svcError != nil {
		log.Printf("coupleSvc.FindByCoupleCode err: %v", svcError)
		return nil, shared.CoupleCodeNotFound{}
	}

	// 자신의 coupleId 와 Code 의 CoupleId 가 같음 -> 자신의 couple 에 연결한다 -> bad request
	if *authContext.CoupleId == *targetCouple.Id {
		log.Printf("member.coupleId == targetCouple.Id %v", *authContext.CoupleId)
		return nil, shared.BadRequestError{}
	}

	// Couple 이 이미 연결되어 있는 경우 Bad Request
	if len(targetCouple.Members) >= 2 {
		log.Printf("couple is already connected to %v", targetCouple.Members)
		return nil, shared.BadRequestError{}
	}

	transaction.BeginTransaction()

	// 기존의 Couple 삭제
	svcErr := coupleSvc.DeleteCouple(*authContext.CoupleId, transaction)
	if svcErr != nil {
		return nil, shared.DBUpdateError{}
	}

	// Update Member
	updateReq := &member.UpdateMemberReq{
		Provider: *authContext.Provider,
		SocialId: *authContext.SocialId,
		CoupleId: targetCouple.Id,
	}
	updatedMember, svcError := memberSvc.UpdateMemberTransaction(updateReq, transaction)
	if svcError != nil {
		return nil, svcError
	}

	// Update Couple
	updateCoupleReq := &couple.UpdateCoupleReq{
		Id:      targetCouple.Id,
		Members: append(targetCouple.Members, updatedMember), // Member 추가
		Code:    ptr.String(""),                              // 코드 삭제
	}
	updatedCouple, svcError := coupleSvc.UpdateCouple(updateCoupleReq, transaction)
	if svcError != nil {
		return nil, svcError
	}

	// Create Anniversary Event
	anniversary100DayReq := &event.SaveReq{
		Recurrence: &event.Recurrence{
			RepeatStartDate: *updatedCouple.RelationDate,
			Interval:        100,
			Frequency:       "DAILY",
		},
		StartDateTime: *updatedCouple.RelationDate,
		IsAllDay:      true,
		IsTogether:    true,
		CreatedBy:     updatedCouple.Members[0].GetId(),
		Title:         "100days",
		EventType:     "ANNIVERSARY",
		CoupleId:      *targetCouple.Id,
	}
	anniversaryYearlyReq := &event.SaveReq{
		Recurrence: &event.Recurrence{
			RepeatStartDate: *updatedCouple.RelationDate,
			Interval:        1,
			Frequency:       "YEARLY",
		},
		StartDateTime: *updatedCouple.RelationDate,
		IsAllDay:      true,
		CreatedBy:     updatedCouple.Members[0].GetId(),
		Title:         "year",
		IsTogether:    true,
		EventType:     "ANNIVERSARY",
		CoupleId:      *targetCouple.Id,
	}
	_, svcErr = eventSvc.SaveTransaction(anniversary100DayReq, transaction)
	if svcErr != nil {
		return nil, svcErr
	}
	_, svcErr = eventSvc.SaveTransaction(anniversaryYearlyReq, transaction)
	if svcErr != nil {
		return nil, svcErr
	}

	// execute transaction
	_, err := transaction.Execute()
	if err != nil {
		log.Println(err)
		return nil, shared.DBUpdateError{}
	}

	// generate new token with couple id
	jwtUtil := &appjwt.Impl{}
	key := os.Getenv("secretKey")
	newToken := jwtUtil.NewToken(jwtUtil.GenerateSubject(updatedMember.SocialId, updatedMember.Provider), updatedCouple.Id, key)

	result := &StartDuaryRes{
		Member: updatedMember,
		Couple: updatedCouple,
		Token:  newToken,
	}

	return result, nil
}
