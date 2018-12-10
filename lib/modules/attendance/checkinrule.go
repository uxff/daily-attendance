package attendance

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/uxff/daily-attendance/lib/modules/attendance/models"
	"strconv"
	"strings"
	"time"
)

const (
	// CheckInKeyType extendible
	CheckInKeyWorkUp      = "WORKUP"
	CheckInKeyWorkOff     = "WORKOFF"
	CheckInKeyHealthDaily = "HEALTHD"
)

// rule for daily work:{"WORKUP":{"timespan":"00:00-10:00"},"WORKOFF":{"timespan":"18:00-23:59"}}
// rule for daily health:[{"HEALTH":{"timespan":"06:00-09:00"}}]
// rule for monthly report:[{"REPORTM":{"dayspan":"01-02"}}]

type CheckInRuleMap map[string]CheckInRule

type CheckInRule struct {
	//CheckInKey string `json:"checkinkey"`
	TimeSpan    string `json:"timespan,omitempty"`
	TimeSpanMap struct {
		Start string
		End   string
	}
	DaySpan    string `json:"dayspan,omitempty"`
	DaySpanMap struct {
		Start  string
		End    string
		StartN int
		EndN   int
	}
}

func (c *CheckInRuleMap) IsValid(checkInPeriodType int8) bool {
	if c == nil {
		return false
	}

	for _, rule := range *c {
		switch checkInPeriodType {
		case models.CheckInPeriodSecondly:
			return true
		case models.CheckInPeriodMinutely:
			if !rule.IsSecondSpanValid() {
				return false
			}
		case models.CheckInPeriodHourly:
			if !rule.IsMinuteSpanValid() {
				return false
			}
		case models.CheckInPeriodDaily:
			if !rule.IsHourSpanValid() {
				return false
			}
		case models.CheckInPeriodWeekly:
			if !rule.IsWeekdaySpanValid() {
				return false
			}
		case models.CheckInPeriodMonthly:
			if !rule.IsDaySpanValid() {
				return false
			}
		}
	}

	return true
}

func (c *CheckInRule) IsSecondSpanValid() bool {
	// use dayspan format
	return c.IsDaySpanValid()
}

func (c *CheckInRule) IsMinuteSpanValid() bool {
	// use timespan format
	return c.IsHourSpanValid()
}

func (c *CheckInRule) IsHourSpanValid() bool {
	if c.TimeSpan != "" {
		timeSpanStartAndEnd := strings.Split(c.TimeSpan, "-")
		if len(timeSpanStartAndEnd) < 2 {
			return false
		}

		timeSpanStart := timeSpanStartAndEnd[0]
		timeSpanEnd := timeSpanStartAndEnd[1]

		startHour, startMin := 0, 0
		endHour, endMin := 0, 0
		_, err := fmt.Sscanf(timeSpanStart, "%d:%d", &startHour, &startMin)
		if err != nil {
			return false
		}

		_, err = fmt.Sscanf(timeSpanEnd, "%d:%d", &endHour, &endMin)
		if err != nil {
			return false
		}

		if startHour*60+startMin > endHour*60+endMin {
			return false
		}

		c.TimeSpanMap = struct {
			Start string
			End   string
		}{
			Start: fmt.Sprintf("%02d:%02d", startHour, startMin),
			End:   fmt.Sprintf("%02d:%02d", endHour, endMin),
		}

		return true
	}

	return true // no limit
}

func (c *CheckInRule) IsWeekdaySpanValid() bool {
	// use timespan format
	return c.IsDaySpanValid()
}

func (c *CheckInRule) IsDaySpanValid() bool {

	if c.DaySpan != "" {
		daySpanStartAndEnd := strings.Split(c.DaySpan, "-")
		if len(daySpanStartAndEnd) < 2 {
			return false
		}

		daySpanStart := daySpanStartAndEnd[0]
		daySpanEnd := daySpanStartAndEnd[1]

		//dayStart, dayEnd := 0, 0
		dayStart, err := strconv.Atoi(daySpanStart)
		if err != nil {
			return false
		}
		dayEnd, err := strconv.Atoi(daySpanEnd)
		if err != nil {
			return false
		}

		if dayStart > dayEnd {
			return false
		}

		c.DaySpanMap = struct {
			Start  string
			End    string
			StartN int
			EndN   int
		}{
			Start:  fmt.Sprintf("%02d", dayStart),
			End:    fmt.Sprintf("%02d", dayEnd),
			StartN: dayStart,
			EndN:   dayEnd,
		}

		return true
	}

	return true // no limit
}

func (c *CheckInRule) IsInSecondSpan(t time.Time) bool {
	hm := fmt.Sprintf("%02d", t.Second())
	if c.DaySpanMap.Start <= hm && hm <= c.DaySpanMap.End {
		return true
	}
	return false
}

func (c *CheckInRule) IsInMinuteSpan(t time.Time) bool {
	hm := fmt.Sprintf("%02d:%02d", t.Minute(), t.Second())
	if c.TimeSpanMap.Start <= hm && hm <= c.TimeSpanMap.End {
		return true
	}
	return false
}

func (c *CheckInRule) IsInHourSpan(t time.Time) bool {
	hm := fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
	if c.TimeSpanMap.Start <= hm && hm <= c.TimeSpanMap.End {
		return true
	}
	return false
}

func (c *CheckInRule) IsInDaySpan(t time.Time) bool {
	hm := fmt.Sprintf("%02d", t.Day())
	if c.DaySpanMap.Start <= hm && hm <= c.DaySpanMap.End {
		return true
	}
	return false
}

func (c *CheckInRule) IsInWeekdaySpan(t time.Time) bool {
	hm := fmt.Sprintf("%02d", t.Weekday())
	if c.DaySpanMap.Start <= hm && hm <= c.DaySpanMap.End {
		return true
	}
	return false
}

