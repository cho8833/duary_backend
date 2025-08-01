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
	"math"
	"strconv"
	"time"
)

type SendFCMReq struct {
	Title          string   `json:"title"`
	Body           string   `json:"body"`
	TargetMemberId []string `json:"target_member_id"`
}

type SendAnniversaryReq struct {
	EventData      event.VO
	TargetMemberId []string
}

type BridgeSchedulerHelper struct {
	schedulerClient *scheduler.Client
}

func NewEventBridgeSchedulerHelper(schedulerClient *scheduler.Client) *BridgeSchedulerHelper {
	return &BridgeSchedulerHelper{schedulerClient: schedulerClient}
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

func (helper *BridgeSchedulerHelper) CreateAnniversarySchedule(ev event.VO, members []member.Member) shared.ApplicationError {
	var input scheduler.CreateScheduleInput

	var targetMemberIds []string
	for _, m := range members {
		targetMemberIds = append(targetMemberIds, m.GetId())
	}
	fcmReq := SendAnniversaryReq{
		EventData:      ev,
		TargetMemberId: targetMemberIds,
	}
	encoded, err := json.Marshal(fcmReq)
	if err != nil {
		log.Println(err)
		return shared.InternalServerError{}
	}
	encodedString := string(encoded)
	target := &types.Target{
		Arn:     aws.String("arn:aws:lambda:ap-northeast-2:922001515124:function:send_anniversary_fcm"),
		RoleArn: aws.String("arn:aws:iam::922001515124:role/SchedulerExecutionRole"),
		Input:   &encodedString,
	}

	if ev.Frequency == event.Yearly {
		yearly := ev.Yearly
		scheduleTime := time.Date(2000, yearly.Month, yearly.Day-1, 9, 0, 0, 0, time.UTC)
		var scheduleName string
		if ev.EventType == event.Anniversary {
			scheduleName = helper.GetYearlyAnniversaryScheduleName(ev.CoupleId)
		} else if ev.EventType == event.Birthday {
			scheduleName = helper.GetBirthdayScheduleName(ev.CoupleId, members[0].GetId())
		} else {
			log.Printf("attempt to create schedule with invalid event type: %+v\n", ev.EventType)
			return shared.BadRequestError{}
		}
		input = scheduler.CreateScheduleInput{
			Name:                       &scheduleName,
			FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			ScheduleExpression:         aws.String(helper.yearlyExpression(scheduleTime)),
			ScheduleExpressionTimezone: shared.StringPtr("UTC"),
			StartDate:                  shared.TimePtr(time.Now()),
			ActionAfterCompletion:      types.ActionAfterCompletionDelete,
			Target:                     target,
		}
	} else if ev.Frequency == event.Daily {
		// 1. 오늘까지 며칠째인지 계산 (UTC 기준)
		// 시간, 분, 초를 0으로 만들어 날짜만 비교
		utcNow := time.Now().UTC().Truncate(24 * time.Hour)
		utcStartDate := ev.RecurStartDate.Truncate(24 * time.Hour)

		// Sub은 Duration을 반환하므로 24시간으로 나누어 일수를 계산
		daysPassed := utcNow.Sub(utcStartDate).Hours() / 24

		// 2. 다음 N00일 기념일 찾기
		// 예: 365.2일 지났으면 -> floor(3.652) + 1 = 4 -> 400일
		// 예: 89일 지났으면 -> floor(0.89) + 1 = 1 -> 100일
		nextN := (math.Floor(daysPassed/100) + 1) * 100
		nextNValue := int(nextN)

		// 3. 목표 날짜 계산 (만난날 + N-1일)
		targetDate := utcStartDate.AddDate(0, 0, nextNValue-1).Add(time.Hour * 9)
		scheduleName := helper.Get100DayAnniversaryScheduleName(ev.CoupleId)
		input = scheduler.CreateScheduleInput{
			Name:                       &scheduleName,
			FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			ScheduleExpression:         aws.String(helper.dailyExpression(*ev.Daily)),
			ScheduleExpressionTimezone: shared.StringPtr("UTC"),
			StartDate:                  &targetDate,
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

func (helper *BridgeSchedulerHelper) UpdateAnniversarySchedule(ev event.VO, members []member.Member) shared.ApplicationError {
	var input scheduler.UpdateScheduleInput

	var targetMemberIds []string
	for _, m := range members {
		targetMemberIds = append(targetMemberIds, m.GetId())
	}
	fcmReq := SendAnniversaryReq{
		EventData:      ev,
		TargetMemberId: targetMemberIds,
	}
	encoded, err := json.Marshal(fcmReq)
	if err != nil {
		log.Println(err)
		return shared.InternalServerError{}
	}
	encodedString := string(encoded)
	target := &types.Target{
		Arn:     aws.String("arn:aws:lambda:ap-northeast-2:922001515124:function:send_anniversary_fcm"),
		RoleArn: aws.String("arn:aws:iam::922001515124:role/SchedulerExecutionRole"),
		Input:   &encodedString,
	}

	if ev.Frequency == event.Yearly {
		yearly := ev.Yearly
		scheduleTime := time.Date(2000, yearly.Month, yearly.Day-1, 9, 0, 0, 0, time.UTC)
		var scheduleName string
		if ev.EventType == event.Anniversary {
			scheduleName = helper.GetYearlyAnniversaryScheduleName(ev.CoupleId)
		} else if ev.EventType == event.Birthday {
			scheduleName = helper.GetBirthdayScheduleName(ev.CoupleId, members[0].GetId())
		} else {
			log.Printf("attempt to create schedule with invalid event type: %+v\n", ev.EventType)
			return shared.BadRequestError{}
		}
		input = scheduler.UpdateScheduleInput{
			Name:                       &scheduleName,
			FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			ScheduleExpression:         aws.String(helper.yearlyExpression(scheduleTime)),
			ScheduleExpressionTimezone: shared.StringPtr("UTC"),
			StartDate:                  shared.TimePtr(time.Now()),
			ActionAfterCompletion:      types.ActionAfterCompletionDelete,
			Target:                     target,
		}
	} else if ev.Frequency == event.Daily {
		// 1. 오늘까지 며칠째인지 계산 (UTC 기준)
		// 시간, 분, 초를 0으로 만들어 날짜만 비교
		utcNow := time.Now().UTC().Truncate(24 * time.Hour)
		utcStartDate := ev.RecurStartDate.Truncate(24 * time.Hour)

		// Sub은 Duration을 반환하므로 24시간으로 나누어 일수를 계산
		daysPassed := utcNow.Sub(utcStartDate).Hours() / 24

		// 2. 다음 N00일 기념일 찾기
		// 예: 365.2일 지났으면 -> floor(3.652) + 1 = 4 -> 400일
		// 예: 89일 지났으면 -> floor(0.89) + 1 = 1 -> 100일
		nextN := (math.Floor(daysPassed/100) + 1) * 100
		nextNValue := int(nextN)

		// 3. 목표 날짜 계산 (만난날 + N-1일 9시)
		targetDate := utcStartDate.AddDate(0, 0, nextNValue-1).Add(time.Hour * 9)
		scheduleName := helper.Get100DayAnniversaryScheduleName(ev.CoupleId)
		input = scheduler.UpdateScheduleInput{
			Name:                       &scheduleName,
			FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			ScheduleExpression:         aws.String(helper.dailyExpression(*ev.Daily)),
			ScheduleExpressionTimezone: shared.StringPtr("UTC"),
			StartDate:                  &targetDate,
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

func (helper *BridgeSchedulerHelper) GetScheduleName(eventId string, memberId string) string {
	return eventId + "_" + memberId
}

func (helper *BridgeSchedulerHelper) GetBirthdayScheduleName(coupleId string, memberId string) string {
	return coupleId + "_" + memberId + "_bir"
}

func (helper *BridgeSchedulerHelper) Get100DayAnniversaryScheduleName(coupleId string) string {
	return coupleId + "_ann100"
}
func (helper *BridgeSchedulerHelper) GetYearlyAnniversaryScheduleName(coupleId string) string {
	return coupleId + "_annYear"
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

func (helper *BridgeSchedulerHelper) yearlyExpression(scheduleTime time.Time) string {

	s := fmt.Sprintf("cron(%d %d %d %d ? *)", scheduleTime.Minute(),
		scheduleTime.Hour(),
		scheduleTime.Day(),
		scheduleTime.Month())
	return s
}
