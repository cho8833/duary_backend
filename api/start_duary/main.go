package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_backend/appjwt"
	"github.com/cho8833/duary_backend/model"
	"github.com/cho8833/duary_backend/model/couple"
	"github.com/cho8833/duary_backend/model/member"
	"github.com/cho8833/duary_backend/shared"
	"log"
	"os"
	"time"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/start_duary/main.go && chmod 755 bootstrap && zip  build/package/start_duary.zip bootstrap && rm bootstrap
*/

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

func startDuaryAPI(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// init
	dynamoDBClient, err := model.GetDynamoDBClient()
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	stage := req.StageVariables["stage"]

	transaction := model.NewWriteTransaction(dynamoDBClient)
	coupleRepo := couple.NewRepository(dynamoDBClient, stage)
	memberRepo := member.NewRepository(dynamoDBClient, stage)
	coupleSvc := couple.NewService(coupleRepo)
	memberSvc := member.NewService(memberRepo)

	initDuaryReq := &StartDuaryReq{}
	err = json.Unmarshal([]byte(req.Body), &initDuaryReq)
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	shared.NewAuthContext(req)

	result, svcError := StartDuary(initDuaryReq, transaction, coupleSvc, memberSvc)

	if svcError != nil {
		return shared.LambdaAppErrorResponse(svcError), nil
	}

	return shared.LambdaResponseWithDataAndHeader(result, appjwt.ApplicationJWTToHeader(*result.Token)), nil
}

func StartDuary(request *StartDuaryReq, transaction *model.DynamoDBWriteTransaction, coupleSvc couple.Service, memberSvc member.Service) (*StartDuaryRes, shared.ApplicationError) {
	// validate request
	// Birthday 는 nil 일 수 있음
	if request.Name == nil || request.MyCharacter == nil {
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
	updatedMember, svcErr := memberSvc.UpdateTransaction(memberReq, transaction)
	if svcErr != nil {
		return nil, svcErr
	}

	// create Couple
	coupleReq := &couple.CreateCoupleReq{
		RelationDate: request.RelationDate,
		Members: []member.Member{
			*updatedMember,
		},
		ConnectedMemberIds: []string{
			updatedMember.GetId(),
		},
	}
	newCouple, err := coupleSvc.CreateWithId(*newCoupleId, coupleReq, transaction)
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

func main() {
	lambda.Start(startDuaryAPI)
}
