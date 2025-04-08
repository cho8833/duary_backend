package test

import (
	"fmt"
	"github.com/cho8833/duary_lambda/internal/event"
	"testing"
	"time"
)

func getDailyDummy() event.Event {

	testDate := time.Date(2025, 4, 1, 0, 0, 0, 0, time.Local)

	dummy := event.Event{
		Title:       "Team Meeting",
		CoupleId:    ptr("couple1"),
		CreatedBy:   ptr("user1"),
		StartDate:   timePtr(testDate),
		EndDate:     timePtr(testDate.Add(2 * time.Hour)),
		StartTime:   timePtr(testDate),
		EndTime:     timePtr(testDate.Add(2 * time.Hour)),
		Content:     ptr("Discuss project updates"),
		IsTogether:  true,
		IsAllDay:    false,
		Location:    ptr("Office"),
		HangOutWith: ptr("Team"),
		Recurrence: &event.Recurrence{
			Frequency:       "daily",
			Interval:        1,
			RepeatStartDate: testDate,
			RepeatEndDate:   testDate.AddDate(0, 3, 0),
		},
	}
	return dummy
}

func getWeeklyDummy() event.Event {
	testDate := time.Date(2025, 4, 1, 0, 0, 0, 0, time.Local)
	dummy := event.Event{
		Title:       "Team Meeting",
		CoupleId:    ptr("couple1"),
		CreatedBy:   ptr("user1"),
		StartDate:   timePtr(testDate),
		EndDate:     timePtr(testDate.Add(2 * time.Hour)),
		StartTime:   timePtr(testDate),
		EndTime:     timePtr(testDate.Add(2 * time.Hour)),
		Content:     ptr("Discuss project updates"),
		IsTogether:  true,
		IsAllDay:    false,
		Location:    ptr("Office"),
		HangOutWith: ptr("Team"),
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
