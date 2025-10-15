package main

import (
	"context"
	"encoding/json"
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
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/update_event/main.go && chmod 755 bootstrap && zip  build/package/update_event.zip bootstrap && rm bootstrap
*/
func updateEvent(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// init
	dynamodbClient, err := model.GetDynamoDBClient()
	if err != nil {
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	eventRepo := event.NewRepository(dynamodbClient)
	eventSvc := event.NewService(eventRepo)
	coupleRepo := couple.NewRepository(dynamodbClient)
	coupleSvc := couple.NewService(coupleRepo)

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

	// parse req
	updateEventReq := &event.EditReq{}
	err = json.Unmarshal([]byte(request.Body), updateEventReq)
	if err != nil {
		log.Println(err.Error())
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}
	eventId := request.QueryStringParameters["id"]

	// update event (dynamodb)
	vo, svcError := eventSvc.Update(*coupleId, eventId, updateEventReq)
	if svcError != nil {
		return shared.LambdaAppErrorResponse(shared.BadRequestError{}), nil
	}

	//----------------------------Send WS Message to lover-------------------------------//
	wsRepo := ws_connection.NewRepository(dynamodbClient)
	wsSvc, err := ws.NewService(&wsRepo)
	if err == nil {
		loginMemberId := *socialId + "-" + *provider
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
				err = wsSvc.Send(idSplit[0], idSplit[1], ws.EventUpdated, vo)
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
	return shared.LambdaResponseWithData(vo), nil
}

func main() {
	lambda.Start(updateEvent)
}
