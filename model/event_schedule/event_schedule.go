package event_schedule

import (
	"github.com/cho8833/duary_lambda/model/event"
	"time"
)

type EventSchedule struct {
	MemberId      string            `dynamodbav:"memberId"` // partition key
	EventId       string            `dynamodbav:"eventId"`  // sort key
	ScheduleName  string            `dynamodbav:"scheduleName"`
	Recurrence    *event.Recurrence `dynamodbav:"recurrence"`
	StartDateTime time.Time         `dynamodbav:"startDateTime"`
	EndDateTime   time.Time         `dynamodbav:"endDateTime"`
}

type UpdateEventScheduleReq struct {
	ScheduleName  *string           `dynamodbav:"scheduleName"`
	Recurrence    *event.Recurrence `dynamodbav:"recurrence"`
	StartDateTime *time.Time        `dynamodbav:"startDateTime"`
	EndDateTime   *time.Time        `dynamodbav:"endDateTime"`
}
