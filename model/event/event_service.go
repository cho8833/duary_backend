package event

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/shared"
	uuid2 "github.com/google/uuid"
	"log"
	"time"
)

type Service interface {
	Save(req *SaveReq) (*VO, shared.ApplicationError)
	GetBetweenStartAndEndDate(coupleId string, rangeStartDate time.Time, rangeEndDate time.Time) ([]VO, shared.ApplicationError)
	Update(coupleId string, id string, req *EditReq) (*VO, shared.ApplicationError)
	SaveTransaction(req *SaveReq, transaction *model.DynamoDBWriteTransaction) (*VO, shared.ApplicationError)
	GenerateOccurrence(vo VO, rangeStartDate time.Time, rangeEndDate time.Time) ([]VO, shared.ApplicationError)
	DeleteEvent(coupleId string, id string) (*VO, shared.ApplicationError)
	DeleteBirthdayTransaction(coupleId string, memberId string, transaction *model.DynamoDBWriteTransaction) shared.ApplicationError
	DeleteByCoupleIdAndTypeTransaction(coupleId string, eventType EventType, transaction *model.DynamoDBWriteTransaction) shared.ApplicationError
	GenerateFirstMetDay(coupleId string, createdBy string, relationDate time.Time) *SaveReq
	GenerateYearlyAnniversary(coupleId string, createdById string, relationDate time.Time) *SaveReq
	Generate100Anniversary(coupleId string, createdById string, relationDate time.Time) *SaveReq
	GenerateBirthday(coupleId string, memberId string, birthday time.Time) *SaveReq
}

type ServiceImpl struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &ServiceImpl{repository: repository}
}

func (service *ServiceImpl) Save(req *SaveReq) (*VO, shared.ApplicationError) {
	svcErr := req.Validate()
	if svcErr != nil {
		return nil, svcErr
	}
	event := FromReq(req, *service.generateUID())

	vo, err := service.repository.Save(event)
	if err != nil {
		log.Printf("save event error\nreq: %+v\nerror: %s", req, err)
		return nil, shared.DBSaveError{}
	}

	return vo, nil
}

func (service *ServiceImpl) SaveTransaction(req *SaveReq, transaction *model.DynamoDBWriteTransaction) (*VO, shared.ApplicationError) {
	svcErr := req.Validate()
	if svcErr != nil {
		return nil, svcErr
	}
	event := FromReq(req, *service.generateUID())
	vo, err := service.repository.SaveTransaction(event, transaction)
	if err != nil {
		log.Printf("save transaction error\nreq: %+v\nerror: %s", req, err)
		return nil, shared.DBSaveError{}
	}

	return vo, nil
}

func (service *ServiceImpl) Update(coupleId string, id string, req *EditReq) (*VO, shared.ApplicationError) {
	svcErr := req.Validate()
	if svcErr != nil {
		return nil, svcErr
	}
	vo, err := service.repository.Update(coupleId, id, *req)

	if err != nil {
		log.Printf("update event error\nreq: %+v\nerror: %s", req, err)
		return nil, shared.DBSaveError{}
	}

	return vo, nil
}

func (service *ServiceImpl) DeleteEvent(coupleId string, id string) (*VO, shared.ApplicationError) {

	authContext := shared.GetAuthContext()

	if coupleId != *authContext.CoupleId {
		log.Printf("unauthorized request : coupleId not matched\n user coupleId: %s, request coupleId: %s", *authContext.CoupleId, coupleId)
		return nil, shared.BadRequestError{}
	}

	deleted, err := service.repository.Delete(coupleId, id)

	if err != nil {
		return nil, shared.DBDeleteError{}
	}
	return deleted, nil
}

func (service *ServiceImpl) DeleteBirthdayTransaction(coupleId string, memberId string, transaction *model.DynamoDBWriteTransaction) shared.ApplicationError {
	events, err := service.repository.QueryByCoupleIdAndType(coupleId, Birthday, true)
	if err != nil {
		log.Printf("query event error")
		return shared.DBDeleteError{}
	}
	for _, ev := range events {
		if ev.CreatedBy == memberId {
			err = service.repository.DeleteTransaction(ev.CoupleId, ev.Id, transaction)
			if err != nil {
				log.Printf("delete event transaction error\n error: %s", err.Error())
				return shared.DBDeleteError{}
			}
		}
	}
	return nil
}

func (service *ServiceImpl) DeleteByCoupleIdAndTypeTransaction(coupleId string, eventType EventType, transaction *model.DynamoDBWriteTransaction) shared.ApplicationError {
	events, err := service.repository.QueryByCoupleIdAndType(coupleId, eventType, false)
	if err != nil {
		log.Printf("query event error")
		return shared.DBDeleteError{}
	}

	for _, ev := range events {
		err = service.repository.DeleteTransaction(ev.CoupleId, ev.Id, transaction)
		if err != nil {
			log.Printf("delete event transaction error\nreq: %+v\nerror: %s", ev, err)
			return shared.DBDeleteError{}
		}
	}
	return nil
}

func (service *ServiceImpl) GetBetweenStartAndEndDate(coupleId string, rangeStartDate time.Time, rangeEndDate time.Time) ([]VO, shared.ApplicationError) {

	candidateEvents, err := service.repository.FindByCoupleIdAndStartDateBefore(coupleId, rangeEndDate)
	if err != nil {
		// 데이터가 없는 건 exception 이 아님
		if temp := new(types.ResourceNotFoundException); errors.As(err, &temp) {
			return []VO{}, nil
		}
		log.Printf("failed to find events by coupleId and rangeStartDate before: %s", err)
		return nil, shared.DBReadError{}
	}

	var result []VO

	for _, event := range candidateEvents {
		// 반복 이벤트인 경우 이벤트 인스턴스(Occurrence) 생성
		if event.Frequency != OneTime {
			occurrences, err := service.GenerateOccurrence(event, rangeStartDate, rangeEndDate)
			if err != nil {
				return nil, err
			}
			result = append(result, occurrences...)

			// 단발성 이벤트인 경우 이베트가 range 내에 포함되는지 검사
		} else {
			startDateTime := event.StartDateTime
			endDateTime := event.EndDateTime
			if rangeStartDate.Before(endDateTime) && rangeEndDate.After(startDateTime) {
				result = append(result, event)
			}
		}
	}

	if len(result) == 0 {
		result = make([]VO, 0)
	}
	return result, nil
}

func (service *ServiceImpl) generateUID() *string {
	uuid := uuid2.New().String()[:32]
	return &uuid
}
