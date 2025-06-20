package event

import (
	"cmp"
	"github.com/cho8833/duary_lambda/shared"
	"slices"
	"time"
)

func (service *ServiceImpl) GenerateOccurrence(vo VO, rangeStartDate time.Time, rangeEndDate time.Time) ([]VO, shared.ApplicationError) {

	var occurrences []VO

	if vo.Frequency == Daily {
		occurrences = service.dailyOccurrences(vo, rangeStartDate, rangeEndDate)
	} else if vo.Frequency == Weekly {
		occurrences = service.weeklyOccurrences(vo, rangeStartDate, rangeEndDate)
	} else if vo.Frequency == Monthly {
		occurrences = service.monthlyOccurrences(vo, rangeStartDate, rangeEndDate)
	} else if vo.Frequency == Yearly {
		occurrences = service.yearlyOccurrences(vo, rangeStartDate, rangeEndDate)
	}

	if occurrences == nil {
		occurrences = make([]VO, 0)
	}
	return occurrences, nil
}

func (service *ServiceImpl) dailyOccurrences(event VO, start time.Time, end time.Time) []VO {
	var result []VO

	interval := event.Daily.Interval

	// 반복이 끝나야 하는 최대 시점 계산
	maxRepeatEnd := end
	if event.RecurEndDate != nil && event.RecurEndDate.Before(end) {
		maxRepeatEnd = *event.RecurEndDate
	}
	end = end.Add(time.Second * -1) // 조회 종료일자와 겹치는 경우 포함하지 않음

	recurCount := 1
	for d := *event.RecurStartDate; !d.After(maxRepeatEnd); d = d.AddDate(0, 0, interval) {
		// range 내에 들어가는 날짜만 포함
		if !d.Before(start) && !d.After(end) {
			occurrence := service.changeDate(event, d)
			occurrence.RecurCount = recurCount
			result = append(result, occurrence)
			recurCount++
		}
	}

	return result
}

func (service *ServiceImpl) weeklyOccurrences(event VO, start time.Time, end time.Time) []VO {
	var result []VO
	schedules := event.Weekly.Schedule

	// 반복 종료일 계산
	maxRepeatEnd := end
	if event.RecurEndDate != nil && event.RecurEndDate.Before(end) {
		maxRepeatEnd = *event.RecurEndDate
	}
	end = end.Add(time.Second * -1) // 조회 종료일자와 겹치는 경우 포함하지 않음

	recurCount := 1

	// 반복 시작일을 기준으로 주 단위로 이동
	startOfWeek := start
	for !startOfWeek.After(maxRepeatEnd) {
		for weekday, timeRange := range schedules {
			// 해당 주의 요일에 해당하는 날짜 계산
			dayOffset := int(time.Weekday(weekday) - startOfWeek.Weekday())
			if dayOffset < 0 {
				dayOffset += 7
			}
			eventDate := startOfWeek.AddDate(0, 0, dayOffset)

			if eventDate.Before(*event.RecurStartDate) || eventDate.After(maxRepeatEnd) {
				continue
			}
			if !eventDate.Before(start) && !eventDate.After(end) {
				occurrence := service.changeDate(event, eventDate)
				occurrence = service.changeTime(occurrence, timeRange.Start, timeRange.End)
				occurrence.RecurCount = recurCount
				result = append(result, occurrence)
				recurCount++
			}
		}
		// 다음 interval 주로 이동
		startOfWeek = startOfWeek.AddDate(0, 0, 7)
	}
	// 날짜 정렬
	slices.SortFunc(result, func(i VO, j VO) int {
		return cmp.Compare(i.StartDateTime.Unix(), j.StartDateTime.Unix())
	})

	return result
}

