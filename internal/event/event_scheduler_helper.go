package event

//
//import (
//	"fmt"
//	"github.com/aws/aws-sdk-go-v2/aws"
//	"github.com/aws/aws-sdk-go-v2/service/scheduler"
//	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
//	"time"
//)
//
//type EventBridgeSchedulerHelper struct {
//	schedulerClient *scheduler.Client
//}
//
//func NewEventBridgeSchedulerHelper(schedulerClient *scheduler.Client) *EventBridgeSchedulerHelper {
//	return &EventBridgeSchedulerHelper{schedulerClient: schedulerClient}
//}
//
//func (helper *EventBridgeSchedulerHelper) CreateEventSchedule(event *VO, scheduleBefore int) {
//	// 15분 전 알림
//	scheduleTime := event.StartDateTime.Add(-15 * time.Minute)
//
//	input := scheduler.CreateScheduleInput{
//		Name:                       aws.String(event.CoupleId + "#" + event.Id),
//		FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
//		ScheduleExpression:         aws.String("at(" + helper.scheduleExpression(&scheduleTime) + ")"),
//		ScheduleExpressionTimezone: aws.String("Asia/Seoul"),
//		Target: &types.Target{
//
//		}
//	}
//}
//
//func (helper *EventBridgeSchedulerHelper) scheduleExpression(t *time.Time) string {
//	s := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
//		t.Year(), t.Month(), t.Day(),
//		t.Hour(), t.Minute(), t.Second())
//
//	return s
//}
