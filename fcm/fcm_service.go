package fcm

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"log"
)

func GetFCMClient() (*messaging.Client, error) {

	ctx := context.Background()

	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to create Firebase app: %v", err)
		return nil, err
	}
	log.Printf("Using Firebase Cloud Messaging client: %+v", app)
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Printf("Error creating Messaging client: %+v", err)
		return nil, err
	}
	return client, nil
}

type SendReq struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	FcmTokens []string `json:"fcmTokens"`
}

type Service struct {
	client *messaging.Client
}

func NewService(client *messaging.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Send(req SendReq) {
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
