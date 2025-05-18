package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/internal/common"
	"github.com/cho8833/duary_lambda/internal/couple"
	"github.com/cho8833/duary_lambda/internal/member"
	"github.com/cho8833/duary_lambda/internal/util"
	"log"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap internal/common/main/update_member_api.go && chmod 755 bootstrap && zip  build/package/common/update_member_api.zip bootstrap && rm bootstrap
*/
func updateMember(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamodbClient, err := util.GetCacheClient().GetDynamoDBClient()
	if err != nil {
		log.Printf(err.Error())
		return util.LambdaAppErrorResponse(util.InternalServerError{}), nil
	}

	memberRepo := member.NewMemberRepository(dynamodbClient)
	memberSvc := member.NewMemberService(memberRepo)
	coupleRepo := couple.NewCoupleRepository(dynamodbClient)
	coupleSvc := couple.NewCoupleService(coupleRepo)

	commonMemberSvc := common.NewMemberService(memberSvc, coupleSvc)

	updateReq := &common.UpdateMemberReq{}
	err = json.Unmarshal([]byte(request.Body), updateReq)
	if err != nil {
		log.Printf(err.Error())
		return util.LambdaAppErrorResponse(util.BadRequestError{}), nil
	}

	authContext := util.NewAuthContext(request)

	transaction := util.NewWriteTransaction(dynamodbClient)
	result, svcErr := commonMemberSvc.UpdateMember(updateReq, *authContext.SocialId, *authContext.Provider, *authContext.CoupleId, transaction)

	if svcErr != nil {
		return util.LambdaAppErrorResponse(svcErr), nil
	}
	return util.LambdaResponseWithData(result), nil
}

func main() {
	lambda.Start(updateMember)
}
