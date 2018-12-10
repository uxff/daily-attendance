package attendance

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/uxff/daily-attendance/lib/modules/attendance/models"
)

// for controller
func UserCheckIn(Uid int, jal *models.JoinActivityLog) error {
	if jal == nil {
		return errors.New("jal model cannot be null when UserCheckIn")
	}
	var act *models.AttendanceActivity = jal.Aid
	if act == nil {
		return errors.New("act model cannot be null when UserCheckIn")
	}

	todayCheckInLog := []*models.CheckInLog{}
	now := time.Now()

	ormObj := orm.NewOrm()
	today, _ := time.Parse("2006-01-02 15:04:05", now.Format("2006-01-02 00:00:00"))

	logs.Debug("today=%v", today)

	cirMap := Json2CheckInRule(act.CheckInRule)

	// find out possible key
	checkInKeyWill := ""
	checkInKeyTypeWill := ""

	switch act.CheckInPeriod {
	case models.CheckInPeriodSecondly:
		for cik, _ := range cirMap {
			checkInKeyTypeWill = cik
			checkInKeyWill = fmt.Sprintf("%s-%d", cik, now.Unix())
			break
		}

	case models.CheckInPeriodMinutely:
		for cik, rule := range cirMap {
			if rule.IsInSecondSpan(now) {
				// not in checkIn timespan
				checkInKeyTypeWill = cik
				checkInKeyWill = fmt.Sprintf("%s-%04d%02d%02d-%02d%02d", cik,
					now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
				break
			}
		}

	case models.CheckInPeriodHourly:
		for cik, rule := range cirMap {
			if rule.IsInMinuteSpan(now) {
				// not in checkIn timespan
				checkInKeyTypeWill = cik
				checkInKeyWill = fmt.Sprintf("%s-%04d%02d%02d-%02d", cik,
					now.Year(), now.Month(), now.Day(), now.Hour())
				break
			}
		}

	case models.CheckInPeriodDaily:
		for cik, rule := range cirMap {
			if rule.IsInHourSpan(now) {
				// not in checkIn timespan
				checkInKeyTypeWill = cik
				checkInKeyWill = fmt.Sprintf("%s-%04d%02d%02d", cik,
					now.Year(), now.Month(), now.Day())
				break
			}
		}
	case models.CheckInPeriodWeekly:
		for cik, rule := range cirMap {
			if rule.IsInWeekdaySpan(now) {
				// not in checkIn timespan
				checkInKeyTypeWill = cik
				checkInKeyWill = fmt.Sprintf("%s-%04dw%02d", cik,
					now.Year(), now.YearDay()-int(now.Weekday()))
				break
			}
		}
	case models.CheckInPeriodMonthly:
		for cik, rule := range cirMap {
			if rule.IsInDaySpan(time.Now()) {
				// not in checkIn timespan
				checkInKeyTypeWill = cik
				checkInKeyWill = fmt.Sprintf("%s-%04dm%02d", cik,
					now.Year(), now.Month())
				break
			}
		}
	}

	logs.Debug("the checkInKeyTypeWill=%s, checkInKeyWill=%s", checkInKeyTypeWill, checkInKeyWill)

	if checkInKeyWill == "" {
		return fmt.Errorf("not in checkIn timespan, act.checkInRule:%s", act.CheckInRule)
	}

	ormObj.QueryTable(models.CheckInLog{}).Filter("uid", Uid).Filter("aid", act.Aid).Filter("check_key", checkInKeyWill).All(&todayCheckInLog)

	if len(todayCheckInLog) > 0 {
		return fmt.Errorf("the checkInKey is checked:%s", checkInKeyWill)
	}

	_, err := ormObj.Insert(&models.CheckInLog{
		Uid:            Uid,
		Aid:            act.Aid,
		CheckInKeyType: checkInKeyTypeWill,
		CheckInKey:     checkInKeyWill,
	})

	return err
}

func ListUserCheckInLog(Uid int, Aid int) []*models.CheckInLog {
	list := []*models.CheckInLog{}

	ormObj := orm.NewOrm()
	filter := ormObj.QueryTable(models.CheckInLog{}).Filter("uid", Uid)
	if Aid > 0 {
		filter.Filter("aid", Aid)
	}
	filter.All(&list)
	return list
}

type CheckInScheduleElem struct {
	KeyMark string
	Key     string
	From    string
	To      string
}

// checkin
func GetJalSchedule(JalId int) {

}
