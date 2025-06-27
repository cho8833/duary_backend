package test

import (
	"fmt"
	"github.com/cho8833/duary_lambda/model/event"
	"github.com/cho8833/duary_lambda/shared"
	"testing"
	"time"
)

func Test_generateOccurrence_of_daily_by_one_day1(t *testing.T) {
	t.Logf("Daily Recur: Interval = 3, rangeStartDate=5/1, rangeEndDate=5/2, recurStartDate=5/1, recurEndDate=5/2 일 때, Occurrence 는 1개이어야 함")

	testTime := time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)
	testDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	svc := &event.ServiceImpl{}

	dummy := event.VO{
		Frequency: event.Daily,
		Daily: &event.DailyRecurrence{
			Interval: 3,
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &testDate,
		RecurEndDate:   &testDate,
	}
	rangeStartDate := testDate
	rangeEndDate := rangeStartDate.AddDate(0, 0, 1)

	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)

	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("GenerateOccurrence result length: expect 1, got %v", len(result))
	}

	if result[0].StartDateTime.Day() != rangeStartDate.Day() {
		t.Fatalf("Wrong Occurence %s", shared.ToString(result[0]))
	}

	fmt.Print(shared.ToString(result))
}

func Test_GenerateOccurrence_of_daily_by_one_day2(t *testing.T) {
	t.Logf("Daily recur: Interval=3, rangeStartDate=5/1, rangeEndDate=5/2, recurStartDate=4/30, recurEndDate=5/2 일 때 Occurence 는 없어야 함")
	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := rangeStartDate.AddDate(0, 0, 1)
	recurStartDate := time.Date(2025, 4, 30, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2025, 5, 2, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 4, 30, 9, 0, 0, 0, time.UTC)

	svc := &event.ServiceImpl{}

	dummy := event.VO{
		Frequency: event.Daily,
		Daily: &event.DailyRecurrence{
			Interval: 3,
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}

	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("GenerateOccurrence result length: expect 0, got %v", len(result))
	}
}

func Test_GenerateOccurrence_of_daily_by_one_day3(t *testing.T) {
	t.Logf("Daily Recur: Interval=1, rangeStartDate=5/1, rangeEndDate=5/2, recurStartDate=4/1, recurEndDate=5/1 일 때, Occurrence 는 1개이어야 함")

	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := rangeStartDate.AddDate(0, 0, 1)
	recurStartDate := time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 4, 1, 9, 0, 0, 0, time.UTC)

	svc := &event.ServiceImpl{}

	dummy := event.VO{
		Frequency: event.Daily,
		Daily: &event.DailyRecurrence{
			Interval: 1,
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}

	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("GenerateOccurrence result length: expect 1, got %v", len(result))
	}
	expected := time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)
	if !result[0].StartDateTime.Equal(expected) {
		t.Fatalf("Wrong Occurence %s", shared.ToString(result[0]))
	}
}

func Test_GenerateOccurrence_of_daily_by_one_day4(t *testing.T) {
	t.Logf("Daily Recur: Interval=1, rangeStartDate=5/1, rangeEndDate=5/2, recurStartDate=5/1, recurEndDate=5/2 일 때 Occurence 는 1개이어야 함")
	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := rangeStartDate.AddDate(0, 0, 1)
	recurStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	recurEndDate := recurStartDate.AddDate(0, 0, 1)
	testTime := time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)

	svc := &event.ServiceImpl{}
	dummy := event.VO{
		Frequency: event.Daily,
		Daily: &event.DailyRecurrence{
			Interval: 1,
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}

	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("GenerateOccurrence result length: expect 1, got %v", len(result))
	}
	expected := time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)
	if !result[0].StartDateTime.Equal(expected) {
		t.Fatalf("Wrong Occurence %s", shared.ToString(result[0]))
	}
}

func Test_GenerateOccurrence_of_daily_by_one_month1(t *testing.T) {
	t.Logf("Daily Recur: Interval=3, rangeStartDate=5/1, rangeEndDate=6/1, recurStartDate=5/1, recurEndDate=7/1 일 때 Occurence 는 11개이어야 함")

	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := rangeStartDate.AddDate(0, 1, 0)
	recurStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)

	svc := &event.ServiceImpl{}

	dummy := event.VO{
		Id:        "event1",
		Title:     "Team Meeting",
		CoupleId:  "couple1",
		CreatedBy: "user1",
		Frequency: event.Daily,
		Daily: &event.DailyRecurrence{
			Interval: 3,
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
		IsTogether:     true,
		IsAllDay:       false,
		EventType:      event.Normal,
	}

	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 11 {
		t.Fatalf("GenerateOccurrence result length: expect 11, got %v", len(result))
	}
	expected := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	for _, e := range result {
		if e.StartDateTime.Day() != expected.Day() {
			t.Fatalf("Wrong Occurence %s", shared.ToString(e))
		}
		expected = expected.AddDate(0, 0, 3)
	}

}

