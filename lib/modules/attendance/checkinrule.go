package attendance

import (
	"strings"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"time"
)

const (
	// CheckInKeyType extendible
	CheckInKeyWorkUp = "WORKUP"
	CheckInKeyWorkOff = "WORKOFF"
	CheckInKeyHealthDaily = "HEALTHD"
)

// rule for daily work:[{"WORKUP":{"timespan":"00:00-10:00"]}},{"WORKOFF":{"timespan":"18:00-23:59}}]
// rule for daily health:[{"HEALTH":{"timespan":"06:00-09:00"}}]
// rule for monthly report:[{"REPORTM":{"dayspan":"01-02"}}]

type CheckInRuleMap map[string]CheckInRule

type CheckInRule struct {
	//CheckInKey string `json:"checkinkey"`
	TimeSpan string `json:"timespan,omitempty"`
	TimeSpanMap []struct {
		Start string
		End string
	}
	DaySpan string `json:"dayspan,omitempty"`
	DaySpanMap []struct {
		Start string
		End string
	}
}

func (c *CheckInRuleMap) IsValid() bool {
	if c == nil {
		return false
	}

	for _, rule := range *c {
		if rule.TimeSpan != "" {
			if !rule.IsTimeSpanValid() {
				return false
			}
		}
		if rule.DaySpan != "" {
			if !rule.IsDaySpanValid() {
				return false
			}
		}
	}

	return true
}

func (c *CheckInRule) IsTimeSpanValid()  bool {
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

		c.TimeSpanMap = []struct{
			Start string
			End string
		}{{
			Start:fmt.Sprintf("%02d:%02d", startHour, startMin),
			End:fmt.Sprintf("%02d:%02d", endHour, endMin),
		}}

		return true
	}

	return false
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

		c.DaySpanMap = []struct{
			Start string
			End string
		}{{
			Start:fmt.Sprintf("%02d", dayStart),
			End:fmt.Sprintf("%02d", dayEnd),
		}}

		return true
	}

	return false
}

func (c *CheckInRule) IsInTimeSpan(t time.Time) bool {
	hm := fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
	for _, se := range c.TimeSpanMap {
		if hm <= se.Start || se.End <= hm {
			return false
		}
	}
	return true
}

func (c *CheckInRule) IsInDaySpan(t time.Time) bool {
	day := fmt.Sprintf("%02d", t.Day())
	for _, se := range c.DaySpanMap {
		if day <= se.Start || se.End <= day {
			return false
		}
	}
	return true
}

func Json2CheckInRule(str string) CheckInRuleMap {
	cir := CheckInRuleMap{}
	err := json.Unmarshal([]byte(str), &cir)
	if err != nil {
		logs.Warn("err when Json2CheckInRule:%v", err)
	}
	return cir
}

