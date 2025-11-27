package invoke_lambda

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/cho8833/duary_backend/model/event"
	"log"
	"os"
	"time"
)

type Service interface {
	SendEventFCM(context context.Context, stage string, vo *event.VO, targetMemberId string, action EventAction) error
}

type ServiceImpl struct {
	client lambda.Client
}

// NewService returns invoke lambda service
// if client is nil, create new default lambda client
func NewService(client *lambda.Client) Service {
	if client == nil {
		newClient, err := GetLambdaClient()
		if err != nil {
			log.Fatal(err)
			return nil
		}
		client = newClient
	}
	return &ServiceImpl{client: *client}
}

func (s *ServiceImpl) SendEventFCM(context context.Context, stage string, vo *event.VO, targetMemberId string, action EventAction) error {
	payload, err := json.Marshal(&SendEventReq{
		Title:          vo.Title,
		TargetMemberId: []string{targetMemberId},
		Body:           s.formatTime(vo.StartDateTime),
		Stage:          stage,
	})

	if err != nil {
		log.Printf("SendEventFCM - json marshal err: %v", err)
		return err
	}
	invokeInput := lambda.InvokeInput{
		FunctionName: aws.String(os.Getenv("send_fcm")),
		Payload:      payload,
		LogType:      types.LogTypeTail,
	}
	invokeOutput, err := s.client.Invoke(context, &invokeInput)
	if err != nil {
		log.Printf("SendEventFCM - invoke err: %v", err)
		return err
	}

	log.Printf("SendEventFCM - invokeOutput: %v", string(invokeOutput.Payload))
	return nil
}

func (s *ServiceImpl) formatTime(date time.Time) string {
	return date.Format("2006년 01월 02일 15시 04분")
}