func Test_GenerateOccurrence_of_daily_by_one_month2(t *testing.T) {
	t.Logf("Daily Recur: interval=3, rangeStartDate=5/1, rangeEndDate=6/1, recurStartDate=4/29, recurEndDate=5/6 일 때 Occurrence는 2개이어야 함")

	testTime := time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)
	testDate := time.Date(2025, 4, 29, 0, 0, 0, 0, time.UTC)
	recurEndDate := testDate.AddDate(0, 0, 7)

	svc := &event.ServiceImpl{}

	dummy := event.VO{
		Id:        "event1",
		Title:     "Team Meeting",
		CoupleId:  "couple1",
		CreatedBy: "user1",
		Frequency: event.Daily,
		Daily: &event.DailyRecurrence{
			Interval: 3,
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &testDate,
		RecurEndDate:   &recurEndDate,
		IsTogether:     true,
		IsAllDay:       false,
		EventType:      event.Normal,
	}

	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := rangeStartDate.AddDate(0, 1, 0)

	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)

	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("GenerateOccurrence result length: expect 2, got %v", len(result))
	}
	expected1 := time.Date(2025, 5, 2, 0, 0, 0, 0, time.UTC)
	expected2 := time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC)
	if result[0].StartDateTime.Day() != expected1.Day() {
		t.Fatalf("Wrong Occurence %s", shared.ToString(result[0]))
	}
	if result[1].StartDateTime.Day() != expected2.Day() {
		t.Fatalf("Wrong Occurence %s", shared.ToString(result[1]))
	}
}

func Test_GenerateOccurrence_of_daily_by_one_month3(t *testing.T) {
	t.Logf("Daily Recur: interval=3, rangeStartDate=5/1, rangeEndDate=6/1, recurStartDate=5/30, recurEndDate=6/10 일 때 Occurence 는 1개이어야 함")

	recurStartDate := time.Date(2025, 5, 30, 9, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 5, 30, 9, 0, 0, 0, time.UTC)

	dummy := event.VO{
		Frequency: event.Daily,
		Daily: &event.DailyRecurrence{
			Interval: 3,
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}

	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := rangeStartDate.AddDate(0, 1, 0)
	svc := &event.ServiceImpl{}

	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("GenerateOccurrence result length: expect 1, got %v", len(result))
	}

	expected1 := time.Date(2025, 5, 30, 0, 0, 0, 0, time.UTC)
	if expected1.Day() != result[0].StartDateTime.Day() {
		t.Fatalf("Wrong Occurence %s", shared.ToString(result[0]))
	}
}

func Test_GenerateOccurrence_of_weekly_by_one_day1(t *testing.T) {
	t.Logf("Weekly Recur: 화, 목 반복, rangeStartDate=5/1(목), rangeEndDate=5/2, recurStartDate=5/1, recurEndDate=6/1 이면 Occurence 는 1개이어야 함")
	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := time.Date(2025, 5, 2, 0, 0, 0, 0, time.UTC)
	recurStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)

	dummy := event.VO{
		Frequency: event.Weekly,
		Weekly: &event.WeeklyRecurrence{
			Weekdays: []event.Weekday{
				event.Weekday(time.Tuesday),
				event.Weekday(time.Thursday),
			},
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}

	svc := &event.ServiceImpl{}
	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("GenerateOccurrence result length: expect 1, got %v", len(result))
	}
	expected1 := time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)
	if !result[0].StartDateTime.Equal(expected1) {
		t.Fatalf("Wrong Occurence %s", shared.ToString(result[0]))
	}
}

