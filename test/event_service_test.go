package test

import (
	"fmt"
	"github.com/cho8833/duary_lambda/model/event"
	"testing"
	"time"
)

func getDailyDummy() event.VO {

	testDate := time.Date(2025, 4, 1, 0, 0, 0, 0, time.Local)

	dummy := event.VO{
		Id:            "event1",
		Title:         "Team Meeting",
		CoupleId:      "couple1",
		CreatedBy:     "user1",
		StartDateTime: testDate,
		EndDateTime:   testDate.Add(time.Hour * 2),
		Content:       ptr("Discuss project updates"),
		IsTogether:    true,
		IsAllDay:      false,
		Location:      ptr("Office"),
		HangOutWith:   ptr("Team"),
		Recurrence: &event.Recurrence{
			Frequency:       "daily",
			Interval:        1,
			RepeatStartDate: testDate,
			RepeatEndDate:   testDate.AddDate(0, 3, 0),
		},
	}
	return dummy
}

func getWeeklyDummy() event.VO {
	testDate := time.Date(2025, 4, 1, 0, 0, 0, 0, time.Local)
	dummy := event.VO{
		Id:            "event2",
		Title:         "Team Meeting",
		CoupleId:      "couple1",
		CreatedBy:     "user2",
		StartDateTime: testDate,
		EndDateTime:   testDate.Add(time.Hour * 2),
		Content:       ptr("Discuss project updates"),
		IsTogether:    true,
		IsAllDay:      false,
		Location:      ptr("Office"),
		HangOutWith:   ptr("Team"),
		Recurrence: &event.Recurrence{
			Frequency:       "weekly",
			Interval:        1,
			RepeatStartDate: testDate,
			RepeatEndDate:   testDate.AddDate(0, 3, 0),
		},
	}
	return dummy
}

func Test_generateOccurrence_of_daily_by_one_day(t *testing.T) {
	dummy := getDailyDummy()

	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.Local)
	rangeEndDate := rangeStartDate.AddDate(0, 0, 1)

	service := &event.ServiceImpl{}

	result, err := service.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)

	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("GenerateOccurrence result length: expect 1, got %v", len(result))
	}

	fmt.Println(result)

}

func Test_generateOccurrence_of_daily_by_one_week(t *testing.T) {
	dummy := getDailyDummy()

	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.Local)
	rangeEndDate := rangeStartDate.AddDate(0, 0, 7)

	service := &event.ServiceImpl{}

	result, err := service.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)

	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}

	if len(result) != 7 {
		t.Fatalf("GenerateOccurrence result length: expect 7, got %v", len(result))
	}

	fmt.Println(result)
}

func Test_generateOccurrence_of_weekly_by_one_month(t *testing.T) {
	dummy := getWeeklyDummy()

	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.Local)
	rangeEndDate := rangeStartDate.AddDate(0, 1, 0)

	service := &event.ServiceImpl{}

	result, err := service.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)

	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}

	if len(result) != 4 {
		t.Fatalf("GenerateOccurrence result length: expect 1, got %v", len(result))
	}

	fmt.Println(result)
}
