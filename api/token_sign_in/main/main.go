package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/appjwt"
	"github.com/cho8833/duary_lambda/auth"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/shared"
	"log"
	"os"
)

/*
GCO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/token_sign_in/main/main.go && chmod 755 bootstrap && zip build/package/token_sign_in.zip bootstrap && rm bootstrap
*/
func tokenSignIn(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	dynamodbClient, err := model.GetDynamoDBClient()
	if err != nil {
		log.Printf("failed to get DynamoDB client: %s", err.Error())
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	memberRepo := member.NewMemberRepository(dynamodbClient)
	coupleRepo := couple.NewCoupleRepository(dynamodbClient)

	memberSvc := member.NewMemberService(memberRepo)
	coupleSvc := couple.NewCoupleService(coupleRepo)

	authContext := shared.NewAuthContext(req)

	signInReq := &auth.SignInReq{}
	err = json.Unmarshal([]byte(req.Body), &signInReq)
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	// update member.fcmToken
	updatedMember, svcErr := memberSvc.UpdateFcmToken(*authContext.SocialId, *authContext.Provider, signInReq.FcmToken)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	// find Couple of member
	var memberCouple *couple.Couple
	if updatedMember.CoupleId != nil {
		memberCouple, svcErr = coupleSvc.FindById(*updatedMember.CoupleId)
		if temp := new(shared.UserNotFound); !errors.As(err, &temp) && svcErr != nil {
			log.Printf("failed to find member. Error: %s", svcErr)
			return shared.LambdaAppErrorResponse(svcErr), nil
		}
	}

	jwtUtil := appjwt.Impl{}
	key := os.Getenv("secretKey")

	newToken := jwtUtil.NewToken(jwtUtil.GenerateSubject(updatedMember.SocialId, updatedMember.Provider), updatedMember.CoupleId, key)

	res := &auth.SignInRes{
		IsRegister: false,
		Member:     updatedMember,
		Couple:     memberCouple,
		Token:      newToken,
	}

	return shared.LambdaResponseWithData(res), nil

}

func main() {
	lambda.Start(tokenSignIn)
}
