package scheduler

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/cho8833/duary_lambda/model/event"
	"github.com/cho8833/duary_lambda/model/member"
	"log"
	"time"
)

type BridgeSchedulerHelper struct {
	schedulerClient *scheduler.Client
}

func NewEventBridgeSchedulerHelper(schedulerClient *scheduler.Client) *BridgeSchedulerHelper {
	return &BridgeSchedulerHelper{schedulerClient: schedulerClient}
}

func (helper *BridgeSchedulerHelper) CreateEventSchedule(ev *event.VO, tm *member.Member) {

	var input scheduler.CreateScheduleInput

	// 단발성 일정인 경우
	if ev.Recurrence == nil {
		var scheduleTime time.Time
		if ev.CreatedBy == tm.GetId() { // 나의 일정인 경우
			scheduleTime = ev.StartDateTime.Add(event.AlarmOffsetMap[tm.MyAlarm].Duration)
		} else { // 상대방 일정인 경우
			scheduleTime = ev.StartDateTime.Add(event.AlarmOffsetMap[tm.LoverAlarm].Duration)
		}
		input = scheduler.CreateScheduleInput{
			Name:                       aws.String(helper.scheduleName(ev, tm.GetId())),
			FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			ScheduleExpression:         aws.String(helper.oneTimeExpression(scheduleTime)),
			ScheduleExpressionTimezone: aws.String("Asia/Seoul"),
			Target: &types.Target{
				Arn:     aws.String("arn:aws:lambda:ap-northeast-2:922001515124:function:send_fcm"),
				RoleArn: aws.String("arn:aws:iam::922001515124:role/SchedulerExecutionRole"),
			},
		}
		// 반복 일정인 경우
	} else {
		var scheduleTime time.Time

	}
	output, err := helper.schedulerClient.CreateSchedule(context.TODO(), &input)
	if err != nil {
		log.Printf("failed to create scheduler: %+v\n", err)
	} else {
		log.Printf("create schedule output: %+v\n", output)
	}
}

func (helper *BridgeSchedulerHelper) DeleteEventSchedule(event *event.VO) {
	input := scheduler.DeleteScheduleInput{
		Name: aws.String(helper.scheduleName(event)),
	}

	output, err := helper.schedulerClient.DeleteSchedule(context.TODO(), &input)
	if err != nil {
		log.Printf("failed to delete scheduler: %+v\n", err)
	} else {
		log.Printf("delete schedule output: %+v\n", output)
	}
}

func (helper *BridgeSchedulerHelper) UpdateEventSchedule(event *event.VO) {

}

func (helper *BridgeSchedulerHelper) scheduleName(event *event.VO, memberId string) string {
	return event.Id + "#" + memberId
}

func (helper *BridgeSchedulerHelper) oneTimeExpression(t time.Time) string {
	s := fmt.Sprintf("at(%d-%02d-%02dT%02d:%02d:%02d)",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	return s
}

func (helper *BridgeSchedulerHelper) recurrenceExpression(scheduleTime time.Time, recurrence event.Recurrence) string {

}
