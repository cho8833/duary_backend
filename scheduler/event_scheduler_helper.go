package scheduler

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/cho8833/duary_lambda/model/event"
	"log"
	"time"
)

type BridgeSchedulerHelper struct {
	schedulerClient *scheduler.Client
}

func NewEventBridgeSchedulerHelper(schedulerClient *scheduler.Client) *BridgeSchedulerHelper {
	return &BridgeSchedulerHelper{schedulerClient: schedulerClient}
}

func (helper *BridgeSchedulerHelper) CreateEventSchedule(event *event.VO, offset event.AlarmOffset) {
	// 15분 전 알림
	scheduleTime := event.StartDateTime.Add(-15 * time.Minute)

	input := scheduler.CreateScheduleInput{
		Name:                       aws.String(helper.scheduleName(event)),
		FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		ScheduleExpression:         aws.String("at(" + helper.scheduleExpression(&scheduleTime) + ")"),
		ScheduleExpressionTimezone: aws.String("Asia/Seoul"),
		Target: &types.Target{
			Arn:     aws.String("arn:aws:lambda:ap-northeast-2:922001515124:function:send_fcm"),
			RoleArn: aws.String("arn:aws:iam::922001515124:role/SchedulerExecutionRole"),
		},
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

func (helper *BridgeSchedulerHelper) scheduleName(event *event.VO) string {
	return event.CoupleId + "#" + event.Id
}

func (helper *BridgeSchedulerHelper) scheduleExpression(t *time.Time) string {
	s := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	return s
}
