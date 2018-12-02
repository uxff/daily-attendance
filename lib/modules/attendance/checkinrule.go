package attendance

import (
	"strings"
	"fmt"
	"strconv"
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
	DaySpan string `json:"dayspan,omitempty"`
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

		return true
	}

	return false
}
