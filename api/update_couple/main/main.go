package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/couple"
	"github.com/cho8833/duary_lambda/shared"
	"log"
)

/*
GCO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/update_couple/main/main.go && chmod 755 bootstrap && zip build/package/update_couple.zip bootstrap && rm bootstrap
*/

func updateCouple(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	authContext := shared.NewAuthContext(request)

	client, err := model.GetDynamoDBClient()

	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	repo := couple.NewCoupleRepository(client)

	svc := couple.NewCoupleService(repo)

	req := &couple.UpdateCoupleReq{}

	err = json.Unmarshal([]byte(request.Body), req)

	if err != nil {
		log.Println(err)
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	req.Id = authContext.CoupleId

	res, svcErr := svc.UpdateCouple(req, nil)

	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	return shared.LambdaResponseWithData(res), nil

}
func main() {
	lambda.Start(updateCouple)
}
