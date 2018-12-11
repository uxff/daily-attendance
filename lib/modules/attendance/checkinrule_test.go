package attendance

import (
	"testing"
	"time"
)

func TestCheckInRule(t *testing.T) {
	jalId := 1
	key := "aaa"
	now := time.Now()

	cirm := Json2CheckInRule(`{"s":{"dayspan":""},"i":{"dayspan":"00-10"},"H":{"timespan":"10:00-29:59"},"D":{"timespan":"18:00-23:59"},"M":{"dayspan":"01-02"},"W":{"dayspan":"03-04"},"Y":{"dayspan":"05-12"}}`)
	t.Errorf("cirm=%+v", cirm)

	t.Logf("elem S=%+v", cirm["s"].GetSecondlyCheckInScheduleElem(jalId, key, now))
	t.Logf("elem i=%+v", cirm["i"].GetMinutelyCheckInScheduleElem(jalId, key, now))
	t.Logf("elem H=%+v", cirm["H"].GetHourlyCheckInScheduleElem(jalId, key, now))
	t.Logf("elem D=%+v", cirm["D"].GetDailyCheckInScheduleElem(jalId, key, now))
	t.Logf("elem M=%+v", cirm["M"].GetMonthlyCheckInScheduleElem(jalId, key, now))
	t.Logf("elem W=%+v", cirm["W"].GetWeeklyCheckInScheduleElem(jalId, key, now))
	t.Logf("elem Y=%+v", cirm["Y"].GetYearlyCheckInScheduleElem(jalId, key, now))
}
