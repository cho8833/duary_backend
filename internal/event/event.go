package event

import "time"

/*
StartDate, EndDate 의 값 중 Time.Year, Time.Month, Time.Day 만 쓰임
StartTime, StartTime 의 값 중 Time.Hour, Time.Minute 만 쓰임

StartDateTime 과 Recurrence.RepeatStartDate(EndDateTime 과 Recurrence.RepeatEndDate) 가 중복으로 있는 이유:

	반복 이벤트에 대한 Occurrence 를 만들 때 StartDateTime(EndDateTime) 에 Occurrence 의 StartDateTime(EndDateTime) 가 들어감.
*/
type Event struct {
	Id            *string     `json:"id" dynamodbav:"id"`
	CoupleId      *string     `json:"coupleId" dynamodbav:"coupleId"` // partitionKey
	StartDateTime time.Time   `json:"startDateTime" dynamodbav:"startDateTime"`
	EndDateTime   time.Time   `json:"endDateTime" dynamodbav:"endDateTime"`
	Title         string      `json:"title" dynamodbav:"title"`
	Content       *string     `json:"content" dynamodbav:"content"` // 메모
	IsTogether    bool        `json:"isTogether" dynamodbav:"isTogether"`
	CreatedBy     *int        `json:"createdBy" dynamodbav:"createdBy"`
	IsAllDay      bool        `json:"isAllDay" dynamodbav:"isAllDay"`
	Location      *string     `json:"location" dynamodbav:"location"`
	HangOutWith   *string     `json:"hangOutWith" dynamodbav:"hangOutWith"`
	Recurrence    *Recurrence `json:"recurrence" dynamodbav:"recurrence"`
}

type Recurrence struct {
	Frequency       string    `json:"frequency" dynamodbav:"frequency"`
	Interval        uint8     `json:"interval" dynamodbav:"interval"`
	RepeatStartDate time.Time `json:"repeatStartDate" dynamodbav:"repeatStartDate"`
	RepeatEndDate   time.Time `json:"repeatEndDate" dynamodbav:"repeatEndDate"`
}

func From(req *SaveEventReq) *Event {
	return &Event{
		Title:         req.Title,
		StartDateTime: req.StartDateTime,
		EndDateTime:   req.EndDateTime,
		IsTogether:    req.IsTogether,
		CoupleId:      req.CoupleId,
		CreatedBy:     req.CreatedBy,
		IsAllDay:      req.IsAllDay,
		Location:      req.Location,
		HangOutWith:   req.HangOutWith,
		Recurrence:    req.Recurrence,
	}
}

type SaveEventReq struct {
	Title         string    `json:"title"`
	CoupleId      *string   `json:"coupleId"`
	CreatedBy     *int      `json:"createdBy"`
	StartDateTime time.Time `json:"startDateTime"`
	EndDateTime   time.Time `json:"endDateTime"`
	Content       *string   `json:"content"`
	IsTogether    bool      `json:"isTogether"`
	IsAllDay      bool      `json:"isAllDay"`
	Location      *string   `json:"location"`
	HangOutWith   *string   `json:"hangOutWith"`

	Recurrence *Recurrence `json:"recurrence"`
}
