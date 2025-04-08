package event

import (
	"github.com/cho8833/duary_lambda/internal/util"
	"log"
	"time"
)

type Service interface {
	SaveEvent(req *SaveEventReq) (*Event, util.ApplicationError)
	GetEventBetweenStartAndEndDate(coupleId *string, rangeStartDate time.Time, rangeEndDate time.Time) ([]Event, util.ApplicationError)
}

type ServiceImpl struct {
	repository Repository
}

func NewEventService(repository Repository) *ServiceImpl {
	return &ServiceImpl{repository: repository}
}

func (service *ServiceImpl) SaveEvent(req *SaveEventReq) (*Event, util.ApplicationError) {
	event := From(req)

	event, err := service.repository.SaveEvent(event)
	if err != nil {
		log.Printf("save event error\nreq: %+v\nerror: %s", req, err)
		return nil, util.DBSaveError{}
	}

	return event, nil
}

func (service *ServiceImpl) GetEventBetweenStartAndEndDate(coupleId string, rangeStartDate time.Time, rangeEndDate time.Time) ([]Event, util.ApplicationError) {

	candidateEvents, err := service.repository.FindByCoupleIdAndStartDateBefore(coupleId, rangeEndDate)
	if err != nil {
		log.Printf("failed to find events by coupleId and rangeStartDate before: %s", err)
		return nil, util.DBReadError{}
	}

	var result []Event

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
		result = make([]Event, 0)
	}
	return result, nil
}

func (service *ServiceImpl) GenerateOccurrence(event Event, rangeStartDate time.Time, rangeEndDate time.Time) ([]Event, util.ApplicationError) {

	var occurrences []Event

	currentDate := event.Recurrence.RepeatStartDate

	for {
		// occurrence 의 날짜가 반복종료(event.Reccurrence.RepeatEndDate) 날짜 혹은 범위(rangeEndDate) 밖이면 종료
		if currentDate.After(rangeEndDate) || currentDate.After(event.Recurrence.RepeatEndDate) {
			break
		}

		if (currentDate.After(rangeStartDate) || currentDate.Equal(rangeStartDate)) && currentDate.Before(rangeEndDate) {
			eventCopy := event
			eventCopy.StartDateTime = time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())
			eventCopy.EndDateTime = time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())

			occurrences = append(occurrences, eventCopy)
		}

		if event.Recurrence.Frequency == "DAILY" {
			currentDate = currentDate.AddDate(0, 0, int(1*event.Recurrence.Interval))
		} else if event.Recurrence.Frequency == "WEEKLY" {
			currentDate = currentDate.AddDate(0, 0, int(7*event.Recurrence.Interval))
		} else if event.Recurrence.Frequency == "MONTHLY" {
			currentDate = currentDate.AddDate(0, int(event.Recurrence.Interval), 0)
		} else if event.Recurrence.Frequency == "YEARLY" {
			currentDate = currentDate.AddDate(int(event.Recurrence.Interval), 0, 0)
		} else {
			return nil, util.DBReadError{}
		}

	}
	return occurrences, nil
}

func timePtr(time time.Time) *time.Time {
	return &time
}
