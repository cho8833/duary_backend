package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cho8833/duary_lambda/fcm"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/event"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/shared"
	"github.com/teambition/rrule-go"
	"log"
	"strconv"
	"time"
)

/*
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -tags lambda.norpc -o bootstrap api/send_anniversary_fcm/main/main.go && chmod 755 bootstrap && zip  build/package/send_anniversary_fcm.zip bootstrap duary-8c5b2-firebase-adminsdk-9d1a5-e86abbedfc.json && rm bootstrap
*/

func sendAnniversaryFCM(ctx context.Context, jsonMsg json.RawMessage) (events.APIGatewayProxyResponse, error) {
	fcmClient, err := fcm.GetFCMClient()
	if err != nil {
		log.Printf("failed to get fcmClient: %+v\n", err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	dynamodbClient, err := model.GetDynamoDBClient()
	if err != nil {
		log.Printf("failed to get DynamoDBClient: %+v\n", err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}
	memberRepo := member.NewRepository(dynamodbClient)
	memberSvc := member.NewService(memberRepo)

	fcmReq := &fcm.SendAnniversaryReq{}
	err = json.Unmarshal(jsonMsg, fcmReq)
	if err != nil {
		log.Printf("failed to get req body: %+v\n\n", err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	var content string
	var title string
	if fcmReq.EventData.EventType == event.Birthday {
		if fcmReq.EventData.Frequency == event.Yearly {
			recurCount, err := getYearlyRecurCount(fcmReq.EventData)
			if err != nil {
				return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
			}
			title = strconv.Itoa(*recurCount) + "주년"
		} else if fcmReq.EventData.Frequency == event.Daily {
			recurCount, err := get100DayRecurCount(fcmReq.EventData)
			if err != nil {
				return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
			}
			title = strconv.Itoa(*recurCount*100) + "일"
		}

		content = "내일은 우리가 만난 지 " + title + " 되는 날이에요!"
	} else if fcmReq.EventData.EventType == event.Anniversary {
		content = "내일은 연인의 생일이에요!"
		title = "연인 생일"
	}

	fcmService := fcm.NewService(fcmClient, memberSvc)

	sendReq := fcm.SendReq{
		Title:          title,
		Body:           content,
		TargetMemberId: fcmReq.TargetMemberId,
	}

	err = fcmService.Send(sendReq)
	if err != nil {
		log.Printf("failed to send fcm: %+v\n", err)
		return shared.LambdaAppErrorResponse(shared.InternalServerError{}), nil
	}

	return shared.LambdaResponseWithData(nil), nil
}

func getYearlyRecurCount(ev event.VO) (*int, error) {
	rule, err := rrule.NewRRule(
		rrule.ROption{
			Freq:    rrule.YEARLY,
			Dtstart: *ev.RecurStartDate,
		})
	if err != nil {
		log.Printf("failed to create rule: %+v\n", err)
		return nil, err
	}

	now := time.Now().UTC()

	allOccurrences := rule.Between(rule.GetDTStart(), now, true)

	recurCount := len(allOccurrences)

	return &recurCount, nil
}

func get100DayRecurCount(ev event.VO) (*int, error) {
	rule, err := rrule.NewRRule(
		rrule.ROption{
			Freq:     rrule.DAILY,
			Interval: 100,
			Dtstart:  *ev.RecurStartDate,
		})
	if err != nil {
		log.Printf("failed to create rule: %+v\n", err)
		return nil, err
	}

	now := time.Now().UTC()
	allOccurrences := rule.Between(rule.GetDTStart(), now, true)
	recurCount := len(allOccurrences)

	return &recurCount, nil

}

func main() {
	lambda.Start(sendAnniversaryFCM)
}
