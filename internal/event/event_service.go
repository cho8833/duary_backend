package event

import (
	"github.com/cho8833/duary_lambda/internal/util"
	uuid2 "github.com/google/uuid"
	"log"
	"time"
)

type Service interface {
	SaveEvent(req *SaveReq) (*VO, util.ApplicationError)
	GetEventBetweenStartAndEndDate(coupleId *string, rangeStartDate time.Time, rangeEndDate time.Time) ([]VO, util.ApplicationError)
	UpdateEvent(req *UpdateReq) (*VO, util.ApplicationError)
}

type ServiceImpl struct {
	repository     Repository
	scheduleHelper BridgeSchedulerHelper
}

func NewEventService(repository Repository, scheduleHelper BridgeSchedulerHelper) *ServiceImpl {
	return &ServiceImpl{repository: repository, scheduleHelper: scheduleHelper}
}

func (service *ServiceImpl) SaveEvent(req *SaveReq) (*VO, util.ApplicationError) {

	event := FromReq(req, *service.generateUID())

	vo, err := service.repository.SaveEvent(event)
	if err != nil {
		log.Printf("save event error\nreq: %+v\nerror: %s", req, err)
		return nil, util.DBSaveError{}
	}

	return vo, nil
}

func (service *ServiceImpl) UpdateEvent(req *UpdateReq) (*VO, util.ApplicationError) {

	vo, err := service.repository.UpdateEvent(req)

	if err != nil {
		log.Printf("update event error\nreq: %+v\nerror: %s", req, err)
		return nil, util.DBSaveError{}
	}

	return vo, nil
}

func (service *ServiceImpl) DeleteEvent(coupleId string, id string) util.ApplicationError {

	authContext := util.GetAuthContext()

	if coupleId != *authContext.CoupleId {
		log.Printf("unauthorized request : coupleId not matched\n user coupleId: %s, request coupleId: %s", *authContext.CoupleId, coupleId)
		return util.BadRequestError{}
	}

	err := service.repository.DeleteEvent(coupleId, id)

	if err != nil {
		return util.DBDeleteError{}
	}
	return nil
}

func (service *ServiceImpl) GetEventBetweenStartAndEndDate(coupleId string, rangeStartDate time.Time, rangeEndDate time.Time) ([]VO, util.ApplicationError) {

	candidateEvents, err := service.repository.FindByCoupleIdAndStartDateBefore(coupleId, rangeEndDate)
	if err != nil {
		log.Printf("failed to find events by coupleId and rangeStartDate before: %s", err)
		return nil, util.DBReadError{}
	}

	var result []VO

	for _, event := range candidateEvents {

		// 반복 이벤트인 경우 이벤트 인스턴스(Occurrence) 생성
		if event.Recurrence != nil {
			ocurrences, err := service.GenerateOccurrence(event, rangeStartDate, rangeEndDate)
			if err != nil {
				return nil, err
			}
			result = append(result, ocurrences...)

			// 비반복 이벤트인 경우 이베트가 range 내에 포함되는지 검사
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

func (service *ServiceImpl) GenerateOccurrence(vo VO, rangeStartDate time.Time, rangeEndDate time.Time) ([]VO, util.ApplicationError) {

	var occurrences []VO

	currentDate := vo.Recurrence.RepeatStartDate

	for {
		// occurrence 의 날짜가 반복종료(vo.Reccurrence.RepeatEndDate) 날짜 혹은 범위(rangeEndDate) 밖이면 종료
		if currentDate.After(rangeEndDate) || currentDate.After(vo.Recurrence.RepeatEndDate) {
			break
		}

		if (currentDate.After(rangeStartDate) || currentDate.Equal(rangeStartDate)) && currentDate.Before(rangeEndDate) {
			eventCopy := vo
			eventCopy.StartDateTime = time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())
			eventCopy.EndDateTime = time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())

			occurrences = append(occurrences, eventCopy)
		}

		if vo.Recurrence.Frequency == "DAILY" {
			currentDate = currentDate.AddDate(0, 0, int(1*vo.Recurrence.Interval))
		} else if vo.Recurrence.Frequency == "WEEKLY" {
			currentDate = currentDate.AddDate(0, 0, int(7*vo.Recurrence.Interval))
		} else if vo.Recurrence.Frequency == "MONTHLY" {
			currentDate = currentDate.AddDate(0, int(vo.Recurrence.Interval), 0)
		} else if vo.Recurrence.Frequency == "YEARLY" {
			currentDate = currentDate.AddDate(int(vo.Recurrence.Interval), 0, 0)
		} else {
			return nil, util.DBReadError{}
		}

	}
	return occurrences, nil
}

func (service *ServiceImpl) generateUID() *string {
	uuid := uuid2.New().String()
	return &uuid
}
