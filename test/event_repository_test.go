package test

import (
	"fmt"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/event"
	"testing"
	"time"
)

//func Test_FindByCoupleIdAndStartDateBefore(t *testing.T) {
//	dynamodbClient := util.CreateLocalDynamoDBClient()
//
//	repository := event.NewRepository(dynamodbClient)
//
//}

/*
make sure dynamodb clean before running save!!!!!!!!!
*/

/*
java -Djava.library.path=./DynamoDBLocal_lib -jar DynamoDBLocal.jar -sharedDb

aws dynamodb create-table --table-name Event --attribute-definitions AttributeName=coupleId,AttributeType=S AttributeName=startDateTime,AttributeType=S --key-schema AttributeName=coupleId,KeyType=HASH AttributeName=startDateTime,KeyType=RANGE --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 --endpoint-url http://localhost:8000

aws dynamodb delete-table --table-name Event --endpoint-url http://localhost:8000
*/

func Test_UpdateEvent(t *testing.T) {
	dynamodbClient := CreateLocalDynamoDBClient()

	eventRepo := event.NewRepository(dynamodbClient)

	startRange := time.Now()
	events, _ := eventRepo.FindByCoupleIdAndStartDateBefore("couple1", startRange)

	vo := events[0]
	updateReq := event.UpdateReq{
		Content: ptr("updated content"),
	}

	updatedVO, err := eventRepo.Update(vo.CoupleId, vo.Id, &updateReq)
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

	eventRepo := event.NewRepository(dynamodbClient)

	testDate := time.Date(2025, 5, 1, 5, 0, 0, 0, time.Local)
	recurEndDate := testDate.AddDate(0, 1, 0)
	eventReq := event.SaveReq{
		Title:          "Team Meeting",
		CoupleId:       "couple1",
		CreatedBy:      "testuser1",
		StartDateTime:  testDate,
		EndDateTime:    testDate.Add(2 * time.Hour),
		RecurStartDate: &testDate,
		RecurEndDate:   &recurEndDate,
		Content:        ptr("Discuss project updates"),
		IsTogether:     true,
		IsAllDay:       false,
		Location:       ptr("Office"),
		HangOutWith:    ptr("Team"),
		Weekly: &event.WeeklyRecurrence{
			Weekdays: []event.Weekday{
				event.Weekday(time.Tuesday),
				event.Weekday(time.Thursday),
			},
		},
	}

	target := event.FromReq(&eventReq, "event1")

	saveEvent, err := eventRepo.Save(target)
	if err != nil {
		return
	}

	print(saveEvent)

}

func Test_GetEvent_match_couple_id_match_start_date_before(t *testing.T) {
	dynamodbClient, _ := model.GetDynamoDBClient()

	eventRepo := event.NewRepository(dynamodbClient)

	now := timePtr(time.Now())
	rangeStart := time.Date(now.Year(), now.Month()+1, now.Day(), 0, 0, 0, 0, now.Location())

	events, err := eventRepo.FindByCoupleIdAndStartDateBefore("dafe98ab-9eac-4b66-ba47-a389783d1a19", rangeStart)

	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("%#v", events)
}

func Test_GetEvent_not_match_couple_id(t *testing.T) {
	dynamodbClient := CreateLocalDynamoDBClient()

	eventRepo := event.NewRepository(dynamodbClient)

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

	eventRepo := event.NewRepository(dynamodbClient)

	rangeStart := time.Date(2024, 12, 12, 0, 0, 0, 0, time.Local)

	events, err := eventRepo.FindByCoupleIdAndStartDateBefore("couple1", rangeStart)

	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("%+v", events)
}

func Test_DeleteEvent_there_is_event_matching_id(t *testing.T) {
	dynamodbClient := CreateLocalDynamoDBClient()

	eventRepo := event.NewRepository(dynamodbClient)

	err := eventRepo.Delete("couple1", "2025-04-28T15:26:09+09:00#event1")

	if err != nil {
		t.Fatalf(err.Error())
	}

}

func ptr(s string) *string {
	return &s
}
