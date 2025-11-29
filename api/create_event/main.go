package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_backend/invoke_lambda"
	"github.com/cho8833/duary_backend/model"
	"github.com/cho8833/duary_backend/model/couple"
	"github.com/cho8833/duary_backend/model/event"
	"github.com/cho8833/duary_backend/model/ws_connection"
	"github.com/cho8833/duary_backend/shared"
	"github.com/cho8833/duary_backend/ws"
	"log"
	"strings"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap main.go && chmod 755 bootstrap && zip  ../../build/package/create_event.zip bootstrap && rm bootstrap
*/

func createEvent(context context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//----------------------------init----------------------------------------//
	dynamoDBClient, err := model.GetDynamoDBClient()
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	stage := req.StageVariables["stage"]

	// get auth
	authCtx := shared.NewAuthContext(req)
	coupleId := authCtx.CoupleId
	if coupleId == nil {
		log.Printf("coupleId is null")
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}
	socialId := authCtx.SocialId
	if socialId == nil {
		log.Printf("socialId is null")
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}
	provider := authCtx.Provider
	if provider == nil {
		log.Printf("provider is null")
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	eventRepo := event.NewRepository(dynamoDBClient, stage)
	coupleRepo := couple.NewRepository(dynamoDBClient, stage)
	coupleSvc := couple.NewService(coupleRepo)
	eventSvc := event.NewService(eventRepo)

	saveEventReq := &event.SaveReq{}
	err = json.Unmarshal([]byte(req.Body), &saveEventReq)
	if err != nil {
		log.Printf(err.Error())
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}
	loginMemberId := *socialId + "-" + *provider
	saveEventReq.CoupleId = *coupleId
	saveEventReq.CreatedBy = loginMemberId

	//----------------------------Save Event(DynamoDB)-----------------------------------//
	vo, svcError := eventSvc.Save(saveEventReq)
	if svcError != nil {
		return shared.LambdaAppErrorResponse(svcError), nil
	}

	//------------------if save completed, try sending WS, FCM to lover-----------------//
	// find couple
	foundCouple, svcErr := coupleSvc.FindById(*coupleId)
	if svcErr == nil {
		if len(foundCouple.Members) > 1 {
			var loverId string
			for _, coupleMember := range foundCouple.Members {
				memberId := coupleMember.GetId()
				if memberId != loginMemberId {
					loverId = memberId
				}
			}
			idSplit := strings.Split(loverId, "-")

			//----------------------------Invoke Send FCM to lover Lambda--------------------------------------//
			lambdaService := invoke_lambda.NewService(nil)
			err = lambdaService.SendEventFCM(context, stage, vo, loverId, invoke_lambda.EventCreated)
			if err != nil {
				log.Printf("failed to send fcm : " + err.Error())
			}

			//----------------------------Send WS Message to lover--------------------------------------//
			wsRepo := ws_connection.NewRepository(dynamoDBClient, stage)
			wsSvc, err := ws.NewService(&wsRepo, stage)
			if err == nil {
				err = wsSvc.Send(idSplit[0], idSplit[1], ws.EventCreated, vo)
				if err != nil {
					log.Printf("failed to send ws : " + err.Error())
				}
			} else {
				log.Printf("failed to create WS client : " + err.Error())
			}

		} else {
			log.Printf("couple is not connected, stop sending WS Message")
		}
	} else {
		log.Printf("faeild to find couple : " + svcErr.Error())
	}

	return shared.LambdaResponseWithData(vo), nil
}

func main() {
	lambda.Start(createEvent)
}
