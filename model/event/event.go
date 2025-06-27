package event

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

/*
Event
StartDate, EndDate 의 값 중 Time.Year, Time.Month, Time.Day 만 쓰임

# StartTime, StartTime 의 값 중 Time.Hour, Time.Minute 만 쓰임

StartTime 과 Recurrence.RepeatStartDate(EndTime 과 Recurrence.RepeatEndDate) 가 중복으로 있는 이유:

	반복 이벤트에 대한 Occurrence 를 만들 때 StartTime(EndTime) 에 Occurrence 의 StartTime(EndTime) 가 들어감.
*/
type Event struct {
	CoupleId  string `dynamodbav:"coupleId"` // partition key
	CreatedBy string `dynamodbav:"createdBy"`

	StartDateTime  string             `dynamodbav:"startDateTime"` // sort key format : {startDateTime(ISO8601)}#{generated id}
	EndDateTime    time.Time          `dynamodbav:"endDateTime"`
	Frequency      Frequency          `dynamodbav:"frequency"`
	recurStartDate *time.Time         `dynamodbav:"recurStartDate"`
	recurEndDate   *time.Time         `dynamodbav:"recurEndDate"`
	Daily          *DailyRecurrence   `dynamodbav:"daily"`
	Weekly         *WeeklyRecurrence  `dynamodbav:"weekly"`
	Monthly        *MonthlyRecurrence `dynamodbav:"monthly"`
	Yearly         *YearlyRecurrence  `dynamodbav:"yearly"`

	Title       string  `dynamodbav:"title"`
	Content     *string `dynamodbav:"content"`
	IsTogether  bool    `dynamodbav:"isTogether"`
	IsAllDay    bool    `dynamodbav:"isAllDay"`
	Location    *string `dynamodbav:"location"`
	HangOutWith *string `dynamodbav:"hangOutWith"`

	EventType EventType `dynamodbav:"eventType"`
}

type DailyRecurrence struct {
	Interval int
}

type WeeklyRecurrence struct {
	Weekdays []Weekday
}

type MonthlyRecurrence struct {
	Days []int
}

type YearlyRecurrence struct {
	Month time.Month
	Day   int
}

type Frequency string

const (
	OneTime Frequency = "NONE"
	Daily   Frequency = "DAILY"
	Weekly  Frequency = "WEEKLY"
	Monthly Frequency = "MONTHLY"
	Yearly  Frequency = "YEARLY"
)

type EventType string

const (
	Normal      EventType = "NORMAL"
	Anniversary EventType = "ANNIVERSARY"
)

func (e Event) isRecurrence() bool {
	return e.Frequency != OneTime
}

type VO struct {
	Id        string `json:"id"` // id format : {startDate(ISO8601)}#{generatedId}
	CoupleId  string `json:"coupleId"`
	CreatedBy string `json:"createdBy"`

	StartDateTime  time.Time          `json:"startDateTime"`
	EndDateTime    time.Time          `json:"endDateTime"`
	Frequency      Frequency          `json:"frequency"`
	RecurStartDate *time.Time         `json:"recurStartDate"`
	RecurEndDate   *time.Time         `json:"recurEndDate"`
	Daily          *DailyRecurrence   `json:"daily"`
	Weekly         *WeeklyRecurrence  `json:"weekly"`
	Monthly        *MonthlyRecurrence `json:"monthly"`
	Yearly         *YearlyRecurrence  `json:"yearly"`
	RecurCount     int                `json:"recurCount"`

	Title       string  `json:"title"`
	Content     *string `json:"content"`
	Location    *string `json:"location"`
	HangOutWith *string `json:"hangOutWith"`
	IsTogether  bool    `json:"isTogether"`
	IsAllDay    bool    `json:"isAllDay"`

	EventType EventType `json:"eventType"`
}

func FromEvent(event Event) VO {
	sortKeySplit := strings.Split(event.StartDateTime, "#")
	startDateTime, _ := time.Parse(time.RFC3339, sortKeySplit[0])
	return VO{
		Id:        sortKeySplit[1],
		CoupleId:  event.CoupleId,
		CreatedBy: event.CreatedBy,

		StartDateTime:  startDateTime,
		EndDateTime:    event.EndDateTime,
		Frequency:      event.Frequency,
		RecurStartDate: event.recurStartDate,
		RecurEndDate:   event.recurEndDate,
		Daily:          event.Daily,
		Weekly:         event.Weekly,
		Monthly:        event.Monthly,
		Yearly:         event.Yearly,

		Title:       event.Title,
		Content:     event.Content,
		IsTogether:  event.IsTogether,
		IsAllDay:    event.IsAllDay,
		Location:    event.Location,
		HangOutWith: event.HangOutWith,

		EventType: event.EventType,
	}
}

/*
time.Weekday 는 Unmarshalling 불가능 -> custom type + UnMarshall/Marshall 구현
*/
type Weekday time.Weekday

var shortDayNames = []string{
	"SUN",
	"MON",
	"TUE",
	"WED",
	"THU",
	"FRI",
	"SAT",
}

var longDayNames = map[string]time.Weekday{
	"Sunday":    time.Sunday,
	"Monday":    time.Monday,
	"Tuesday":   time.Tuesday,
	"Wednesday": time.Wednesday,
	"Thursday":  time.Thursday,
	"Friday":    time.Friday,
	"Saturday":  time.Saturday,
}

func (w *Weekday) UnmarshalJSON(data []byte) error {
	// 문자열로 변환 (예: "Thursday")
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if val, ok := longDayNames[s]; ok {
		*w = Weekday(val)
		return nil
	}

	return fmt.Errorf("invalid weekday string: %s", s)
}

func (w Weekday) MarshalJSON() ([]byte, error) {
	day := time.Weekday(w).String() // "Thursday"
	return json.Marshal(day)
}

func (w Weekday) ShortDayName() string {
	return shortDayNames[w]
}
