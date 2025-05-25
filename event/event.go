package event

import (
	"strings"
	"time"
)

/*
Event
StartDate, EndDate 의 값 중 Time.Year, Time.Month, Time.Day 만 쓰임

# StartTime, StartTime 의 값 중 Time.Hour, Time.Minute 만 쓰임

StartDateTime 과 Recurrence.RepeatStartDate(EndDateTime 과 Recurrence.RepeatEndDate) 가 중복으로 있는 이유:

	반복 이벤트에 대한 Occurrence 를 만들 때 StartDateTime(EndDateTime) 에 Occurrence 의 StartDateTime(EndDateTime) 가 들어감.
*/
type Event struct {
	CoupleId string `dynamodbav:"coupleId"` // partition key
	// sort key format : {startDateTime(ISO8601)}#{generated id}
	StartDateTime string      `dynamodbav:"startDateTime"` // sort key
	EndDateTime   time.Time   `dynamodbav:"endDateTime"`
	Title         string      `dynamodbav:"title"`
	Content       *string     `dynamodbav:"content"` // 메모
	IsTogether    bool        `dynamodbav:"isTogether"`
	CreatedBy     int         `dynamodbav:"createdBy"`
	IsAllDay      bool        `dynamodbav:"isAllDay"`
	Location      *string     `dynamodbav:"location"`
	HangOutWith   *string     `dynamodbav:"hangOutWith"`
	Recurrence    *Recurrence `dynamodbav:"recurrence"`
}

type Recurrence struct {
	Frequency       string    `json:"frequency" dynamodbav:"frequency"`
	Interval        uint8     `json:"interval" dynamodbav:"interval"`
	RepeatStartDate time.Time `json:"repeatStartDate" dynamodbav:"repeatStartDate"`
	RepeatEndDate   time.Time `json:"repeatEndDate" dynamodbav:"repeatEndDate"`
}

func FromReq(req *SaveReq, id string) *Event {
	sortKey := req.StartDateTime.Format(time.RFC3339) + "#" + id
	return &Event{
		Title:         req.Title,
		StartDateTime: sortKey,
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

type VO struct {
	// id format : {startDateTime(ISO8601)}#{generatedId}
	Id            string      `json:"id"`
	Title         string      `json:"title"`
	CoupleId      string      `json:"coupleId"`
	CreatedBy     int         `json:"createdBy"`
	StartDateTime time.Time   `json:"startDateTime"`
	EndDateTime   time.Time   `json:"endDateTime"`
	Content       *string     `json:"content"`
	IsTogether    bool        `json:"isTogether"`
	IsAllDay      bool        `json:"isAllDay"`
	Location      *string     `json:"location"`
	HangOutWith   *string     `json:"hangOutWith"`
	Recurrence    *Recurrence `json:"recurrence"`
}

func FromEvent(event Event) VO {
	sortKeySplit := strings.Split(event.StartDateTime, "#")
	startDateTime, _ := time.Parse(time.RFC3339, sortKeySplit[0])
	return VO{
		Id:            event.StartDateTime,
		CoupleId:      event.CoupleId,
		Title:         event.Title,
		StartDateTime: startDateTime,
		EndDateTime:   event.EndDateTime,
		Content:       event.Content,
		IsTogether:    event.IsTogether,
		IsAllDay:      event.IsAllDay,
		Location:      event.Location,
		HangOutWith:   event.HangOutWith,
		Recurrence:    event.Recurrence,
	}
}

type SaveReq struct {
	Title         string      `json:"title"`
	CoupleId      string      `json:"coupleId"`
	CreatedBy     int         `json:"createdBy"`
	StartDateTime time.Time   `json:"startDateTime"`
	EndDateTime   time.Time   `json:"endDateTime"`
	Content       *string     `json:"content"`
	IsTogether    bool        `json:"isTogether"`
	IsAllDay      bool        `json:"isAllDay"`
	Location      *string     `json:"location"`
	HangOutWith   *string     `json:"hangOutWith"`
	Recurrence    *Recurrence `json:"recurrence"`
}

type UpdateReq struct {
	Id            *string     `json:"id" dynamodbav:"id,omitempty"`
	CoupleId      *string     `json:"coupleId" dynamodbav:"coupleId,omitempty"`
	Title         *string     `json:"title" dynamodbav:"title,omitempty"`
	StartDateTime *time.Time  `json:"startDateTime" dynamodbav:"startDateTime,omitempty"`
	EndDateTime   *time.Time  `json:"endDateTime" dynamodbav:"endDateTime,omitempty"`
	Content       *string     `json:"content" dynamodbav:"content,omitempty"`
	IsTogether    *bool       `json:"isTogether" dynamodbav:"isTogether,omitempty"`
	IsAllDay      *bool       `json:"isAllDay" dynamodbav:"isAllDay,omitempty"`
	Location      *string     `json:"location" dynamodbav:"location,omitempty"`
	HangOutWith   *string     `json:"hangOutWith" dynamodbav:"hangOutWith,omitempty"`
	Recurrence    *Recurrence `json:"recurrence" dynamodbav:"recurrence,omitempty"`
}
