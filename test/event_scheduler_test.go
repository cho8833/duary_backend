package test

import (
	"github.com/cho8833/duary_lambda/model/event"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/scheduler"
	"testing"
	"time"
)

func Test_CreateEventSchedule_Daily(t *testing.T) {
	schedulerClient, err := scheduler.GetSchedulerClient()
	if err != nil {
		t.Fatal(err)
	}

	schedulerHelper := scheduler.NewEventBridgeSchedulerHelper(schedulerClient)

	recurStartDate := time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)

	data := event.VO{
		Id:        "scheduler_test_event",
		Frequency: event.Daily,
		Daily: &event.DailyRecurrence{
			Interval: 3,
		},
		StartDateTime:  time.Date(2025, 5, 15, 9, 0, 0, 0, time.UTC),
		EndDateTime:    time.Date(2025, 5, 15, 11, 0, 0, 0, time.UTC),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}

	createdBy := member.Member{
		SocialId: "test_member",
		Provider: "GOOGLE",
	}
	err = schedulerHelper.CreateEventSchedule(data, createdBy)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_CreateEventSchedule_Weekly(t *testing.T) {
	schedulerClient, err := scheduler.GetSchedulerClient()
	if err != nil {
		t.Fatal(err)
	}
	schedulerHelper := scheduler.NewEventBridgeSchedulerHelper(schedulerClient)
	recurStartDate := time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)

	data := event.VO{
		Id:        "scheduler_test_event",
		Frequency: event.Weekly,
		Weekly: &event.WeeklyRecurrence{
			Weekdays: []event.Weekday{
				event.Weekday(time.Tuesday),
				event.Weekday(time.Thursday),
			},
		},
		StartDateTime:  time.Date(2025, 5, 15, 9, 0, 0, 0, time.UTC),
		EndDateTime:    time.Date(2025, 5, 15, 11, 0, 0, 0, time.UTC),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}
	createdBy := member.Member{
		SocialId: "test_member",
		Provider: "GOOGLE",
	}

	err = schedulerHelper.CreateEventSchedule(data, createdBy)
	if err != nil {
		t.Fatal(err)
	}
}
