package test

import (
	"fmt"
	"github.com/cho8833/duary_lambda/internal/event"
	"github.com/cho8833/duary_lambda/internal/test/util"
	"testing"
	"time"
)

//func Test_FindByCoupleIdAndStartDateBefore(t *testing.T) {
//	dynamodbClient := util.CreateLocalDynamoDBClient()
//
//	repository := event.NewEventRepository(dynamodbClient)
//
//}

/*
ensure dynamodb clean!!!!!!!!!
*/
func Test_SaveEvent(t *testing.T) {
	dynamodbClient := util.CreateLocalDynamoDBClient()

	eventRepo := event.NewEventRepository(dynamodbClient)

	now := timePtr(time.Now())
	eventReq := event.SaveEventReq{
		Title:       "Team Meeting",
		CoupleId:    ptr("couple1"),
		CreatedBy:   ptr("user1"),
		StartDate:   now,
		EndDate:     timePtr(now.Add(2 * time.Hour)),
		StartTime:   now,
		EndTime:     timePtr(now.Add(2 * time.Hour)),
		Content:     ptr("Discuss project updates"),
		IsTogether:  true,
		IsAllDay:    false,
		Location:    ptr("Office"),
		HangOutWith: ptr("Team"),
		Recurrence: &event.Recurrence{
			Frequency:       "weekly",
			Interval:        1,
			RepeatStartDate: *now,
			RepeatEndDate:   now.AddDate(0, 3, 0),
		},
	}

	target := event.From(&eventReq)

	saveEvent, err := eventRepo.SaveEvent(target)
	if err != nil {
		return
	}

	print(saveEvent)

}

func Test_GetEvent_match_couple_id_match_start_date_before(t *testing.T) {
	dynamodbClient := util.CreateLocalDynamoDBClient()

	eventRepo := event.NewEventRepository(dynamodbClient)

	now := timePtr(time.Now())
	rangeStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	events, err := eventRepo.FindByCoupleIdAndStartDateBefore("couple1", rangeStart)

	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("%+v", events)
}

func Test_GetEvent_not_match_couple_id(t *testing.T) {
	dynamodbClient := util.CreateLocalDynamoDBClient()

	eventRepo := event.NewEventRepository(dynamodbClient)

	now := timePtr(time.Now())

	rangeStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	events, err := eventRepo.FindByCoupleIdAndStartDateBefore("couple2", rangeStart)

	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("%+v", events)
}

func Test_GetEvent_match_couple_id_not_match_start_date_before(t *testing.T) {
	dynamodbClient := util.CreateLocalDynamoDBClient()

	eventRepo := event.NewEventRepository(dynamodbClient)

	rangeStart := time.Date(2024, 12, 12, 0, 0, 0, 0, time.Local)

	events, err := eventRepo.FindByCoupleIdAndStartDateBefore("couple1", rangeStart)

	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("%+v", events)
}

func ptr(s string) *string {
	return &s
}
func timePtr(time time.Time) *time.Time {
	return &time
}
