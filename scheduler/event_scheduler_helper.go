package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/cho8833/duary_lambda/model/event"
	"github.com/cho8833/duary_lambda/model/member"
	"github.com/cho8833/duary_lambda/shared"
	"log"
	"strconv"
	"time"
)

type SendFCMReq struct {
	Title          string `json:"title"`
	Body           string `json:"body"`
	TargetMemberId string `json:"target_member_id"`
}

type BridgeSchedulerHelper struct {
	schedulerClient *scheduler.Client
}

func NewEventBridgeSchedulerHelper(schedulerClient *scheduler.Client) *BridgeSchedulerHelper {
	return &BridgeSchedulerHelper{schedulerClient: schedulerClient}
}

func (helper *BridgeSchedulerHelper) CreateEventSchedule(ev event.VO, tm member.Member) shared.ApplicationError {

	var input scheduler.CreateScheduleInput

	var duration time.Duration
	if ev.CreatedBy == tm.GetId() { // 나의 일정인 경우
		if tm.MyAlarm == event.None { // 나의 일정 알람을 끈 경우 아무것도 하지 않음
			return nil
		}
		duration = event.AlarmOffsetMap[tm.MyAlarm].Duration
	} else { // 상대방 일정인 경우
		if tm.LoverAlarm == event.None { // 연인의 일정 알람을 끈 경우 아무것도 하지 않음
			return nil
		}
		duration = event.AlarmOffsetMap[tm.LoverAlarm].Duration
	}

	fcmReq := SendFCMReq{
		Title:          ev.Title,
		TargetMemberId: tm.GetId(),
	}
	encoded, err := json.Marshal(fcmReq)
	if err != nil {
		log.Println(err)
		return shared.InternalServerError{}
	}
	encodedString := string(encoded)
	target := &types.Target{
		Arn:     aws.String("arn:aws:lambda:ap-northeast-2:922001515124:function:send_fcm"),
		RoleArn: aws.String("arn:aws:iam::922001515124:role/SchedulerExecutionRole"),
		Input:   &encodedString,
	}

	scheduleName := helper.GetScheduleName(ev.GetIdOnly(), tm.GetId())

	timeZone := "UTC"

	if ev.Frequency == event.OneTime {
		scheduleTime := ev.StartDateTime.Add(duration)
		input = scheduler.CreateScheduleInput{
			Name:                       &scheduleName,
			FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			ScheduleExpression:         aws.String(helper.oneTimeExpression(scheduleTime)),
			ScheduleExpressionTimezone: &timeZone,
			Target:                     target,
		}
	} else if ev.Frequency == event.Daily {
		startDate := time.Date(ev.RecurStartDate.Year(), ev.RecurStartDate.Month(), ev.RecurStartDate.Day(), ev.StartDateTime.Hour(), ev.StartDateTime.Minute(), 0, 0, time.UTC)
		var endDate *time.Time
		if ev.RecurEndDate != nil {
			temp := time.Date(ev.RecurEndDate.Year(), ev.RecurEndDate.Month(), ev.RecurEndDate.Day(), ev.EndDateTime.Hour(), ev.EndDateTime.Minute(), 0, 0, time.UTC)
			endDate = &temp
		}
		input = scheduler.CreateScheduleInput{
			Name:                       &scheduleName,
			FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			ScheduleExpression:         aws.String(helper.dailyExpression(*ev.Daily)),
			ScheduleExpressionTimezone: &timeZone,
			StartDate:                  &startDate,
			EndDate:                    endDate,
			ActionAfterCompletion:      types.ActionAfterCompletionDelete,
			Target:                     target,
		}
	} else if ev.Frequency == event.Weekly {
		scheduleTime := ev.StartDateTime.Add(duration)
		input = scheduler.CreateScheduleInput{
			Name:                       &scheduleName,
			FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			ScheduleExpression:         aws.String(helper.weeklyExpression(scheduleTime, *ev.Weekly)),
			ScheduleExpressionTimezone: &timeZone,
			StartDate:                  ev.RecurStartDate,
			EndDate:                    ev.RecurEndDate,
			ActionAfterCompletion:      types.ActionAfterCompletionDelete,
			Target:                     target,
		}
	} else if ev.Frequency == event.Monthly {
		scheduleTime := ev.StartDateTime.Add(duration)
		input = scheduler.CreateScheduleInput{
			Name:                       &scheduleName,
			FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			ScheduleExpression:         aws.String(helper.monthlyExpression(scheduleTime, *ev.Monthly)),
			ScheduleExpressionTimezone: &timeZone,
			StartDate:                  ev.RecurStartDate,
			EndDate:                    ev.RecurEndDate,
			ActionAfterCompletion:      types.ActionAfterCompletionDelete,
			Target:                     target,
		}
	} else if ev.Frequency == event.Yearly {
		scheduleTime := ev.StartDateTime.Add(duration)
		input = scheduler.CreateScheduleInput{
			Name:                       &scheduleName,
			FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			ScheduleExpression:         aws.String(helper.yearlyExpression(scheduleTime, *ev.Yearly)),
			ScheduleExpressionTimezone: &timeZone,
			StartDate:                  ev.RecurStartDate,
			EndDate:                    ev.RecurEndDate,
			ActionAfterCompletion:      types.ActionAfterCompletionDelete,
			Target:                     target,
		}
	}
	output, err := helper.schedulerClient.CreateSchedule(context.TODO(), &input)
	if err != nil {
		log.Printf("failed to create scheduler: %+v\n", err)
		return shared.InternalServerError{}
	} else {
		log.Printf("create schedule output: %+v\n", output)
	}
	return nil
}

