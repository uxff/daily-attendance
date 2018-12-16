package attendance

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/uxff/daily-attendance/lib/modules/attendance/models"
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

type CheckInRuleMap map[string]*CheckInRule

type CheckInRule struct {
	//CheckInKey string `json:"checkinkey"`
	TimeSpan    string `json:"timespan,omitempty"`
	TimeSpanMap struct {
		Start string
		End   string
	} `json:"-"`
	DaySpan    string `json:"dayspan,omitempty"`
	DaySpanMap struct {
		Start  string
		End    string
		StartN int
		EndN   int
	} `json:"-"`
}

func (c *CheckInRuleMap) IsValid(checkInPeriodType int8) bool {
	if c == nil {
		return false
	}

	for _, rule := range *c {
		switch checkInPeriodType {
		case models.CheckInPeriodSecondly:
			rule.DaySpanMap.EndN = 999
			return true
		case models.CheckInPeriodMinutely:
			rule.DaySpanMap.EndN = 60
			if !rule.IsSecondSpanValid() {
				return false
			}
		case models.CheckInPeriodHourly:
			rule.DaySpanMap.EndN = 60
			if !rule.IsMinuteSpanValid() {
				return false
			}
		case models.CheckInPeriodDaily:
			rule.DaySpanMap.EndN = 24
			if !rule.IsHourSpanValid() {
				return false
			}
		case models.CheckInPeriodWeekly:
			rule.DaySpanMap.EndN = 7
			if !rule.IsWeekdaySpanValid() {
				return false
			}
		case models.CheckInPeriodMonthly:
			rule.DaySpanMap.EndN = 31
			if !rule.IsDaySpanValid() {
				return false
			}
		case models.CheckInPeriodYearly:
			rule.DaySpanMap.EndN = 12
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
		// timespan must be as 01:02-03:04
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
			logs.Error("sscanf %s error:%v", timeSpanStart, err)
			return false
		}

		_, err = fmt.Sscanf(timeSpanEnd, "%d:%d", &endHour, &endMin)
		if err != nil {
			logs.Error("sscanf %s error:%v", timeSpanEnd, err)
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
		// dayspan must be as 01-02
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

func (c *CheckInRuleMap) IsInTimeSpan(t time.Time, checkInPeriodType int8) (ruleKey string, rule *CheckInRule) {
	if c == nil {
		return "", nil
	}

	if !c.IsValid(checkInPeriodType) {
		return "", nil
	}

	for cirKey, rule := range *c {
		switch checkInPeriodType {
		case models.CheckInPeriodSecondly:
			return cirKey, rule
		case models.CheckInPeriodMinutely:
			if rule.IsInSecondSpan(t) {
				return cirKey, rule
			}
		case models.CheckInPeriodHourly:
			if rule.IsInMinuteSpan(t) {
				return cirKey, rule
			}
		case models.CheckInPeriodDaily:
			if rule.IsInHourSpan(t) {
				return cirKey, rule
			}
		case models.CheckInPeriodWeekly:
			if rule.IsInWeekdaySpan(t) {
				return cirKey, rule
			}
		case models.CheckInPeriodMonthly:
			if rule.IsInDaySpan(t) {
				return cirKey, rule
			}
		case models.CheckInPeriodYearly:
			if rule.IsInDaySpan(t) {
				return cirKey, rule
			}
		}
	}

	return "", nil
}

func Json2CheckInRule(str string) CheckInRuleMap {
	cir := CheckInRuleMap{}
	err := json.Unmarshal([]byte(str), &cir)
	if err != nil {
		logs.Warn("err when Json2CheckInRule:%v", err)
	}
	return cir
}

type CheckInScheduleElem struct {
	KeyMark string
	Key     string
	From    string
	To      string
	CilId   int
}

// 获取某一时间的key和对应的起始时间
func (c *CheckInRule) GetSecondlyCheckInScheduleElem(jalId int, checkInKeyMark string, t time.Time) *CheckInScheduleElem {
	return &CheckInScheduleElem{
		KeyMark: checkInKeyMark,
		Key:     fmt.Sprintf("%d-%s-%s", jalId, checkInKeyMark, t.Format("20060102150405")),
		From:    fmt.Sprintf("%s", t.Format("2006-01-02 15:04:05")),
		To:      fmt.Sprintf("%s", t.Format("2006-01-02 15:04:05")),
	}
}

func (c *CheckInRule) GetMinutelyCheckInScheduleElem(jalId int, checkInKeyMark string, t time.Time) *CheckInScheduleElem {
	if t.Second() > c.DaySpanMap.EndN {
		t = t.Add(time.Minute)
	}
	return &CheckInScheduleElem{
		KeyMark: checkInKeyMark,
		Key:     fmt.Sprintf("%d-%s-%s", jalId, checkInKeyMark, t.Format("200601021504")),
		From:    fmt.Sprintf("%s:%02d", t.Format("2006-01-02 15:04"), c.DaySpanMap.StartN),
		To:      fmt.Sprintf("%s:%02d", t.Format("2006-01-02 15:04"), c.DaySpanMap.EndN),
	}
}

func (c *CheckInRule) GetHourlyCheckInScheduleElem(jalId int, checkInKeyMark string, t time.Time) *CheckInScheduleElem {
	if t.Format("04:05") > c.TimeSpanMap.End {
		t = t.Add(time.Hour)
	}
	return &CheckInScheduleElem{
		KeyMark: checkInKeyMark,
		Key:     fmt.Sprintf("%d-%s-%s", jalId, checkInKeyMark, t.Format("2006010215")),
		From:    fmt.Sprintf("%s:%s", t.Format("2006-01-02 15"), c.TimeSpanMap.Start),
		To:      fmt.Sprintf("%s:%s", t.Format("2006-01-02 15"), c.TimeSpanMap.End),
	}
}

func (c *CheckInRule) GetDailyCheckInScheduleElem(jalId int, checkInKeyMark string, t time.Time) *CheckInScheduleElem {
	if t.Format("15:04") > c.TimeSpanMap.End {
		t = t.Add(time.Hour * 24)
	}
	return &CheckInScheduleElem{
		KeyMark: checkInKeyMark,
		Key:     fmt.Sprintf("%d-%s-%s", jalId, checkInKeyMark, t.Format("20060102")),
		From:    fmt.Sprintf("%s %s:00", t.Format("2006-01-02"), c.TimeSpanMap.Start),
		To:      fmt.Sprintf("%s %s:59", t.Format("2006-01-02"), c.TimeSpanMap.End),
	}
}

// week day start from 0 // test ok
func (c *CheckInRule) GetWeeklyCheckInScheduleElem(jalId int, checkInKeyMark string, t time.Time) *CheckInScheduleElem {
	if int(t.Weekday()) > c.DaySpanMap.EndN {
		t = t.Add(time.Hour*24*7 - time.Hour*24*time.Duration(int(t.Weekday())-c.DaySpanMap.StartN))
	} else {
		t = t.Add(-time.Hour * 24 * time.Duration(int(t.Weekday())-c.DaySpanMap.StartN))
	}
	return &CheckInScheduleElem{
		KeyMark: checkInKeyMark,
		Key:     fmt.Sprintf("%d-%s-%sw%02d", jalId, checkInKeyMark, t.Format("2006"), t.YearDay()-int(t.Weekday())),
		From:    fmt.Sprintf("%s 00:00:00", t.Format("2006-01-02")),
		To:      fmt.Sprintf("%s 23:59:59", t.Add(time.Hour*24*time.Duration(c.DaySpanMap.EndN-c.DaySpanMap.StartN)).Format("2006-01-02")),
	}
}

// month start from 1
func (c *CheckInRule) GetMonthlyCheckInScheduleElem(jalId int, checkInKeyMark string, t time.Time) *CheckInScheduleElem {
	if t.Day() > c.DaySpanMap.EndN {
		t = t.Add(time.Hour * 24 * 30)
	}
	return &CheckInScheduleElem{
		KeyMark: checkInKeyMark,
		Key:     fmt.Sprintf("%d-%s-%s", jalId, checkInKeyMark, t.Format("200601")),
		From:    fmt.Sprintf("%s-%s 00:00:00", t.Format("2006-01"), c.DaySpanMap.Start),
		To:      fmt.Sprintf("%s-%s 23:59:59", t.Format("2006-01"), c.DaySpanMap.End),
	}
}

// not supported
func (c *CheckInRule) GetYearlyCheckInScheduleElem(jalId int, checkInKeyMark string, t time.Time) *CheckInScheduleElem {
	if int(t.Month()) > c.DaySpanMap.EndN {
		t = t.Add(time.Hour * 24 * 30)
	}
	return &CheckInScheduleElem{
		KeyMark: checkInKeyMark,
		Key:     fmt.Sprintf("%d-%s-%s", jalId, checkInKeyMark, t.Format("2006")),
		From:    fmt.Sprintf("%s-%s-01 00:00:00", t.Format("2006"), c.DaySpanMap.Start),
		To:      fmt.Sprintf("%s-%s-30 23:59:59", t.Format("2006"), c.DaySpanMap.End),
	}
}

// 用于:创建jal;checkin前检查
// 返回：一天可能多条记录 一天一条map
func (c *CheckInRuleMap) GetCheckInScheduleElems(checkInPeriodType int8, jalId int, t time.Time) (elems []*CheckInScheduleElem) {
	if c == nil {
		return nil
	}

	for cirKey, rule := range *c {
		elems = append(elems, rule.GetCheckInScheduleElem(checkInPeriodType, jalId, t, cirKey))
	}

	return elems
}

func (c *CheckInRule) GetCheckInScheduleElem(checkInPeriodType int8, jalId int, t time.Time, cirKey string) *CheckInScheduleElem {
	if c == nil {
		return nil
	}

	switch checkInPeriodType {
	case models.CheckInPeriodSecondly:
		return c.GetSecondlyCheckInScheduleElem(jalId, cirKey, t)
	case models.CheckInPeriodMinutely:
		return c.GetMinutelyCheckInScheduleElem(jalId, cirKey, t)
	case models.CheckInPeriodHourly:
		return c.GetHourlyCheckInScheduleElem(jalId, cirKey, t)
	case models.CheckInPeriodDaily:
		return c.GetDailyCheckInScheduleElem(jalId, cirKey, t)
	case models.CheckInPeriodWeekly:
		return c.GetWeeklyCheckInScheduleElem(jalId, cirKey, t)
	case models.CheckInPeriodMonthly:
		return c.GetMonthlyCheckInScheduleElem(jalId, cirKey, t)
	case models.CheckInPeriodYearly:
		return c.GetYearlyCheckInScheduleElem(jalId, cirKey, t)
	}

	return nil
}

// 因为map乱序，所以说还是用slice
type CheckInSchedules [][]*CheckInScheduleElem

func Json2CheckInSchedules(str string) CheckInSchedules {
	cism := new(CheckInSchedules)

	err := json.Unmarshal([]byte(str), cism)
	if err != nil {
		return nil
		//return fmt.Errorf("error format jal schedules:%s json.Unmarshal error:%v", jal.Schedule, err)
	}
	return *cism
}

func (c *CheckInSchedules) IsTimeIn(t time.Time) (step int, elemIdx int) {
	tstr := t.Format("2006-01-02 15:04:05")
	for step, elems := range *c {
		for i, elem := range elems {
			if elem.From <= tstr && tstr <= elem.To {
				return step, i
			}
		}
	}
	return -1, -1
}

func (c *CheckInSchedules) GetMinMax() (min, max string) {
	for _, elems := range *c {
		for _, elem := range elems {
			if min == "" {
				min = elem.From
			}
			if max == "" {
				max = elem.To
			}
			if min > elem.From {
				min = elem.From
			}
			if max < elem.To {
				max = elem.To
			}
		}
	}
	return min, max
}

func (c *CheckInSchedules) ToJson() string {
	scheduleJson, err := json.Marshal(c)
	if err != nil {
		//return fmt.Errorf("error format jal schedules:%s json.Unmarshal error:%v", jal.Schedule, err)
	}
	return string(scheduleJson)
}