func Test_GenerateOccurrence_of_weekly_by_one_day2(t *testing.T) {
	t.Logf("Weekly recur: 화, 목 반복, rangeStartDate=4/30(수), rangeEndDate=5/1, recurStartDate=4/1, recurEndDate=5/1 이면 Occurrence 는 없어야 함")
	rangeStartDate := time.Date(2025, 4, 30, 0, 0, 0, 0, time.UTC)
	rangeEndDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	recurStartDate := time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 4, 1, 9, 0, 0, 0, time.UTC)

	svc := &event.ServiceImpl{}

	dummy := event.VO{
		Frequency: event.Weekly,
		Weekly: &event.WeeklyRecurrence{
			Weekdays: []event.Weekday{
				event.Weekday(time.Tuesday),
				event.Weekday(time.Thursday),
			},
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}

	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("GenerateOccurrence result length: expect 0, got %v", len(result))
	}
}

func Test_GenerateOccurrence_of_weekly_by_one_day3(t *testing.T) {
	t.Logf("Weekly Recur: 화 목 반복, rangeStartDate=5/1(목), rangeEndDate=5/2, recurStartDate=4/1, recurEndDate=5/1 이면 Occurrence 는 1개이어야 함")

	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := time.Date(2025, 5, 2, 0, 0, 0, 0, time.UTC)
	recurStartDate := time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 4, 1, 9, 0, 0, 0, time.UTC)
	dummy := event.VO{
		Frequency: event.Weekly,
		Weekly: &event.WeeklyRecurrence{
			Weekdays: []event.Weekday{
				event.Weekday(time.Tuesday),
				event.Weekday(time.Thursday),
			},
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}

	svc := &event.ServiceImpl{}
	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("GenerateOccurrence result length: expect 1, got %v", len(result))
	}
	expected1 := time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)
	if !result[0].StartDateTime.Equal(expected1) {
		t.Fatalf("Wrong Occurence %s", shared.ToString(result[0]))
	}
}

func Test_GenerateOccurrence_of_weekly_by_one_month1(t *testing.T) {
	t.Logf("Weekly Recur: 수 금 반복, rangeStartDate=5/1(목), rangeEndDate=6/1(일), recurStartDate=4/15, recurEndDate=5/15 일 때 Occurrence 는 4개이어야 함")
	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	recurStartDate := time.Date(2025, 4, 15, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 4, 15, 9, 0, 0, 0, time.UTC)

	dummy := event.VO{
		Frequency: event.Weekly,
		Weekly: &event.WeeklyRecurrence{
			Weekdays: []event.Weekday{
				event.Weekday(time.Wednesday),
				event.Weekday(time.Friday),
			},
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}

	svc := &event.ServiceImpl{}
	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 4 {
		t.Fatalf("GenerateOccurrence result length: expect 4, got %v", len(result))
	}
	expected := []time.Time{
		time.Date(2025, 5, 2, 9, 0, 0, 0, time.UTC),
		time.Date(2025, 5, 7, 9, 0, 0, 0, time.UTC),
		time.Date(2025, 5, 9, 9, 0, 0, 0, time.UTC),
		time.Date(2025, 5, 14, 9, 0, 0, 0, time.UTC),
	}

	for i := range result {
		if !result[i].StartDateTime.Equal(expected[i]) {
			t.Fatalf("Wrong Occurence %s", shared.ToString(result[i]))
		}
	}
}

func Test_GenerateOccurrence_of_weekly_by_one_month2(t *testing.T) {
	t.Logf("Weekly Recur: 수 금 반복, rangeStartDate=5/1(목), rangeEndDate=6/1, recurStartDate=5/15(목), recurEndDate=6/15일 때, Occurrence 는 5개이어야 함")

	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	recurStartDate := time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 5, 15, 9, 0, 0, 0, time.UTC)

	svc := &event.ServiceImpl{}

	dummy := event.VO{
		Frequency: event.Weekly,
		Weekly: &event.WeeklyRecurrence{
			Weekdays: []event.Weekday{
				event.Weekday(time.Wednesday),
				event.Weekday(time.Friday),
			},
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}

	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 5 {
		t.Fatalf("GenerateOccurrence result length: expect 5, got %v", len(result))
	}

	expected := []time.Time{
		time.Date(2025, 5, 16, 9, 0, 0, 0, time.UTC),
		time.Date(2025, 5, 21, 9, 0, 0, 0, time.UTC),
		time.Date(2025, 5, 23, 9, 0, 0, 0, time.UTC),
		time.Date(2025, 5, 28, 9, 0, 0, 0, time.UTC),
		time.Date(2025, 5, 30, 9, 0, 0, 0, time.UTC),
	}

	for i := range result {
		if !result[i].StartDateTime.Equal(expected[i]) {
			t.Fatalf("Wrong Occurence %s", shared.ToString(result[i]))
		}
	}
}

func Test_GenerateOccurrences_of_monthly_by_one_day1(t *testing.T) {
	t.Logf("Monthly Recur: 1, 10일 반복, rangeStartDate=5/1, rangeEndDate=5/2, recurStartDate=5/1, recurEndDate=7/1 일 때 Occurrence 는 1개이어야 함")

	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := time.Date(2025, 5, 2, 0, 0, 0, 0, time.UTC)
	recurStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)
	svc := &event.ServiceImpl{}

	dummy := event.VO{
		Frequency: event.Monthly,
		Monthly: &event.MonthlyRecurrence{
			Days: []int{
				1, 10,
			},
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}

	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("GenerateOccurrence result length: expect 1, got %v", len(result))
	}
	expected := time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)
	if !result[0].StartDateTime.Equal(expected) {
		t.Fatalf("Wrong Occurence %s", shared.ToString(result[0]))
	}
}

