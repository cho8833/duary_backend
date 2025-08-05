package ws

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/cho8833/duary_lambda/model/ws_connection"
	"log"
)

type Action string

const (
	CoupleConnected Action = "COUPLE_CONNECTED"
	EventDeleted    Action = "EVENT_DELETED"
	EventCreated    Action = "EVENT_CREATED"
	EventUpdated    Action = "EVENT_UPDATED"
	LoverUpdated    Action = "LOVER_UPDATED"
)

type Data[T any] struct {
	Action Action `json:"action"`
	Data   T      `json:"data"`
}

type Service interface {
	Send(socialId string, provider string, action Action, value interface{}) error
}

type ServiceImpl struct {
	repo        ws_connection.Repository
	apigwClient apigatewaymanagementapi.Client
}

func NewService(repo ws_connection.Repository) (*ServiceImpl, error) {
	apigwClient, err := GetApiGwClient()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &ServiceImpl{repo: repo, apigwClient: *apigwClient}, nil

}

func (svc *ServiceImpl) Send(socialId string, provider string, action Action, value interface{}) error {
	conn, err := svc.repo.FindBySocialIdAndProvider(socialId, provider)
	if err != nil {
		log.Println(err)
		return err
	}

	data := Data[any]{
		Action: action,
		Data:   value,
	}

	encoded, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return err
	}

	err = svc.publish(conn.ConnectionId, encoded)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (svc *ServiceImpl) publish(connId string, data []byte) error {
	_, err := svc.apigwClient.PostToConnection(context.TODO(), &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: &connId,
		Data:         data,
	})

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