func (service *ServiceImpl) monthlyOccurrences(event VO, start time.Time, end time.Time) []VO {
	var result []VO
	days := event.Monthly.Days

	// 반복 종료일 계산
	maxRepeatEnd := end
	if event.RecurEndDate != nil && event.RecurEndDate.Before(end) {
		maxRepeatEnd = *event.RecurEndDate
	}
	end = end.Add(time.Second * -1) // 조회 종료일자와 겹치는 경우 포함하지 않음

	recurCount := 1

	// 반복 시작 연/월 기준으로 루프
	current := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, start.Location())

	for !current.After(maxRepeatEnd) {
		lastDay := service.lastDayOfMonth(current)

		for _, day := range days {
			d := day
			if d > lastDay {
				d = lastDay
			}

			date := time.Date(current.Year(), current.Month(), d,
				0, 0, 0, 0, start.Location())

			if date.Before(*event.RecurStartDate) || date.After(maxRepeatEnd) {
				continue
			}
			if !date.Before(start) && !date.After(end) {
				occurrence := service.changeDate(event, date)
				occurrence.RecurCount = recurCount
				result = append(result, occurrence)
				recurCount++
			}
		}
		// 다음 interval 개월로 이동
		current = current.AddDate(0, 1, 0)
	}

	// 날짜 정렬
	slices.SortFunc(result, func(i VO, j VO) int {
		return cmp.Compare(i.StartDateTime.Unix(), j.StartDateTime.Unix())
	})
	return result
}

func (service *ServiceImpl) yearlyOccurrences(event VO, start time.Time, end time.Time) []VO {
	var result []VO

	// 반복 종료일 계산
	maxRepeatEnd := end
	if event.RecurEndDate != nil && event.RecurEndDate.Before(end) {
		maxRepeatEnd = *event.RecurEndDate
	}
	end = end.Add(time.Second * -1) // 조회 종료일자와 겹치는 경우 포함하지 않음

	month := event.Yearly.Month
	day := event.Yearly.Day

	recurCount := 1

	// 시작 연도부터 반복
	for year := event.RecurStartDate.Year(); year <= maxRepeatEnd.Year(); year += 1 {
		base := time.Date(year, month, 1, 0, 0, 0, 0, event.RecurStartDate.Location())
		last := service.lastDayOfMonth(base)

		d := day
		if d > last {
			d = last
		}

		date := time.Date(year, month, d, 0, 0, 0, 0, event.RecurStartDate.Location())

		if date.Before(*event.RecurStartDate) || date.After(maxRepeatEnd) {
			continue
		}
		if !date.Before(start) && !date.After(end) {
			occurrence := service.changeDate(event, date)
			occurrence.RecurCount = recurCount
			result = append(result, occurrence)
			recurCount++
		}
	}

	return result
}

func (service *ServiceImpl) changeDate(target VO, date time.Time) VO {
	eventCopy := target
	startDateTime := time.Date(date.Year(), date.Month(), date.Day(), target.StartDateTime.Hour(), target.StartDateTime.Minute(), 0, 0, time.UTC)
	endDateTime := time.Date(date.Year(), date.Month(), date.Day(), target.EndDateTime.Hour(), target.EndDateTime.Minute(), 0, 0, time.UTC)

	if service.isCrossOneDay(startDateTime, endDateTime) {
		endDateTime = endDateTime.AddDate(0, 0, 1)
	}

	eventCopy.StartDateTime = startDateTime
	eventCopy.EndDateTime = endDateTime
	return eventCopy
}

func (service *ServiceImpl) changeTime(target VO, start time.Time, end time.Time) VO {
	eventCopy := target
	startDateTime := time.Date(target.StartDateTime.Year(), target.StartDateTime.Month(), target.StartDateTime.Day(), start.Hour(), start.Minute(), start.Second(), 0, time.UTC)
	endDateTime := time.Date(target.EndDateTime.Year(), target.EndDateTime.Month(), target.EndDateTime.Day(), end.Hour(), end.Minute(), end.Second(), 0, time.UTC)

	if service.isCrossOneDay(startDateTime, endDateTime) {
		endDateTime = endDateTime.AddDate(0, 0, 1)
	}

	eventCopy.StartDateTime = startDateTime
	eventCopy.EndDateTime = endDateTime
	return eventCopy
}

// 해당 월의 마지막 날짜 반환
func (service *ServiceImpl) lastDayOfMonth(t time.Time) int {
	firstOfNextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)
	return lastOfMonth.Day()
}

func (service *ServiceImpl) isCrossOneDay(start time.Time, end time.Time) bool {
	return end.Hour() < start.Hour() || (end.Hour() == start.Hour() && end.Minute() < start.Minute())
}
