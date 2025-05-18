package common

import (
	"github.com/aws/smithy-go/ptr"
	"github.com/cho8833/duary_lambda/internal/auth"
	"github.com/cho8833/duary_lambda/internal/auth/jwtutil"
	"github.com/cho8833/duary_lambda/internal/couple"
	"github.com/cho8833/duary_lambda/internal/member"
	"github.com/cho8833/duary_lambda/internal/util"
	"log"
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
}

type ServiceImpl struct {
	memberSvc member.Service
	coupleSvc couple.Service
}

func NewCommonService(memberSvc member.Service, coupleSvc couple.Service) *ServiceImpl {
	return &ServiceImpl{memberSvc: memberSvc, coupleSvc: coupleSvc}
}

func (svc *ServiceImpl) StartDuary(request *StartDuaryReq, transaction *util.DynamoDBWriteTransaction) (*StartDuaryRes, util.ApplicationError) {
	// validate request
	if request.Name == nil || request.Birthday == nil || request.MyCharacter == nil || request.RelationDate == nil {
		log.Printf("invalid request. request %+v", request)
		return nil, util.BadRequestError{}
	}

	authContext := util.GetAuthContext()
	socialId := authContext.SocialId
	provider := authContext.Provider
	coupleId := authContext.CoupleId

	if coupleId != nil {
		log.Printf("couple already exists, coupleId: %s", *coupleId)
		return nil, util.CoupleAlreadyExists{}
	}

	if socialId == nil || provider == nil {
		log.Printf("socialId or provider is nil. socialId: %v, provider: %v", socialId, provider)
		return nil, util.BadRequestError{}
	}

	// begin transaction
	transaction.BeginTransaction()

	newCoupleId := svc.coupleSvc.GenerateUID()

	// update member
	memberReq := &member.UpdateMemberReq{
		CoupleId:  newCoupleId,
		Name:      request.Name,
		Birthday:  request.Birthday,
		SocialId:  *socialId,
		Character: request.MyCharacter,
		Provider:  *provider,
	}
	updatedMember, svcErr := svc.memberSvc.UpdateMember(memberReq, transaction)
	if svcErr != nil {
		return nil, svcErr
	}

	// create Couple
	coupleReq := &couple.CreateCoupleReq{
		RelationDate: *request.RelationDate,
		Members: []*member.Member{
			updatedMember,
		},
	}
	newCouple, err := svc.coupleSvc.CreateCoupleWithId(*newCoupleId, coupleReq, transaction)
	if err != nil {
		return nil, err
	}

	// execute transaction
	_, transactionError := transaction.Execute()
	if transactionError != nil {
		return nil, util.DBError{}
	}

	// generate new token: add coupleId in token
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
	if req.CoupleCode == nil || len(*req.CoupleCode) == 0 {
		return nil, util.BadRequestError{}
	}

	authContext := util.GetAuthContext()

	// findMember.coupleId == nil 은 startDuary 를 거치지 않았다는 의미 => bad request
	if authContext.CoupleId == nil {
		return nil, util.BadRequestError{}
	}

	// Find Couple
	targetCouple, svcError := svc.coupleSvc.FindByCoupleCode(req.CoupleCode)
	if svcError != nil {
		return nil, util.CoupleCodeNotFound{}
	}

	// 자신의 coupleId 와 Code 의 CoupleId 가 같음 -> 자신의 couple 에 연결한다 -> bad request
	if *authContext.CoupleId == *targetCouple.Id {
		return nil, util.BadRequestError{}
	}

	// Couple 이 이미 연결되어 있는 경우 Bad Request
	if len(targetCouple.Members) >= 2 {
		return nil, util.BadRequestError{}
	}

	transaction.BeginTransaction()

	// 기존의 Couple 삭제
	svcErr := svc.coupleSvc.DeleteCouple(*authContext.CoupleId, transaction)
	if svcErr != nil {
		return nil, util.DBUpdateError{}
	}

	// Update Member
	updateReq := &member.UpdateMemberReq{
		Provider: loginMember.Provider,
		SocialId: loginMember.SocialId,
		CoupleId: targetCouple.Id,
	}
	updatedMember, svcError := svc.memberSvc.UpdateMember(updateReq, transaction)
	if svcError != nil {
		return nil, svcError
	}

	// Update Couple
	updateCoupleReq := &couple.UpdateCoupleReq{
		Id:      targetCouple.Id,
		Members: append(targetCouple.Members, updatedMember), // Member 추가
		Code:    ptr.String(""),                              // 코드 삭제
	}
	updatedCouple, svcError := svc.coupleSvc.UpdateCouple(updateCoupleReq, transaction)
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
	newToken := jwtUtil.NewToken(jwtUtil.GenerateSubject(updatedMember), updatedCouple.Id, key)

	result := &StartDuaryRes{
		Member: updatedMember,
		Couple: updatedCouple,
		Token:  newToken,
	}

	return result, nil
}