func Test_GenerateOccurrences_of_monthly_by_one_day2(t *testing.T) {
	t.Logf("Monthnly Recur: 2, 10일 반복, rangeStartDate=5/1, rangeEndDate=5/2, recurStartDate=5/1, recurEndDate=7/1 일 때 Occurrence 는 없어야 함")
	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := time.Date(2025, 5, 2, 0, 0, 0, 0, time.UTC)
	recurStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)
	svc := &event.ServiceImpl{}

	dummy := event.VO{
		Frequency: event.Monthly,
		Monthly: &event.MonthlyRecurrence{
			Days: []int{
				2, 10,
			},
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}

	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("GenerateOccurrence result length: expect 0, got %v", len(result))
	}
}

func Test_GenerateOccurrences_of_monthly_by_one_day3(t *testing.T) {
	t.Logf("Monthly Recur: 1, 10일 반복, rangeStartDate=5/1, rangeEndate=5/2, recurStartDate=4/1, recurEndDate=5/1일 때 occurrence 는 1개이어야 함")
	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := time.Date(2025, 5, 2, 0, 0, 0, 0, time.UTC)
	recurStartDate := time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 4, 1, 9, 0, 0, 0, time.UTC)
	svc := &event.ServiceImpl{}
	dummy := event.VO{
		Frequency: event.Monthly,
		Monthly: &event.MonthlyRecurrence{
			Days: []int{
				1, 10,
			},
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}
	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("GenerateOccurrence result length: expect 1, got %v", len(result))
	}
	expected := time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)
	if !result[0].StartDateTime.Equal(expected) {
		t.Fatalf("Wrong Occurence %s", shared.ToString(result[0]))
	}
}

func Test_GenerateOccurrences_of_monthly_by_one_month1(t *testing.T) {
	t.Logf("Monthly Recur: 5, 10일 반복, rangeStartDate=5/1, rangeEndDate=6/1, recurStartDate=5/6, recurEndDate=6/1 일 때 Occurrence 는 1개이어야 함")

	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	recurStartDate := time.Date(2025, 5, 6, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 5, 6, 9, 0, 0, 0, time.UTC)

	svc := &event.ServiceImpl{}

	dummy := event.VO{
		Frequency: event.Monthly,
		Monthly: &event.MonthlyRecurrence{
			Days: []int{
				5, 10,
			},
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}

	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("GenerateOccurrence result length: expect 1, got %v", len(result))
	}
	expected := time.Date(2025, 5, 10, 9, 0, 0, 0, time.UTC)
	if !result[0].StartDateTime.Equal(expected) {
		t.Fatalf("Wrong Occurence %s", shared.ToString(result[0]))
	}
}

func Test_GenerateOccurrences_of_monthly_by_one_month2(t *testing.T) {
	t.Logf("Monthly Recur: 5, 10일 반복, rangeStartDate=5/1, rangeEndDate=6/1, recurStartDate=4/15, recurEndDate=5/9 일 때 occurrence 는 1 개 이어야 함")

	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	recurStartDate := time.Date(2025, 4, 15, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2025, 5, 9, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 4, 15, 9, 0, 0, 0, time.UTC)
	svc := &event.ServiceImpl{}
	dummy := event.VO{
		Frequency: event.Monthly,
		Monthly: &event.MonthlyRecurrence{
			Days: []int{
				5, 10,
			},
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}

	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("GenerateOccurrence result length: expect 1, got %v", len(result))
	}
	expected := time.Date(2025, 5, 5, 9, 0, 0, 0, time.UTC)
	if !result[0].StartDateTime.Equal(expected) {
		t.Fatalf("Wrong Occurence %s", shared.ToString(result[0]))
	}
}

