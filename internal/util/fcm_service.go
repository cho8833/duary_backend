package util

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"log"
)

type FCMReq struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	FcmTokens []string `json:"fcmTokens"`
}

type FCMService struct {
	client *messaging.Client
}

func FcmService(client *messaging.Client) *FCMService {
	return &FCMService{client: client}
}

func (s *FCMService) Send(req FCMReq) {
	notification := &messaging.Notification{
		Title: req.Title,
		Body:  req.Body,
	}
	if len(req.FcmTokens) == 1 {
		message := &messaging.Message{
			Data: map[string]string{
				"title": req.Title,
				"body":  req.Body,
			},
			Notification: notification,
			Token:        req.FcmTokens[0],
		}
		response, err := s.client.Send(context.Background(), message)
		if err != nil {
			log.Printf("failed to send: %+v\n", err)
		} else {
			log.Printf("send response: %+v\n\n", response)
		}
	} else {
		message := &messaging.MulticastMessage{
			Data: map[string]string{
				"title": req.Title,
				"body":  req.Body,
			},
			Notification: notification,
			Tokens:       req.FcmTokens,
		}
		response, err := s.client.SendEachForMulticast(context.Background(), message)
		if err != nil {
			log.Printf("failed to send: %+v\n", err)
		} else {
			log.Printf("send response: %+v\n\n", response)
		}
	}
}
