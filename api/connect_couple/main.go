package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_backend/appjwt"
	"github.com/cho8833/duary_backend/model"
	"github.com/cho8833/duary_backend/model/couple"
	"github.com/cho8833/duary_backend/model/event"
	"github.com/cho8833/duary_backend/model/member"
	"github.com/cho8833/duary_backend/model/ws_connection"
	"github.com/cho8833/duary_backend/scheduler"
	"github.com/cho8833/duary_backend/shared"
	"github.com/cho8833/duary_backend/ws"
	"log"
	"os"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/connect_couple/main.go && chmod 755 bootstrap && zip  build/package/connect_couple.zip bootstrap && rm bootstrap
*/

type ConnectCoupleReq struct {
	CoupleCode *string `json:"coupleCode"`
}

type StartDuaryRes struct {
	Member *member.Member         `json:"member"`
	Couple *couple.Couple         `json:"couple"`
	Token  *appjwt.ApplicationJWT `json:"token"`
}

func handler(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamodbClient, err := model.GetDynamoDBClient()
	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	schedulerClient, err := scheduler.GetSchedulerClient()
	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	schedulerHelper := scheduler.NewEventBridgeSchedulerHelper(schedulerClient)

	stage := request.StageVariables["stage"]

	transaction := model.NewWriteTransaction(dynamodbClient)
	memberRepo := member.NewRepository(dynamodbClient, stage)
	coupleRepo := couple.NewRepository(dynamodbClient, stage)
	eventRepo := event.NewRepository(dynamodbClient, stage)
	wsRepo := ws_connection.NewRepository(dynamodbClient, stage)

	memberSvc := member.NewService(memberRepo)
	coupleSvc := couple.NewService(coupleRepo)
	eventSvc := event.NewService(eventRepo)
	wsSvc, err := ws.NewService(&wsRepo, stage)
	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	connectCoupleReq := &ConnectCoupleReq{}
	err = json.Unmarshal([]byte(request.Body), &connectCoupleReq)
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	shared.NewAuthContext(request)

	res, svcError := ConnectCouple(connectCoupleReq, transaction, coupleSvc, memberSvc, eventSvc, *schedulerHelper, wsSvc, stage)
	if svcError != nil {
		return shared.LambdaAppErrorResponse(svcError), nil
	}

	return shared.LambdaResponseWithDataAndHeader(res, appjwt.ApplicationJWTToHeader(*res.Token)), nil
}

