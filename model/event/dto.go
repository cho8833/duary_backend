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

	// 3. endDateTime 은 startDateTime 이후 이어야 함
	if req.EndDateTime.Before(req.StartDateTime) {
		return shared.ValidateError{Message: "시작 시간은 종료 시간보다 이전일 수 없습니다"}
	}

	// 4. endDate >= startDateTime + minute * 5
	difference := req.EndDateTime.Sub(req.StartDateTime)
	if difference.Minutes() < 5 {
		return shared.ValidateError{Message: "일정 진행 시간은 5분보다 길어야 합니다"}
	}

	// 5.1. Daily Event 인 경우 일간 반복 데이터가 있어야 함
	if req.Frequency == Daily {
		if req.Daily == nil || req.Daily.Interval <= 0 {
			return shared.ValidateError{Message: "일 반복 이벤트의 정보가 지정되지 않았습니다"}
		}
	}

	// 5,2. Weekly Event 인 경우 주간 반복 데이터가 있어야 함
	if req.Frequency == Weekly {
		if req.Weekly == nil || len(req.Weekly.Weekdays) == 0 {
			return shared.ValidateError{Message: "주 반복 이벤트의 정보가 지정되지 않았습니다"}
		}
	}
	// 5.3. Monthly Event 인 경우 월 반복 데이터가 있어야 함
	if req.Frequency == Monthly {
		if req.Monthly == nil || len(req.Monthly.Days) == 0 {
			return shared.ValidateError{Message: "월 반복 이벤트의 정보가 지정되지 않았습니다"}
		}
	}
	// 5.4. Yearly Event 인 경우 년 반복 데이터가 있어야 함
	if req.Frequency == Yearly {
		if req.Yearly == nil || req.Yearly.Day == 0 {
			return shared.ValidateError{Message: "년 반복 이벤트의 정보가 지정되지 않았습니다"}
		}
	}

	if req.EventType == "" {
		req.EventType = "NORMAL"
	}

	// 반복 데이터 정리
	switch req.Frequency {
	case Daily:
		req.Weekly = nil
		req.Monthly = nil
		req.Yearly = nil
	case Weekly:
		req.Daily = nil
		req.Monthly = nil
		req.Yearly = nil
	case Monthly:
		req.Daily = nil
		req.Weekly = nil
		req.Yearly = nil
	case Yearly:
		req.Daily = nil
		req.Weekly = nil
		req.Monthly = nil
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
		RecurStartDate: req.RecurStartDate,
		RecurEndDate:   req.RecurEndDate,
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

type EditReq struct {
	StartDateTime  *time.Time         `json:"startDateTime" dynamodbav:"-"`
	EndDateTime    time.Time          `json:"endDateTime" dynamodbav:"endDateTime"`
	Frequency      Frequency          `json:"frequency" dynamodbav:"frequency"`
	RecurStartDate *time.Time         `json:"recurStartDate" dynamodbav:"recurStartDate"`
	RecurEndDate   *time.Time         `json:"recurEndDate" dynamodbav:"recurEndDate"`
	Daily          *DailyRecurrence   `json:"daily" dynamodbav:"daily"`
	Weekly         *WeeklyRecurrence  `json:"weekly" dynamodbav:"weekly"`
	Monthly        *MonthlyRecurrence `json:"monthly" dynamodbav:"monthly"`
	Yearly         *YearlyRecurrence  `json:"yearly" dynamodbav:"yearly"`

	Title       string  `json:"title" dynamodbav:"title"`
	Content     *string `json:"content" dynamodbav:"content"`
	IsTogether  bool    `json:"isTogether" dynamodbav:"isTogether"`
	IsAllDay    bool    `json:"isAllDay" dynamodbav:"isAllDay"`
	Location    *string `json:"location" dynamodbav:"location"`
	HangOutWith *string `json:"hangOutWith" dynamodbav:"hangOutWith"`
}

func (req *EditReq) Validate() shared.ApplicationError {
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

	// 3. endDateTime 은 startDateTime 이후 이어야 함
	if req.EndDateTime.Before(*req.StartDateTime) {
		return shared.ValidateError{Message: "시작 시간은 종료 시간보다 이전일 수 없습니다"}
	}

	// 4. endDate >= startDateTime + minute * 5
	difference := req.EndDateTime.Sub(*req.StartDateTime)
	if difference.Minutes() < 5 {
		return shared.ValidateError{Message: "일정 진행 시간은 5분보다 길어야 합니다"}
	}

	// 5.1. Daily Event 인 경우 일간 반복 데이터가 있어야 함
	if req.Frequency == Daily {
		if req.Daily == nil || req.Daily.Interval <= 0 {
			return shared.ValidateError{Message: "일 반복 이벤트의 정보가 지정되지 않았습니다"}
		}
	}

	// 5,2. Weekly Event 인 경우 주간 반복 데이터가 있어야 함
	if req.Frequency == Weekly {
		if req.Weekly == nil || len(req.Weekly.Weekdays) == 0 {
			return shared.ValidateError{Message: "주 반복 이벤트의 정보가 지정되지 않았습니다"}
		}
	}
	// 5.3. Monthly Event 인 경우 월 반복 데이터가 있어야 함
	if req.Frequency == Monthly {
		if req.Monthly == nil || len(req.Monthly.Days) == 0 {
			return shared.ValidateError{Message: "월 반복 이벤트의 정보가 지정되지 않았습니다"}
		}
	}
	// 5.4. Yearly Event 인 경우 년 반복 데이터가 있어야 함
	if req.Frequency == Yearly {
		if req.Yearly == nil || req.Yearly.Day == 0 {
			return shared.ValidateError{Message: "년 반복 이벤트의 정보가 지정되지 않았습니다"}
		}
	}

	// 반복 데이터 정리
	switch req.Frequency {
	case Daily:
		req.Weekly = nil
		req.Monthly = nil
		req.Yearly = nil
	case Weekly:
		req.Daily = nil
		req.Monthly = nil
		req.Yearly = nil
	case Monthly:
		req.Daily = nil
		req.Weekly = nil
		req.Yearly = nil
	case Yearly:
		req.Daily = nil
		req.Weekly = nil
		req.Monthly = nil
	}

	return nil
}
