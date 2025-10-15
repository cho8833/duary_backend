package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_backend/model"
	"github.com/cho8833/duary_backend/model/couple"
	"github.com/cho8833/duary_backend/model/event"
	"github.com/cho8833/duary_backend/model/member"
	"github.com/cho8833/duary_backend/scheduler"
	"github.com/cho8833/duary_backend/shared"
	"log"
	"time"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/update_member/main.go && chmod 755 bootstrap && zip  build/package/update_member.zip bootstrap && rm bootstrap
*/

type UpdateMemberReq struct {
	Birthday            *time.Time             `json:"birthday"`
	Name                *string                `json:"name"`
	Character           *string                `json:"character"`
	FcmToken            *string                `json:"fcmToken"`
	MyAlarm             *event.AlarmOffset     `json:"myAlarm"`
	LoverAlarm          *event.AlarmOffset     `json:"loverAlarm"`
	SyncedAppleCalendar []member.AppleCalendar `json:"syncedAppleCalendar"`
}

type UpdateMemberRes struct {
	Member *member.Member `json:"member"`
	Couple *couple.Couple `json:"couple"`
}

func updateMember(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

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

	memberRepo := member.NewRepository(dynamodbClient, stage)
	memberSvc := member.NewService(memberRepo)
	coupleRepo := couple.NewRepository(dynamodbClient, stage)
	coupleSvc := couple.NewService(coupleRepo)
	eventRepo := event.NewRepository(dynamodbClient, stage)
	eventSvc := event.NewService(eventRepo)

	updateReq := &UpdateMemberReq{}
	err = json.Unmarshal([]byte(request.Body), updateReq)
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	transaction := model.NewWriteTransaction(dynamodbClient)

	shared.NewAuthContext(request)

	result, svcErr := UpdateMember(updateReq, transaction, coupleSvc, memberSvc, eventSvc, *schedulerHelper, stage)

	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}
	return shared.LambdaResponseWithData(result), nil
}

func UpdateMember(req *UpdateMemberReq, transaction *model.DynamoDBWriteTransaction, coupleSvc couple.Service, memberSvc member.Service, eventSvc event.Service, schedulerHelper scheduler.BridgeSchedulerHelper, stage string) (*UpdateMemberRes, shared.ApplicationError) {

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
		SocialId:            *socialId,
		Provider:            *provider,
		Name:                req.Name,
		Birthday:            req.Birthday,
		Character:           req.Character,
		FcmToken:            req.FcmToken,
		MyAlarm:             req.MyAlarm,
		LoverAlarm:          req.LoverAlarm,
		SyncedAppleCalendar: req.SyncedAppleCalendar,
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
	// 커플이 연결된 상황에서만 Birthday Event 를 Update
	var birthday *event.VO
	if req.Birthday != nil && len(updatedCouple.Members) > 1 {
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
	// 커플이 연결된 상황에서만 Birthday Scheduler 를 Update
	if birthday != nil && len(updatedCouple.Members) > 1 {
		// 먼저 Schedule 을 찾은 후 있으면 Update, 없으면 Create
		_, err := schedulerHelper.GetEventSchedule(schedulerHelper.GetBirthdayScheduleName(*coupleId, updatedMember.GetId()))
		if err != nil {
			err = schedulerHelper.CreateAnniversarySchedule(*birthday, []member.Member{*updatedMember}, stage)
			if err != nil {
				log.Printf(err.Error())
			}
		} else {
			err = schedulerHelper.UpdateAnniversarySchedule(*birthday, []member.Member{*updatedMember}, stage)
			if err != nil {
				log.Printf(err.Error())
			}
		}
	}

	return &UpdateMemberRes{
		Member: updatedMember,
		Couple: updatedCouple,
	}, nil
}

// UpdateBirthday
// 기존 Birthday Event 를 삭제하고 새 Birthday Event 를 생성
// StartDateTime 이 바뀌기 때문에 SortKey 가 바뀌는 상황이고, Update 는 불가능하여 삭제하고 새로 생성
func UpdateBirthday(coupleId string, targetMember member.Member, eventSvc event.Service, transaction *model.DynamoDBWriteTransaction) (*event.VO, shared.ApplicationError) {
	// target member 의 Birthday 가 nil 이면 Birthday Event 도 생성되어 있지 않을 가능성이 높은데,
	// 삭제 대상이 없어도 Transaction 이 실패하지 않으니 따로 검사할 필요 없음
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

func main() {
	lambda.Start(updateMember)
}
