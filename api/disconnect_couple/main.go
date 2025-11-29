package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_backend/appjwt"
	"github.com/cho8833/duary_backend/model"
	"github.com/cho8833/duary_backend/model/couple"
	"github.com/cho8833/duary_backend/model/event"
	"github.com/cho8833/duary_backend/model/member"
	"github.com/cho8833/duary_backend/scheduler"
	"github.com/cho8833/duary_backend/shared"
	"log"
	"os"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap main.go && chmod 755 bootstrap && zip  ../../build/package/disconnect_couple.zip bootstrap && rm bootstrap
*/

type StartDuaryRes struct {
	Member *member.Member         `json:"member"`
	Couple *couple.Couple         `json:"couple"`
	Token  *appjwt.ApplicationJWT `json:"token"`
}

// StartDuary 전, SignIn 완료 후 상태로 돌아감
func handler(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	authContext := shared.NewAuthContext(request)

	coupleId := authContext.CoupleId
	if coupleId == nil {
		log.Printf("couple ID  is null")
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}
	socialId := authContext.SocialId
	if socialId == nil {
		log.Printf("social ID  is null")
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}
	provider := authContext.Provider
	if provider == nil {
		log.Printf("provider is null")
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

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
	coupleRepo := couple.NewRepository(dynamodbClient, stage)
	eventRepo := event.NewRepository(dynamodbClient, stage)

	memberSvc := member.NewService(memberRepo)
	coupleSvc := couple.NewService(coupleRepo)
	eventSvc := event.NewService(eventRepo)

	targetCouple, svcErr := coupleSvc.FindById(*coupleId)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	// begin transaction
	transaction := model.NewWriteTransaction(dynamodbClient)
	transaction.BeginTransaction()

	// remove coupleId of requested Member
	updatedMember, svcErr := memberSvc.RemoveCoupleIdTransaction(*socialId, *provider, transaction)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	// remove requested Member from couple
	var remainMember []member.Member
	for _, m := range targetCouple.Members {
		if m.SocialId != *socialId {
			remainMember = append(remainMember, m)
		}
	}
	updateCoupleReq := &couple.UpdateCoupleReq{
		Id:      coupleId,
		Members: remainMember,
	}
	_, svcErr = coupleSvc.Update(updateCoupleReq, transaction)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	// delete anniversary event
	svcErr = eventSvc.DeleteByCoupleIdAndTypeTransaction(*coupleId, event.Anniversary, transaction)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}
	svcErr = eventSvc.DeleteByCoupleIdAndTypeTransaction(*coupleId, event.Birthday, transaction)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	_, err = transaction.Execute()
	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	// delete anniversary notification schedule
	svcErr = schedulerHelper.DeleteEventSchedule(schedulerHelper.Get100DayAnniversaryScheduleName(*coupleId))
	if svcErr != nil {
		log.Printf(svcErr.Error())
	}
	svcErr = schedulerHelper.DeleteEventSchedule(schedulerHelper.GetYearlyAnniversaryScheduleName(*coupleId))
	if svcErr != nil {
		log.Printf(svcErr.Error())
	}
	for _, mId := range targetCouple.ConnectedMemberIds {
		svcErr = schedulerHelper.DeleteEventSchedule(schedulerHelper.GetBirthdayScheduleName(*coupleId, mId))
		if svcErr != nil {
			log.Printf(svcErr.Error())
		}
	}

	jwtUtil := &appjwt.Impl{}
	key := os.Getenv("secretKey")
	newToken := jwtUtil.NewToken(jwtUtil.GenerateSubject(updatedMember.SocialId, updatedMember.Provider), nil, key)

	result := &StartDuaryRes{
		Token:  newToken,
		Member: updatedMember,
		Couple: nil,
	}

	return shared.LambdaResponseWithDataAndHeader(result, appjwt.ApplicationJWTToHeader(*result.Token)), nil
}

func main() {
	lambda.Start(handler)
}
