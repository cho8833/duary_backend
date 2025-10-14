package fcm

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/cho8833/duary_backend/model/event"
	"github.com/cho8833/duary_backend/model/member"
	"log"
	"strings"
)

func GetFCMClient() (*messaging.Client, error) {

	ctx := context.Background()

	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to create Firebase app: %v", err)
		return nil, err
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Printf("Error creating Messaging client: %+v", err)
		return nil, err
	}
	return client, nil
}

type SendReq struct {
	Title          string   `json:"title"`
	Body           string   `json:"body"`
	TargetMemberId []string `json:"target_member_id"`
}

type SendAnniversaryReq struct {
	EventData      event.VO
	TargetMemberId []string
}

type Service struct {
	client        *messaging.Client
	memberService member.Service
}

func NewService(client *messaging.Client, memberService member.Service) *Service {
	return &Service{client: client, memberService: memberService}
}

func (s *Service) Send(req SendReq) error {
	notification := &messaging.Notification{
		Title: req.Title,
		Body:  req.Body,
	}

	var tokens []string

	for _, mId := range req.TargetMemberId {
		temp := strings.Split(mId, "-")
		targetMember, svcErr := s.memberService.FindById(temp[0], temp[1])
		if svcErr != nil {
			log.Printf("Error finding target member: %v", svcErr)
		}
		if targetMember.FcmToken == nil {
			log.Printf("target member has no fcm token")
		}
		tokens = append(tokens, *targetMember.FcmToken)
	}

	message := &messaging.MulticastMessage{
		Data: map[string]string{
			"title": req.Title,
			"body":  req.Body,
		},
		Notification: notification,
		Tokens:       tokens,
	}

	response, err := s.client.SendEachForMulticast(context.Background(), message)
	if err != nil {
		log.Printf("failed to send: %+v\n", err)
		return err
	} else {
		log.Printf("send response: %+v\n\n", response)
		return nil
	}
}
