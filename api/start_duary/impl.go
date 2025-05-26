package start_duary

import (
	"github.com/cho8833/duary_lambda/appjwt"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/shared"
	"log"
	"os"
	"time"
)

type StartDuaryReq struct {
	Name         *string    `json:"name"`
	Birthday     *time.Time `json:"birthday"`
	RelationDate *time.Time `json:"relationDate"`
	MyCharacter  *string    `json:"myCharacter"`
}

type StartDuaryRes struct {
	Member *member.Member         `json:"member"`
	Couple *couple.Couple         `json:"couple"`
	Token  *appjwt.ApplicationJWT `json:"token"`
}

func StartDuary(request *StartDuaryReq, transaction *model.DynamoDBWriteTransaction, coupleSvc couple.Service, memberSvc member.Service) (*StartDuaryRes, shared.ApplicationError) {
	// validate request
	if request.Name == nil || request.Birthday == nil || request.MyCharacter == nil || request.RelationDate == nil {
		log.Printf("invalid request. request %+v", request)
		return nil, shared.BadRequestError{}
	}

	authContext := shared.GetAuthContext()
	socialId := authContext.SocialId
	provider := authContext.Provider
	coupleId := authContext.CoupleId

	if coupleId != nil {
		log.Printf("couple already exists, coupleId: %s", *coupleId)
		return nil, shared.CoupleAlreadyExists{}
	}

	if socialId == nil || provider == nil {
		log.Printf("socialId or provider is nil. socialId: %v, provider: %v", socialId, provider)
		return nil, shared.BadRequestError{}
	}

	// begin transaction
	transaction.BeginTransaction()

	newCoupleId := coupleSvc.GenerateUID()

	// update member
	memberReq := &member.UpdateMemberReq{
		CoupleId:  newCoupleId,
		Name:      request.Name,
		Birthday:  request.Birthday,
		SocialId:  *socialId,
		Character: request.MyCharacter,
		Provider:  *provider,
	}
	updatedMember, svcErr := memberSvc.UpdateMember(memberReq, transaction)
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
	newCouple, err := coupleSvc.CreateCoupleWithId(*newCoupleId, coupleReq, transaction)
	if err != nil {
		return nil, err
	}

	// execute transaction
	_, transactionError := transaction.Execute()
	if transactionError != nil {
		return nil, shared.DBError{}
	}

	// generate new token: add coupleId in token
	jwtUtil := &appjwt.Impl{}
	key := os.Getenv("secretKey")
	appToken := jwtUtil.NewToken(jwtUtil.GenerateSubject(updatedMember.SocialId, updatedMember.Provider), updatedMember.CoupleId, key)

	res := &StartDuaryRes{
		Member: updatedMember,
		Couple: newCouple,
		Token:  appToken,
	}
	return res, nil
}
