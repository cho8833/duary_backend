package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap main.go && chmod 755 bootstrap && zip  ../../build/package/delete_event.zip bootstrap && rm bootstrap
*/
func deleteEvent(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// get auth
	authCtx := shared.NewAuthContext(request)
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

	dynamodbClient, err := model.GetDynamoDBClient()
	if err != nil {
		return shared.LambdaAppErrorResponse(shared.DBError{}), nil
	}
	stage := request.StageVariables["stage"]

	eventRepo := event.NewRepository(dynamodbClient, stage)
	coupleRepo := couple.NewRepository(dynamodbClient, stage)
	coupleSvc := couple.NewService(coupleRepo)
	eventSvc := event.NewService(eventRepo)

	eventId := request.QueryStringParameters["id"]

	authContext := shared.NewAuthContext(request)

	if eventId == "" {
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	//----------------------------Delete Event(DynamoDB)-------------------------------//
	deleted, svcErr := eventSvc.DeleteEvent(*coupleId, eventId)
	if svcErr != nil {
		return shared.LambdaAppErrorResponse(svcErr), nil
	}

	//----------------------------Send WS Message to lover-------------------------------//
	wsRepo := ws_connection.NewRepository(dynamodbClient, stage)
	wsSvc, err := ws.NewService(&wsRepo, stage)
	if err == nil {
		loginMemberId := *authContext.SocialId + "-" + *authContext.Provider
		foundCouple, svcErr := coupleSvc.FindById(*coupleId)
		if svcErr == nil {
			if len(foundCouple.Members) > 1 {
				var otherMemberId string
				for _, coupleMember := range foundCouple.Members {
					memberId := coupleMember.GetId()
					if memberId != loginMemberId {
						otherMemberId = memberId
					}
				}
				idSplit := strings.Split(otherMemberId, "-")
				err = wsSvc.Send(idSplit[0], idSplit[1], ws.EventDeleted, deleted)
				if err != nil {
					log.Println("lover is not connected to ws", err)
				}
			} else {
				log.Printf("couple is not connected, stop sending WS Message")
			}
		}
	} else {
		log.Printf(err.Error())
	}

	return shared.LambdaResponseWithData(deleted), nil
}

func main() {
	lambda.Start(deleteEvent)
}
