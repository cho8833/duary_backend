package event

import "time"

type Event struct {
	Title      *string    `json:"title" dynamodbav:"title"`
	StartTime  *time.Time `json:"startTime" dynamodbav:"startDateTime"`
	EndTime    *time.Time `json:"endTime" dynamodbav:"endDateTime"`
	Content    *string    `json:"content" dynamodbav:"content"`
	IsTogether bool       `json:"isTogether" dynamodbav:"isTogether"`
	CoupleId   *string    `json:"coupleId" dynamodbav:"coupleId"`
	MemberId   *string    `json:"memberId" dynamodbav:"memberId"`
	IsAllDay   bool       `json:"isAllDay" dynamodbav:"isAllDay"`
	Repeat     *string    `json:"repeat" dynamodbav:"repeat"`
}