func Json2CheckInRule(str string) CheckInRuleMap {
	cir := CheckInRuleMap{}
	err := json.Unmarshal([]byte(str), &cir)
	if err != nil {
		logs.Warn("err when Json2CheckInRule:%v", err)
	}
	return cir
}

// 获取某一时间的key和对应的起始时间
func (c *CheckInRule) GetSecondlyCheckInScheduleKey(jalId int, checkInKeyMark string, t time.Time) CheckInScheduleElem {
	return CheckInScheduleElem{
		KeyMark: checkInKeyMark,
		Key:     fmt.Sprintf("%d-%s-%s", jalId, checkInKeyMark, t.Format("20060102150405")),
		From:    fmt.Sprintf("%s", t.Format("2006-01-02 15:04:05")),
		To:      fmt.Sprintf("%s", t.Format("2006-01-02 15:04:05")),
	}
}

func (c *CheckInRule) GetMinutelyCheckInScheduleKey(jalId int, checkInKeyMark string, t time.Time) CheckInScheduleElem {
	if t.Second() < c.DaySpanMap.StartN {
		t = t.Add(time.Minute)
	}
	return CheckInScheduleElem{
		KeyMark: checkInKeyMark,
		Key:     fmt.Sprintf("%d-%s-%s", jalId, checkInKeyMark, t.Format("200601021504")),
		From:    fmt.Sprintf("%s:%s", t.Format("2006-01-02 15:04"), c.DaySpanMap.Start),
		To:      fmt.Sprintf("%s:%s", t.Format("2006-01-02 15:04"), c.DaySpanMap.End),
	}
}

func (c *CheckInRule) GetHourlyCheckInScheduleKey(jalId int, checkInKeyMark string, t time.Time) CheckInScheduleElem {
	if t.Format("04:05") < c.TimeSpanMap.Start {
		t = t.Add(time.Hour)
	}
	return CheckInScheduleElem{
		KeyMark: checkInKeyMark,
		Key:     fmt.Sprintf("%d-%s-%s", jalId, checkInKeyMark, t.Format("2006010215")),
		From:    fmt.Sprintf("%s:%s", t.Format("2006-01-02 15"), c.TimeSpanMap.Start),
		To:      fmt.Sprintf("%s:%s", t.Format("2006-01-02 15"), c.TimeSpanMap.End),
	}
}

func (c *CheckInRule) GetDailyCheckInScheduleKey(jalId int, checkInKeyMark string, t time.Time) CheckInScheduleElem {
	if t.Format("15:04") < c.TimeSpanMap.Start {
		t = t.Add(time.Hour * 24)
	}
	return CheckInScheduleElem{
		KeyMark: checkInKeyMark,
		Key:     fmt.Sprintf("%d-%s-%s", jalId, checkInKeyMark, t.Format("20060102")),
		From:    fmt.Sprintf("%s %s:00", t.Format("2006-01-02"), c.TimeSpanMap.Start),
		To:      fmt.Sprintf("%s %s:59", t.Format("2006-01-02"), c.TimeSpanMap.End),
	}
}

// week day from 0
func (c *CheckInRule) GetWeeklyCheckInScheduleKey(jalId int, checkInKeyMark string, t time.Time) CheckInScheduleElem {
	if int(t.Weekday()) < c.DaySpanMap.StartN {
		t = t.Add(time.Hour * 24 * 7)
	}
	//tStart := t.YearDay() - int(t.Weekday())
	return CheckInScheduleElem{
		KeyMark: checkInKeyMark,
		Key:     fmt.Sprintf("%d-%s-%sw%02d", jalId, checkInKeyMark, t.Format("2006"), t.YearDay()-int(t.Weekday())),
		From:    fmt.Sprintf("%s 00:00:00", t.Format("2006-01-02"), c.DaySpanMap.Start),
		To:      fmt.Sprintf("%s 23:59:59", t.Format("2006-01-02"), c.DaySpanMap.End),
	}
}

// month from 1
func (c *CheckInRule) GetMonthlyCheckInScheduleKey(jalId int, checkInKeyMark string, t time.Time) CheckInScheduleElem {
	if t.Day() < c.DaySpanMap.StartN {
		t = t.Add(time.Hour * 24)
	}
	return CheckInScheduleElem{
		KeyMark: checkInKeyMark,
		Key:     fmt.Sprintf("%d-%s-%s", jalId, checkInKeyMark, t.Format("200601")),
		From:    fmt.Sprintf("%s-%02d 00:00:00", t.Format("2006-01"), c.DaySpanMap.Start),
		To:      fmt.Sprintf("%s-%02d 23:59:59", t.Format("2006-01"), c.DaySpanMap.End),
	}
}

// not supported
func (c *CheckInRule) GetYearlyCheckInScheduleKey(jalId int, checkInKeyMark string, t time.Time) CheckInScheduleElem {
	if int(t.Month()) < c.DaySpanMap.StartN {
		t = t.Add(time.Hour * 24 * 30)
	}
	return CheckInScheduleElem{
		KeyMark: checkInKeyMark,
		Key:     fmt.Sprintf("%d-%s-%s", jalId, checkInKeyMark, t.Format("2006")),
		From:    fmt.Sprintf("%s-%02d-01 00:00:00", t.Format("2006"), c.DaySpanMap.Start),
		To:      fmt.Sprintf("%s-%02d-30 23:59:59", t.Format("2006"), c.DaySpanMap.End),
	}
}

func (c *CheckInRule) GetSecondlyCheckInScheduleElem(jalId int, checkInMark string, startTime time.Time, step int) string {

	return ""
}
