package attendance

import (
	"testing"
	"time"
)

func TestCheckInRule(t *testing.T) {
	jalId := 1
	key := "aaa"
	now := time.Now()

	cirm := Json2CheckInRule(`{"s":{"dayspan":""},"i":{"dayspan":"00-10"},"H":{"timespan":"10:00-29:59"},"D":{"timespan":"18:00-23:59"},"W":{"dayspan":"03-04"},"M":{"dayspan":"25-28"},"Y":{"dayspan":"05-12"}}`)
	t.Errorf("cirm=%+v", cirm)

	t.Logf("elem %v S=%+v", cirm["s"].IsSecondSpanValid(), cirm["s"].GetSecondlyCheckInScheduleElem(jalId, key, now))
	t.Logf("elem %v i=%+v", cirm["i"].IsSecondSpanValid(), cirm["i"].GetMinutelyCheckInScheduleElem(jalId, key, now))
	t.Logf("elem %v H=%+v", cirm["H"].IsMinuteSpanValid(), cirm["H"].GetHourlyCheckInScheduleElem(jalId, key, now))
	t.Logf("elem %v D=%+v", cirm["D"].IsHourSpanValid(), cirm["D"].GetDailyCheckInScheduleElem(jalId, key, now))
	t.Logf("elem %v W=%+v", cirm["W"].IsWeekdaySpanValid(), cirm["W"].GetWeeklyCheckInScheduleElem(jalId, key, now))
	t.Logf("elem %v M=%+v", cirm["M"].IsDaySpanValid(), cirm["M"].GetMonthlyCheckInScheduleElem(jalId, key, now))
	t.Logf("elem %v Y=%+v", cirm["Y"].IsDaySpanValid(), cirm["Y"].GetYearlyCheckInScheduleElem(jalId, key, now))
}