func ConnectCouple(req *ConnectCoupleReq, transaction *model.DynamoDBWriteTransaction, coupleSvc couple.Service, memberSvc member.Service, eventSvc event.Service, schedulerHelper scheduler.BridgeSchedulerHelper, wsSvc ws.Service, stage string) (*StartDuaryRes, shared.ApplicationError) {

	// *********** Validate ***************
	if req.CoupleCode == nil || len(*req.CoupleCode) == 0 {
		return nil, shared.BadRequestError{}
	}

	authContext := shared.GetAuthContext()

	loginMemberId := *authContext.SocialId + "-" + *authContext.Provider

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

	// target Couple 이 이미 연결된 적이 있는 경우, target Couple 의 연결된 적 있는 Member 가 아니면 Bad Request
	if len(targetCouple.ConnectedMemberIds) > 1 {
		hasConnectedMember := false
		for _, memberId := range targetCouple.ConnectedMemberIds {
			if memberId == loginMemberId {
				hasConnectedMember = true
			}
		}
		if !hasConnectedMember {
			log.Printf("couple has already connected")
			return nil, shared.BadRequestError{}
		}
	}

	// *********** Start DB Transaction ***************
	transaction.BeginTransaction()

	// 기존의 Couple 삭제 : 커플 연결 요청한 Member 의 Couple 삭제
	svcErr := coupleSvc.Delete(*authContext.CoupleId, transaction)
	if svcErr != nil {
		return nil, shared.DBUpdateError{}
	}

	// Update Member : 커플 연결 요청한 Member 의 Couple Id 변경
	updateReq := &member.UpdateMemberReq{
		Provider: *authContext.Provider,
		SocialId: *authContext.SocialId,
		CoupleId: targetCouple.Id,
	}
	updatedMember, svcError := memberSvc.UpdateTransaction(updateReq, transaction)
	if svcError != nil {
		return nil, svcError
	}

	// Update Couple
	connectedIds := targetCouple.ConnectedMemberIds
	if len(targetCouple.ConnectedMemberIds) < 2 {
		connectedIds = append(targetCouple.ConnectedMemberIds, updatedMember.GetId()) // couple 의 Member 연결 기록에 추가
	}

	updateCoupleReq := &couple.UpdateCoupleReq{
		Id:                 targetCouple.Id,
		Members:            append(targetCouple.Members, *updatedMember), // Member 추가
		ConnectedMemberIds: connectedIds,
	}
	updatedCouple, svcError := coupleSvc.Update(updateCoupleReq, transaction)
	if svcError != nil {
		return nil, svcError
	}

	// Save Couple Anniversary Events
	var day100VO *event.VO
	var yearlyVO *event.VO
	if updatedCouple.RelationDate != nil {
		// Create Anniversary Event
		firstMetDayReq := eventSvc.GenerateFirstMetDay(*updatedCouple.Id, updatedCouple.Members[0].GetId(), *updatedCouple.RelationDate)

		anniversary100DayReq := eventSvc.Generate100Anniversary(*updatedCouple.Id, updatedCouple.Members[0].GetId(), *updatedCouple.RelationDate)
		yearlyAnniversaryReq := eventSvc.GenerateYearlyAnniversary(*updatedCouple.Id, updatedCouple.Members[0].GetId(), *updatedCouple.RelationDate)

		_, svcErr = eventSvc.SaveTransaction(firstMetDayReq, transaction)
		if svcErr != nil {
			return nil, svcErr
		}
		day100VO, svcErr = eventSvc.SaveTransaction(anniversary100DayReq, transaction)
		if svcErr != nil {
			return nil, svcErr
		}
		yearlyVO, svcErr = eventSvc.SaveTransaction(yearlyAnniversaryReq, transaction)
		if svcErr != nil {
			return nil, svcErr
		}
	}

	birthdaySaveReqs := make(map[string]event.SaveReq)
	for _, m := range updatedCouple.Members {
		if m.Birthday != nil { // Member 의 Birthday 는 nil 일 수 있음
			birthdaySaveReqs[m.GetId()] = *eventSvc.GenerateBirthday(*updatedCouple.Id, m.GetId(), *m.Birthday)
		}
	}
	birthdays := make(map[string]event.VO)
	for m, rq := range birthdaySaveReqs {
		birthday, svcErr := eventSvc.SaveTransaction(&rq, transaction)
		if svcErr != nil {
			return nil, svcErr
		}
		birthdays[m] = *birthday
	}

	// execute transaction
	_, err := transaction.Execute()
	if err != nil {
		log.Println(err)
		return nil, shared.DBUpdateError{}
	}

	// 기념일 알림 schedule 생성
	// TODO: 오류 처리 필요. DB 작업이 성공해도 schedule 작업은 실패할 수 있음
	if day100VO != nil {
		_ = schedulerHelper.CreateAnniversarySchedule(*day100VO, updatedCouple.Members, stage)
	}
	if yearlyVO != nil {
		_ = schedulerHelper.CreateAnniversarySchedule(*yearlyVO, updatedCouple.Members, stage)
	}

	for mid, birthday := range birthdays {
		for _, m := range updatedCouple.Members {
			if m.GetId() == mid {
				_ = schedulerHelper.CreateAnniversarySchedule(birthday, []member.Member{m}, stage)
			}
		}
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

	// 상대방에게 couple connected notify
	var otherMember *member.Member
	for _, coupleMember := range targetCouple.Members {
		memberId := coupleMember.GetId()
		if memberId != loginMemberId {
			otherMember = &coupleMember
		}
	}
	if otherMember != nil {
		err = wsSvc.Send(otherMember.SocialId, otherMember.Provider, ws.CoupleConnected, result)
		if err != nil {
			log.Println("lover is not connected to ws", err)
		}
	}

	return result, nil
}

func main() {
	lambda.Start(handler)
}
