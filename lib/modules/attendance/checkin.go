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

// for controller, exclusive checkin, not common checkin
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

	// 查看指定好的计划，只用判断计划中时间是否满足
	schedules := make([]*CheckInScheduleElem, 0)
	err := json.Unmarshal([]byte(jal.Schedule), &schedules)
	if err != nil {
		return fmt.Errorf("error format jal schedules:%s json.Unmarshal error:%v", jal.Schedule, err)
	}

	scheduleElemIdx := IsInSchedule(now, schedules)
	if scheduleElemIdx == -1 {
		return fmt.Errorf("now(%v) is not in jal(%d)s schedules", now, jal.JalId)
	}

	scheduleElem := schedules[scheduleElemIdx]

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

	cilId, err := ormObj.Insert(&models.CheckInLog{
		Uid:            Uid,
		Aid:            act.Aid,
		CheckInKeyType: checkInKeyTypeWill,
		CheckInKey:     checkInKeyWill,
		Map:            map[string]string{"a": "c"},
	})

	if err != nil {
		return fmt.Errorf("insert checkin error:%v", err)
	}
	// todo: update user schedule and step
	err = UpdateJalStep(jal, scheduleElem, int(cilId))
	if err != nil {
		return fmt.Errorf("update jal step error:%v", err)
	}

	return nil
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
func MakeJalSchedule(jal *models.JoinActivityLog) []CheckInScheduleElem {
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

func IsInSchedule(t time.Time, elems []*CheckInScheduleElem) int {
	tstr := t.Format("2006-01-02 15:04:05")
	for i, elem := range elems {
		if elem.From <= tstr && tstr <= elem.To {
			return i
		}
	}
	return -1
}

func UpdateJalStep(jal *models.JoinActivityLog, elem *CheckInScheduleElem, CilId int) error {

	ormObj := orm.NewOrm()

	schedules := make([]*CheckInScheduleElem, 0)
	err := json.Unmarshal([]byte(jal.Schedule), &schedules)
	if err != nil {
		return fmt.Errorf("error format jal schedules:%s json.Unmarshal error:%v", jal.Schedule, err)
	}

	for _, e := range schedules {
		if e.KeyMark == elem.KeyMark && e.Key == elem.Key {
			e.CilId = CilId
		}
	}

	scheduleJson, err := json.Marshal(&schedules)
	if err != nil {
		return fmt.Errorf("error format jal schedules:%s json.Unmarshal error:%v", jal.Schedule, err)
	}

	jal.Schedule = string(scheduleJson)

	_, err = ormObj.Update(jal, "schedule")
	if err != nil {
		return fmt.Errorf("error when update jal schedule, jalid:%d error:%v", jal.JalId, err)
	}

	return nil
}
