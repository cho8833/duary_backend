package test

import (
	"fmt"
	"github.com/cho8833/duary_lambda/model/event"
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
ensure dynamodb clean before running save!!!!!!!!!
*/

func Test_UpdateEvent(t *testing.T) {
	dynamodbClient := CreateLocalDynamoDBClient()

	eventRepo := event.NewEventRepository(dynamodbClient)

	startRange := time.Now()
	events, _ := eventRepo.FindByCoupleIdAndStartDateBefore("couple1", startRange)

	vo := events[0]
	updateReq := event.UpdateReq{
		Id:       &vo.Id,
		CoupleId: &vo.CoupleId,
		Content:  ptr("updated content"),
	}

	updatedVO, err := eventRepo.UpdateEvent(&updateReq)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if *updatedVO.Content != "updated content" {
		t.Fatalf("%+v", updatedVO.Content)
	}
}

func Test_UpdateEvent_updateStartDateTime(t *testing.T) {

}
func Test_SaveEvent(t *testing.T) {
	dynamodbClient := CreateLocalDynamoDBClient()

	eventRepo := event.NewEventRepository(dynamodbClient)

	now := time.Now()
	eventReq := event.SaveReq{
		Title:         "Team Meeting",
		CoupleId:      "couple1",
		CreatedBy:     "testuser1",
		StartDateTime: now,
		EndDateTime:   now.Add(2 * time.Hour),
		Content:       ptr("Discuss project updates"),
		IsTogether:    true,
		IsAllDay:      false,
		Location:      ptr("Office"),
		HangOutWith:   ptr("Team"),
		Recurrence: &event.Recurrence{
			Frequency:       "weekly",
			Interval:        1,
			RepeatStartDate: now,
			RepeatEndDate:   now.AddDate(0, 3, 0),
		},
	}

	target := event.FromReq(&eventReq, "event1")

	saveEvent, err := eventRepo.SaveEvent(target)
	if err != nil {
		return
	}

	print(saveEvent)

}

func Test_GetEvent_match_couple_id_match_start_date_before(t *testing.T) {
	dynamodbClient := CreateLocalDynamoDBClient()

	eventRepo := event.NewEventRepository(dynamodbClient)

	now := timePtr(time.Now())
	rangeStart := time.Date(now.Year(), now.Month()+1, now.Day(), 0, 0, 0, 0, now.Location())

	events, err := eventRepo.FindByCoupleIdAndStartDateBefore("couple1", rangeStart)

	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("%#v", events)
}

func Test_GetEvent_not_match_couple_id(t *testing.T) {
	dynamodbClient := CreateLocalDynamoDBClient()

	eventRepo := event.NewEventRepository(dynamodbClient)

	now := timePtr(time.Now())

	rangeStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	events, err := eventRepo.FindByCoupleIdAndStartDateBefore("couple2", rangeStart)

	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("%#v", events)
}

func Test_GetEvent_match_couple_id_not_match_start_date_before(t *testing.T) {
	dynamodbClient := CreateLocalDynamoDBClient()

	eventRepo := event.NewEventRepository(dynamodbClient)

	rangeStart := time.Date(2024, 12, 12, 0, 0, 0, 0, time.Local)

	events, err := eventRepo.FindByCoupleIdAndStartDateBefore("couple1", rangeStart)

	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("%+v", events)
}

func Test_DeleteEvent_there_is_event_matching_id(t *testing.T) {
	dynamodbClient := CreateLocalDynamoDBClient()

	eventRepo := event.NewEventRepository(dynamodbClient)

	err := eventRepo.DeleteEvent("couple1", "2025-04-28T15:26:09+09:00#event1")

	if err != nil {
		t.Fatalf(err.Error())
	}

}

func ptr(s string) *string {
	return &s
}
func timePtr(time time.Time) *time.Time {
	return &time
}
