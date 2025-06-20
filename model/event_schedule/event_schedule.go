package event_schedule

import (
	"github.com/cho8833/duary_lambda/model/event"
	"time"
)

type EventSchedule struct {
	MemberId     string           `dynamodbav:"memberId"` // partition key
	EventId      string           `dynamodbav:"eventId"`  // sort key
	ScheduleName string           `dynamodbav:"scheduleName"`
	StartDate    time.Time        `dynamodbav:"startDate"`
	EndDate      time.Time        `dynamodbav:"endDate"`
	StartTime    time.Time        `dynamodbav:"startTime"`
	Frequency    *event.Frequency `dynamodbav:"frequency"`
	DayOfWeek    *time.Weekday    `dynamodbav:"dayOfWeek"`
}

type UpdateEventScheduleReq struct {
	ScheduleName *string `dynamodbav:"scheduleName"`

	StartDate time.Time        `dynamodbav:"startDate"`
	EndDate   time.Time        `dynamodbav:"endDate"`
	StartTime time.Time        `dynamodbav:"startTime"`
	Frequency *event.Frequency `dynamodbav:"frequency"`
	DayOfWeek *time.Weekday    `dynamodbav:"dayOfWeek"`
}
