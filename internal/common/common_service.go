package common

import (
	"github.com/cho8833/duary_lambda/internal/auth"
	"github.com/cho8833/duary_lambda/internal/auth/jwtutil"
	"github.com/cho8833/duary_lambda/internal/couple"
	"github.com/cho8833/duary_lambda/internal/member"
	"github.com/cho8833/duary_lambda/internal/util"
	"os"
)

type Service interface {

	/*
		Duary 시작하기

		1. create Couple(CoupleId, RelationDate, Add Member to Members, Code)
			- put random generated CoupleId
			- put RelationDate from StartDuaryReq
			- add Member to Members
			- put random generated Code
		2. update Member
			- put Character from StartDuaryReq
			- put name from StartDuaryReq
			- put CoupleId from (1)

		return 생성된 Couple 정보
	*/
	InitDuaryInfo(request *StartDuaryReq, transaction *util.DynamoDBWriteTransaction) (*StartDuaryRes, util.ApplicationError)

	/*
		커플 연결

		1. get Couple with ConnectCoupleReq.CoupleCode
		2. delete existing Couple of LoginMember
		3. put CoupleId of retrieved Couple to LoginMember.CoupleId
		4. put LoginMember to Couple.Members

		return 업데이트된 couple 정보
	*/
	ConnectCouple(loginMember *auth.LoginMember, req *ConnectCoupleReq) (*StartDuaryRes, util.ApplicationError)

	CheckConnected(loginMember *auth.LoginMember) (*CheckConnectedRes, util.ApplicationError)
}

type ServiceImpl struct {
	memberSvc member.Service
	coupleSvc couple.Service
}

func NewCommonService(memberSvc member.Service, coupleSvc couple.Service) *ServiceImpl {
	return &ServiceImpl{memberSvc: memberSvc, coupleSvc: coupleSvc}
}

func (svc *ServiceImpl) StartDuary(request *StartDuaryReq, transaction *util.DynamoDBWriteTransaction) (*StartDuaryRes, util.ApplicationError) {
	// begin transaction
	transaction.BeginTransaction()

	// create new couple
	coupleReq := &couple.CreateCoupleReq{
		RelationDate: *request.RelationDate,
	}
	newCouple, err := svc.coupleSvc.CreateCouple(coupleReq, transaction)
	if err != nil {
		return nil, err
	}

	// update member
	memberReq := &member.UpdateMemberReq{
		CoupleId:  newCouple.Id,
		Name:      request.Name,
		Birthday:  request.Birthday,
		SocialId:  request.SocialId,
		Character: request.MyCharacter,
		Provider:  request.Provider,
	}
	updatedMember, err := svc.memberSvc.UpdateMember(memberReq, transaction)
	if err != nil {
		return nil, err
	}

	// put updated Member in Couple.Members
	updateCoupleReq := &couple.UpdateCoupleReq{
		Members: []*member.Member{updatedMember},
	}
	err = svc.coupleSvc.UpdateCouple(updateCoupleReq, transaction)
	if err != nil {
		return nil, err
	}

	// execute transaction
	_, transactionError := transaction.Execute()
	if transactionError != nil {
		return nil, util.DBError{}
	}

	jwtUtil := &jwtutil.Impl{}
	key := os.Getenv("secretKey")
	appToken := jwtUtil.NewToken(jwtUtil.GenerateSubject(updatedMember), updatedMember.CoupleId, key)

	res := &StartDuaryRes{
		Member: updatedMember,
		Couple: newCouple,
		Token:  appToken,
	}
	return res, nil
}

func (svc *ServiceImpl) ConnectCouple(loginMember *auth.LoginMember, req *ConnectCoupleReq, transaction *util.DynamoDBWriteTransaction) (*StartDuaryRes, util.ApplicationError) {
	findMember, svcError := svc.memberSvc.GetMember(loginMember.SocialId, loginMember.Provider)
	if svcError != nil {
		return nil, svcError
	}

	// Find Couple
	findCouple, svcError := svc.coupleSvc.FindByCoupleCode(req.CoupleCode)
	if svcError != nil {
		return nil, svcError
	}

	findMember.CoupleId = findCouple.Id

	transaction.BeginTransaction()

	// Couple.Members 에 member 추가
	// Member.CoupleId 에 연결될 coupleId 를 넣어줌
	updateCoupleReq := &couple.UpdateCoupleReq{
		Members: append(findCouple.Members, findMember),
	}
	svcError = svc.coupleSvc.UpdateCouple(updateCoupleReq, transaction)
	if svcError != nil {
		return nil, svcError
	}

	// Update Member
	updateReq := &member.UpdateMemberReq{
		Provider: loginMember.Provider,
		SocialId: loginMember.SocialId,
		CoupleId: findCouple.Id,
	}
	updatedMember, svcError := svc.memberSvc.UpdateMember(updateReq, transaction)
	if svcError != nil {
		return nil, svcError
	}

	// execute transaction
	_, err := transaction.Execute()
	if err != nil {
		return nil, util.DBUpdateError{}
	}

	// generate new token with couple id
	jwtUtil := &jwtutil.Impl{}
	key := os.Getenv("secretKey")
	newToken := jwtUtil.NewToken(jwtUtil.GenerateSubject(findMember), findCouple.Id, key)

	result := &StartDuaryRes{
		Member: updatedMember,
		Couple: findCouple,
		Token:  newToken,
	}

	return result, nil
}

func (svc *ServiceImpl) CheckConnected(loginMember *auth.LoginMember) (*CheckConnectedRes, util.ApplicationError) {
	findMember, err := svc.memberSvc.GetMember(loginMember.SocialId, loginMember.Provider)
	if err != nil {
		return nil, err
	}
	if findMember.CoupleId == nil {
		return &CheckConnectedRes{}, nil
	}

	findCouple, err := svc.coupleSvc.FindById(*findMember.CoupleId)
	if err != nil {
		return nil, err
	}
	jwtUtil := &jwtutil.Impl{}
	key := os.Getenv("secretKey")
	newToken := jwtUtil.NewToken(jwtUtil.GenerateSubject(findMember), findCouple.Id, key)

	return &CheckConnectedRes{
		Token: newToken, Couple: findCouple,
	}, nil
}
