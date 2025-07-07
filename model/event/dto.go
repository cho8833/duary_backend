package event

import (
	"github.com/cho8833/duary_lambda/shared"
	"time"
)

type SaveReq struct {
	CoupleId  string
	CreatedBy string

	StartDateTime  time.Time          `json:"startDateTime"`
	EndDateTime    time.Time          `json:"endDateTime"`
	Frequency      Frequency          `json:"frequency"`
	RecurStartDate *time.Time         `json:"recurStartDate"`
	RecurEndDate   *time.Time         `json:"recurEndDate"`
	Daily          *DailyRecurrence   `json:"daily"`
	Weekly         *WeeklyRecurrence  `json:"weekly"`
	Monthly        *MonthlyRecurrence `json:"monthly"`
	Yearly         *YearlyRecurrence  `json:"yearly"`

	Title       string  `json:"title"`
	Content     *string `json:"content"`
	IsTogether  bool    `json:"isTogether"`
	IsAllDay    bool    `json:"isAllDay"`
	Location    *string `json:"location"`
	HangOutWith *string `json:"hangOutWith"`

	EventType EventType `json:"eventType"`
}

func (req *SaveReq) Validate() shared.ApplicationError {
	// 1. 단발성 이벤트인 경우 반복일정 필드는 nil 이어야 함
	if req.Frequency == OneTime {
		if req.RecurEndDate != nil || req.RecurStartDate != nil {
			return shared.ValidateError{Message: "단발성 이벤트는 반복일정 데이터가 없어야 합니다"}
		}
		// 2. 반복 이벤트인 경우 반복 시작 날짜가 지정되어야 함
	} else {
		if req.RecurStartDate == nil {
			return shared.ValidateError{Message: "반복 이벤트는 반복 시작 날짜가 지정되어야 합니다"}
		}
	}
	if req.Frequency == Daily {
		if req.Daily == nil || req.Daily.Interval <= 0 {
			return shared.ValidateError{Message: "일 반복 이벤트의 정보가 지정되지 않았습니다"}
		}
	}

	// 3. Weekly Event 인 경우 주간 반복 데이터가 있어야 함
	if req.Frequency == Weekly {
		if req.Weekly == nil || len(req.Weekly.Weekdays) == 0 {
			return shared.ValidateError{Message: "주 반복 이벤트의 정보가 지정되지 않았습니다"}
		}
	}
	// 4. Monthly Event 인 경우 월 반복 데이터가 있어야 함
	if req.Frequency == Monthly {
		if req.Monthly == nil || len(req.Monthly.Days) == 0 {
			return shared.ValidateError{Message: "월 반복 이벤트의 정보가 지정되지 않았습니다"}
		}
	}
	// 5. Yearly Event 인 경우 년 반복 데이터가 있어야 함
	if req.Frequency == Yearly {
		if req.Yearly == nil || req.Yearly.Day == 0 {
			return shared.ValidateError{Message: "년 반복 이벤트의 정보가 지정되지 않았습니다"}
		}
	}

	if req.EventType == "" {
		req.EventType = "NORMAL"
	}

	return nil
}

func FromReq(req *SaveReq, id string) *Event {
	sortKey := req.StartDateTime.Format(time.RFC3339) + "#" + id
	return &Event{
		CoupleId:  req.CoupleId,
		CreatedBy: req.CreatedBy,

		StartDateTime:  sortKey,
		EndDateTime:    req.EndDateTime,
		recurStartDate: req.RecurStartDate,
		recurEndDate:   req.RecurEndDate,
		Frequency:      req.Frequency,
		Daily:          req.Daily,
		Weekly:         req.Weekly,
		Monthly:        req.Monthly,
		Yearly:         req.Yearly,

		Title:       req.Title,
		Content:     req.Content,
		IsTogether:  req.IsTogether,
		IsAllDay:    req.IsAllDay,
		Location:    req.Location,
		HangOutWith: req.HangOutWith,

		EventType: req.EventType,
	}
}

type UpdateReq struct {
	Id string `json:"id"`

	StartDateTime  *time.Time         `json:"startDateTime" dynamodbav:"startDateTime,omitempty"`
	EndDateTime    *time.Time         `json:"endDateTime" dynamodbav:"endDateTime,omitempty"`
	RecurStartDate *time.Time         `json:"recurStartDate" dynamodbav:"recurStartDate,omitempty"`
	RecurEndDate   *time.Time         `json:"recurEndDate" dynamodbav:"recurEndDate,omitempty"`
	Frequency      *Frequency         `json:"frequency" dynamodbav:"frequency"`
	Daily          *DailyRecurrence   `json:"daily" dynamodbav:"daily"`
	Weekly         *WeeklyRecurrence  `json:"weekly" dynamodbav:"weekly"`
	Monthly        *MonthlyRecurrence `json:"monthly" dynamodbav:"monthly"`
	Yearly         *YearlyRecurrence  `json:"yearly" dynamodbav:"yearly"`

	Title       *string `json:"title" dynamodbav:"title,omitempty"`
	Content     *string `json:"content" dynamodbav:"content,omitempty"`
	IsTogether  *bool   `json:"isTogether" dynamodbav:"isTogether,omitempty"`
	IsAllDay    *bool   `json:"isAllDay" dynamodbav:"isAllDay,omitempty"`
	Location    *string `json:"location" dynamodbav:"location,omitempty"`
	HangOutWith *string `json:"hangOutWith" dynamodbav:"hangOutWith,omitempty"`
}

func (req *UpdateReq) Validate() shared.ApplicationError {
	if req.Frequency != nil {
		if *req.Frequency == Daily && req.Daily == nil {
			return shared.ValidateError{Message: "일 반복 이벤트의 정보가 지정되지 않았습니다"}
		}
		if *req.Frequency == Weekly && req.Weekly == nil {
			return shared.ValidateError{Message: "주 반복 이벤트의 정보가 지정되지 않았습니다"}
		}
		if *req.Frequency == Monthly && req.Monthly == nil {
			return shared.ValidateError{Message: "월 반복 이벤트의 정보가 지정되지 않았습니다"}
		}
		if *req.Frequency == Yearly && req.Yearly == nil {
			return shared.ValidateError{Message: "년 반복 이벤트의 정보가 지정되지 않았습니다"}
		}
	}

	return nil
}