func (helper *BridgeSchedulerHelper) UpdateEventSchedule(ev event.VO, tm member.Member) shared.ApplicationError {
	var input scheduler.UpdateScheduleInput

	var duration time.Duration
	if ev.CreatedBy == tm.GetId() { // 나의 일정인 경우
		if tm.MyAlarm == event.None { // 나의 일정 알람을 끈 경우 아무것도 하지 않음
			return nil
		}
		duration = event.AlarmOffsetMap[tm.MyAlarm].Duration
	} else { // 상대방 일정인 경우
		if tm.LoverAlarm == event.None { // 연인의 일정 알람을 끈 경우 아무것도 하지 않음
			return nil
		}
		duration = event.AlarmOffsetMap[tm.LoverAlarm].Duration
	}

	fcmReq := SendFCMReq{
		Title:          ev.Title,
		TargetMemberId: tm.GetId(),
	}
	encoded, err := json.Marshal(fcmReq)
	if err != nil {
		log.Println(err)
		return shared.InternalServerError{}
	}
	encodedString := string(encoded)
	target := &types.Target{
		Arn:     aws.String("arn:aws:lambda:ap-northeast-2:922001515124:function:send_fcm"),
		RoleArn: aws.String("arn:aws:iam::922001515124:role/SchedulerExecutionRole"),
		Input:   &encodedString,
	}

	scheduleName := helper.GetScheduleName(ev.GetIdOnly(), tm.GetId())

	timeZone := "UTC"

	if ev.Frequency == event.OneTime {
		scheduleTime := ev.StartDateTime.Add(duration)
		input = scheduler.UpdateScheduleInput{
			Name:                       &scheduleName,
			FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			ScheduleExpression:         aws.String(helper.oneTimeExpression(scheduleTime)),
			ScheduleExpressionTimezone: &timeZone,
			Target:                     target,
		}
	} else if ev.Frequency == event.Daily {
		startDate := time.Date(ev.RecurStartDate.Year(), ev.RecurStartDate.Month(), ev.RecurStartDate.Day(), ev.StartDateTime.Hour(), ev.StartDateTime.Minute(), 0, 0, time.UTC)
		var endDate *time.Time
		if ev.RecurEndDate != nil {
			temp := time.Date(ev.RecurEndDate.Year(), ev.RecurEndDate.Month(), ev.RecurEndDate.Day(), ev.EndDateTime.Hour(), ev.EndDateTime.Minute(), 0, 0, time.UTC)
			endDate = &temp
		}
		input = scheduler.UpdateScheduleInput{
			Name:                       &scheduleName,
			FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			ScheduleExpression:         aws.String(helper.dailyExpression(*ev.Daily)),
			ScheduleExpressionTimezone: &timeZone,
			StartDate:                  &startDate,
			EndDate:                    endDate,
			ActionAfterCompletion:      types.ActionAfterCompletionDelete,
			Target:                     target,
		}
	} else if ev.Frequency == event.Weekly {
		scheduleTime := ev.StartDateTime.Add(duration)
		input = scheduler.UpdateScheduleInput{
			Name:                       &scheduleName,
			FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			ScheduleExpression:         aws.String(helper.weeklyExpression(scheduleTime, *ev.Weekly)),
			ScheduleExpressionTimezone: &timeZone,
			StartDate:                  ev.RecurStartDate,
			EndDate:                    ev.RecurEndDate,
			ActionAfterCompletion:      types.ActionAfterCompletionDelete,
			Target:                     target,
		}
	} else if ev.Frequency == event.Monthly {
		scheduleTime := ev.StartDateTime.Add(duration)
		input = scheduler.UpdateScheduleInput{
			Name:                       &scheduleName,
			FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			ScheduleExpression:         aws.String(helper.monthlyExpression(scheduleTime, *ev.Monthly)),
			ScheduleExpressionTimezone: &timeZone,
			StartDate:                  ev.RecurStartDate,
			EndDate:                    ev.RecurEndDate,
			ActionAfterCompletion:      types.ActionAfterCompletionDelete,
			Target:                     target,
		}
	} else if ev.Frequency == event.Yearly {
		scheduleTime := ev.StartDateTime.Add(duration)
		input = scheduler.UpdateScheduleInput{
			Name:                       &scheduleName,
			FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			ScheduleExpression:         aws.String(helper.yearlyExpression(scheduleTime, *ev.Yearly)),
			ScheduleExpressionTimezone: &timeZone,
			StartDate:                  ev.RecurStartDate,
			EndDate:                    ev.RecurEndDate,
			ActionAfterCompletion:      types.ActionAfterCompletionDelete,
			Target:                     target,
		}
	}
	output, err := helper.schedulerClient.UpdateSchedule(context.TODO(), &input)
	if err != nil {
		log.Printf("failed to create scheduler: %+v\n", err)
		return shared.InternalServerError{}
	} else {
		log.Printf("create schedule output: %+v\n", output)
	}
	return nil
}