func Test_GenerateOccurrences_of_monthly_by_one_month3(t *testing.T) {
	t.Logf("Monthly Recur: 5, 10일 반복, rangeStartDate=5/1, rangeEndDate=6/1, recurStartDate=4/1, recurEndDate=6/1 일 때 occurrence 는 2개이어야 함")
	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	recurStartDate := time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 4, 1, 9, 0, 0, 0, time.UTC)
	svc := &event.ServiceImpl{}
	dummy := event.VO{
		Frequency: event.Monthly,
		Monthly: &event.MonthlyRecurrence{
			Days: []int{
				5, 10,
			},
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}
	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("GenerateOccurrence result length: expect 2, got %v", len(result))
	}
	expected := []time.Time{
		time.Date(2025, 5, 5, 9, 0, 0, 0, time.UTC),
		time.Date(2025, 5, 10, 9, 0, 0, 0, time.UTC),
	}
	for i := range result {
		if !result[i].StartDateTime.Equal(expected[i]) {
			t.Fatalf("Wrong Occurence %s", shared.ToString(result[i]))
		}
	}
}

func Test_GenerateOccurrences_of_yearly_by_one_day1(t *testing.T) {
	t.Logf("Yearly Recur: 5/1, rangeStartDate=5/1, rangeEndDate=5/2, recurStartDate=5/1, recurEndDate=2026/5/1일 때 occurrence 는 1개이어야 함")
	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := time.Date(2025, 5, 2, 0, 0, 0, 0, time.UTC)
	recurStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)
	svc := &event.ServiceImpl{}

	dummy := event.VO{
		Frequency: event.Yearly,
		Yearly: &event.YearlyRecurrence{
			Month: 5,
			Day:   1,
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}

	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("GenerateOccurrence result length: expect 1, got %v", len(result))
	}
	expected := time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)
	if !result[0].StartDateTime.Equal(expected) {
		t.Fatalf("Wrong Occurence %s", shared.ToString(result[0]))
	}
}

func Test_GenerateOccurrences_of_yearly_by_one_day2(t *testing.T) {
	t.Logf("Yearly Recur: 5/2, rangeStartDate=5/1, rangeEndDate=5/2, recurStartDate=5/1, recurEndDate=2026/5/1일 때 occurrence 는 없어야 함")
	rangeStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := time.Date(2025, 5, 2, 0, 0, 0, 0, time.UTC)
	recurStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)
	svc := &event.ServiceImpl{}

	dummy := event.VO{
		Frequency: event.Yearly,
		Yearly: &event.YearlyRecurrence{
			Month: 5,
			Day:   2,
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}

	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("GenerateOccurrence result length: expect 0, got %v", len(result))
	}
}

func Test_GenerateOccurrences_of_yearly_by_one_day3(t *testing.T) {
	t.Logf("Yearly Recur: 5/2, rangeStartDate=2026/5/1, rangeEndDate=2026/5/2, recurStartDate=2025/5/1, recurEndDate=2026/5/1일 때 occurrence 는 1개이어야 함")
	rangeStartDate := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	rangeEndDate := time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)
	recurStartDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	recurEndDate := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	testTime := time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)
	svc := &event.ServiceImpl{}

	dummy := event.VO{
		Frequency: event.Yearly,
		Yearly: &event.YearlyRecurrence{
			Month: 5,
			Day:   1,
		},
		StartDateTime:  testTime,
		EndDateTime:    testTime.Add(time.Hour * 2),
		RecurStartDate: &recurStartDate,
		RecurEndDate:   &recurEndDate,
	}

	result, err := svc.GenerateOccurrence(dummy, rangeStartDate, rangeEndDate)
	if err != nil {
		t.Fatalf("GenerateOccurrence err: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("GenerateOccurrence result length: expect 1, got %v", len(result))
	}
	expected := time.Date(2026, 5, 1, 9, 0, 0, 0, time.UTC)
	if !result[0].StartDateTime.Equal(expected) {
		t.Fatalf("Wrong Occurence %s", shared.ToString(result[0]))
	}
}
