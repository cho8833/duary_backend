package event_schedule

import (
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/shared"
	"log"
)

type Service interface {
	UpdateTransaction(memberId string, eventId string, req UpdateEventScheduleReq, transaction *model.DynamoDBWriteTransaction) shared.ApplicationError
}

type ServiceImpl struct {
	repository Repository
}

func NewService(repository Repository) *ServiceImpl {
	return &ServiceImpl{repository: repository}
}

func (s *ServiceImpl) UpdateTransaction(memberId string, eventId string, req UpdateEventScheduleReq, transaction *model.DynamoDBWriteTransaction) shared.ApplicationError {
	err := s.repository.UpdateTransaction(memberId, eventId, req, transaction)
	if err != nil {
		log.Println(err)
		return shared.DBUpdateError{}
	}
	return nil
}
