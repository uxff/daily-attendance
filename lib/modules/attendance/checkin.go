package attendance

import (
	"errors"
	"fmt"
	"time"

	"encoding/json"
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

	schedules := make([]*CheckInScheduleElem, 0)
	err := json.Unmarshal([]byte(jal.Schedule), &schedules)
	if err != nil {
		return fmt.Errorf("error format jal schedules:%s json.Unmarshal error:%v", jal.Schedule, err)
	}

	scheduleElem := IsInSchedule(now, schedules)
	if scheduleElem == nil {
		return fmt.Errorf("now(%v) is not in jal(%d)s schedules", now, jal.JalId)
	}

	checkInKeyTypeWill := scheduleElem.KeyMark
	checkInKeyWill := scheduleElem.Key

	logs.Debug("the checkInKeyTypeWill=%s, checkInKeyWill=%s elem=%+v", checkInKeyTypeWill, checkInKeyWill, scheduleElem)

	if checkInKeyWill == "" {
		return fmt.Errorf("not in checkIn timespan, act.checkInRule:%s", act.CheckInRule)
	}

	ormObj.QueryTable(models.CheckInLog{}).Filter("uid", Uid).Filter("aid", act.Aid).Filter("check_in_key", checkInKeyWill).All(&todayCheckInLog)

	if len(todayCheckInLog) > 0 {
		return fmt.Errorf("the checkInKey is already checked:%s", checkInKeyWill)
	}

	_, err = ormObj.Insert(&models.CheckInLog{
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

// checkin
func GetJalSchedule(jal *models.JoinActivityLog) []CheckInScheduleElem {
	elemArr := make([]CheckInScheduleElem, 0)
	t := jal.Created
	cirm := Json2CheckInRule(jal.Aid.CheckInRule)

	cirm.IsValid(jal.Aid.CheckInPeriod)
	//logs.Debug("after IsValid, cirm=%+v", cirm)

	for i := 0; i < jal.BonusNeedStep; i++ {
		d, stepElems := cirm.GetCheckInScheduleElems(jal.Aid.CheckInPeriod, jal.JalId, t)
		t = t.Add(d)
		elemArr = append(elemArr, stepElems...)
	}

	return elemArr
}

func IsInSchedule(t time.Time, elems []*CheckInScheduleElem) *CheckInScheduleElem {
	tstr := t.Format("2006-01-02 15:04:05")
	for _, elem := range elems {
		if elem.From <= tstr && tstr <= elem.To {
			return elem
		}
	}
	return nil
}
