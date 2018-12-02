package attendance

import (
	"fmt"
	"errors"
	"time"

	"github.com/uxff/daily-attendance/lib/modules/attendance/models"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

// for controller
func UserCheckIn(Uid int, act *models.AttendanceActivity) error {
	if act == nil {
		return errors.New("act cannot be null when UserCheckIn")
	}

	todayCheckInLog := []*models.CheckInLog{}

	ormObj := orm.NewOrm()
	today, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 00:00:00"))

	logs.Debug("today=%v", today)

	cirMap := Json2CheckInRule(act.CheckInRule)

	checkInKeyWill := ""
	checkInKeyTypeWill := ""

	switch act.CheckInPeriod {
	case models.CheckInPeriodHourly:

	case models.CheckInPeriodDaily:
		for cik, rule := range cirMap {
			if rule.TimeSpan == "" {
				//return errors.New("act has no TimeSpan when its CheckInPeriodDaily")
				continue
			}
			if rule.IsInTimeSpan(time.Now()) {
				// not in checkIn timespan
				//return errors.New()
				checkInKeyTypeWill = cik
				checkInKeyWill = fmt.Sprintf("%s-%d%d", cik, time.Now().Year(), time.Now().Day())
				break
			}
		}
	case models.CheckInPeriodMonthly:
	}

	if checkInKeyWill == "" {
		return fmt.Errorf("not in checkIn timespan, act.checkInRule:%s", act.CheckInRule)
	}

	ormObj.QueryTable(models.CheckInLog{}).Filter("uid", Uid).Filter("check_key", checkInKeyWill).All(&todayCheckInLog)

	if len(todayCheckInLog) > 0 {
		return fmt.Errorf("the checkInKey is checked:%s", checkInKeyWill)
	}

	_, err := ormObj.Insert(&models.CheckInLog{
		Uid:Uid,
		Aid:act.Aid,
		CheckInKeyType:checkInKeyTypeWill,
		CheckInKey:checkInKeyWill,
	})

	return err
}

func ListUserCheckInLog(Uid int) []*models.CheckInLog {
	list := []*models.CheckInLog{}

	ormObj := orm.NewOrm()
	ormObj.QueryTable(models.CheckInLog{}).Filter("uid", Uid).All(&list)
	return list
}

// checkin


