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
)

type Service interface {
	SendEventFCM(context context.Context, vo *event.VO) error
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

func (s *ServiceImpl) SendEventFCM(context context.Context, vo *event.VO) error {
	payload, err := json.Marshal(vo)
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