func (helper *BridgeSchedulerHelper) DeleteEventSchedule(scheduleName string) shared.ApplicationError {
	input := scheduler.DeleteScheduleInput{
		Name: &scheduleName,
	}

	output, err := helper.schedulerClient.DeleteSchedule(context.TODO(), &input)
	if err != nil {
		log.Printf("failed to delete scheduler: %+v\n", err)
		return shared.InternalServerError{}
	} else {
		log.Printf("delete schedule output: %+v\n", output)
		return nil
	}
}

func (helper *BridgeSchedulerHelper) GetScheduleName(eventId string, memberId string) string {
	return eventId + "-" + memberId
}

func (helper *BridgeSchedulerHelper) oneTimeExpression(t time.Time) string {
	s := fmt.Sprintf("at(%d-%02d-%02dT%02d:%02d:%02d)",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	return s
}

func (helper *BridgeSchedulerHelper) dailyExpression(daily event.DailyRecurrence) string {
	s := fmt.Sprintf("rate(%d days)", daily.Interval)
	return s
}

func (helper *BridgeSchedulerHelper) weeklyExpression(scheduleTime time.Time, weekly event.WeeklyRecurrence) string {

	s := "cron(" + strconv.Itoa(scheduleTime.Minute()) + " " + strconv.Itoa(scheduleTime.Hour()) + " " + "?" + " " + "*" + " "

	weekdays := weekly.Weekdays
	var days string
	for index, weekday := range weekdays {
		days += weekday.ShortDayName()
		if index < len(weekdays)-1 {
			days += "/"
		}
	}
	s += days + " " + "*)"
	return s
}

func (helper *BridgeSchedulerHelper) monthlyExpression(scheduleTime time.Time, monthly event.MonthlyRecurrence) string {
	s := "cron(" + strconv.Itoa(scheduleTime.Minute()) + " " + strconv.Itoa(scheduleTime.Hour()) + " "

	for index, day := range monthly.Days {
		s += strconv.Itoa(day)
		if index < len(monthly.Days)-1 {
			s += ","
		}
	}
	s += " * ? *)"
	return s
}

func (helper *BridgeSchedulerHelper) yearlyExpression(scheduleTime time.Time, yearly event.YearlyRecurrence) string {

	s := fmt.Sprintf("cron(%d %d %d %d ? *)", scheduleTime.Minute(),
		scheduleTime.Hour(),
		yearly.Day,
		int(yearly.Month))
	return s
}
