package attendance

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/uxff/daily-attendance/lib/modules/attendance/models"
)

// for controller, attendance checkin is an exclusive checkin type, not common checkin
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
	nowStr := now.Format("2006-01-02 15:04:05")

	ormObj := orm.NewOrm()

	// 查看指定好的计划，只用判断计划中时间是否满足
	schedules := Json2CheckInSchedules(jal.Schedule)
	if schedules == nil {
		return fmt.Errorf("error format jal schedules:%s", jal.Schedule)
	}

	minTime, maxTime := schedules.GetMinMax()
	if nowStr < minTime {
		return fmt.Errorf("too ealier, jal schedule min:%s max:%s", minTime, maxTime)
	}

	// 当前时间是否在活动的时间规则内，非几天内
	cirm := Json2CheckInRule(jal.Aid.CheckInRule)
	cirKey, cir := cirm.IsInTimeSpan(now, jal.Aid.CheckInPeriod)

	if cir == nil {
		return fmt.Errorf("now is not in activity(%d)s rule(%v)", jal.Aid.Aid, jal.Aid.CheckInRule)
	}

	assumeStep := schedules.EstimateStep(minTime, nowStr, checkInPeriodToDuration(jal.Aid.CheckInPeriod))

	// 本时间点的对应的key
	checkInElem := cir.GetCheckInScheduleElem(jal.Aid.CheckInPeriod, jal.JalId, now, cirKey)
	checkInKeyTypeWill := cirKey
	checkInKeyWill := checkInElem.Key

	logs.Debug("the checkInKeyTypeWill=%s, checkInKeyWill=%s assumeStep=:%d minTime=%s now=%s", checkInKeyTypeWill, checkInKeyWill, assumeStep, minTime, nowStr)

	ormObj.QueryTable(models.CheckInLog{}).Filter("uid", Uid).Filter("aid", act.Aid).Filter("check_in_key", checkInKeyWill).All(&todayCheckInLog)

	if len(todayCheckInLog) > 0 {
		return fmt.Errorf("the checkInKey(%s) has already checked for jal:%d", checkInKeyWill, jal.JalId)
	}

	switch jal.Status {
	case models.JalStatusInited:
		// is going to achieved
		// must be sequence step
		// need step up
		// 当前时间是否在已计划5天时间内 检查5天内
		stepIdx, scheduleElemIdx := schedules.IsTimeIn(now)
		if stepIdx == -1 || scheduleElemIdx == -1 {
			// 如果已经达标，则不再时间段内
			return fmt.Errorf("now is not in jal(%d)s schedules(%v), min:%s max:%s ", jal.JalId, jal.Schedule, minTime, maxTime)
		}

		if jal.Step == stepIdx {
			// its the ok time
			jal.Step++
			if jal.Step >= jal.BonusNeedStep {
				jal.Status = models.JalStatusAchieved
				// todo: notify to calc jal bonus next step
			}
			if jal.Step < assumeStep {
				jal.Status = models.JalStatusMissed
			}

			// save jal, insert checkInLog

			// insert db
			cilId, err := ormObj.Insert(&models.CheckInLog{
				JalId:          jal.JalId,
				Uid:            Uid,
				Aid:            act.Aid,
				CheckInKeyType: checkInKeyTypeWill,
				CheckInKey:     checkInKeyWill,
			})

			if err != nil {
				return fmt.Errorf("insert checkin error:%v", err)
			}

			schedules[stepIdx][scheduleElemIdx].CilId = int(cilId)
			jal.Schedule = schedules.ToJson()

			// update db
			_, err = ormObj.Update(jal, "step", "schedule", "status")
			//err = UpdateJalStep(jal, schedules, stepIdx, int(cilId))
			if err != nil {
				return fmt.Errorf("update jal step error:%v", err)
			}

		}
		logs.Debug("jal:%d jal.Step=%d/%d stepIdx=%d ", jal.JalId, jal.Step, jal.BonusNeedStep, stepIdx)

	case models.JalStatusAchieved:
		// will get bonus
		jal.Step++
		if jal.Step < assumeStep {
			jal.Status = models.JalStatusMissed
		}

		// insert db
		_, err := ormObj.Insert(&models.CheckInLog{
			JalId:          jal.JalId,
			Uid:            Uid,
			Aid:            act.Aid,
			CheckInKeyType: checkInKeyTypeWill,
			CheckInKey:     checkInKeyWill,
		})

		if err != nil {
			return fmt.Errorf("insert checkin error:%v", err)
		}

		_, err = ormObj.Update(jal, "step", "status")
		//err = UpdateJalStep(jal, schedules, stepIdx, int(cilId))
		if err != nil {
			return fmt.Errorf("update jal step error:%v", err)
		}
		logs.Debug("jal:%d jal.Step=%d/%d ", jal.JalId, jal.Step, jal.BonusNeedStep)

	case models.JalStatusStopped, models.JalStatusMissed, models.JalStatusShared:
		return fmt.Errorf("jal(%d) is in unmormal stauts:%v", jal.JalId, models.JalStatusMap[jal.Status])
	}

	return nil
}

func ListUserCheckInLog(Uid int, JalId int, Aid int) []*models.CheckInLog {
	list := []*models.CheckInLog{}

	ormObj := orm.NewOrm()
	filter := ormObj.QueryTable(models.CheckInLog{}).Filter("uid", Uid)
	if JalId > 0 {
		filter = filter.Filter("jal_id", JalId)
	}
	if Aid > 0 {
		filter = filter.Filter("aid", Aid)
	}
	filter.All(&list)
	return list
}

// checkin
// @return map[step][]*CheckInScheduleElem
func MakeJalSchedule(jal *models.JoinActivityLog) CheckInSchedules {
	elemArr := make(CheckInSchedules, 0)
	t := jal.Created
	cirm := Json2CheckInRule(jal.Aid.CheckInRule)

	cirm.IsValid(jal.Aid.CheckInPeriod)
	d := checkInPeriodToDuration(jal.Aid.CheckInPeriod)
	//logs.Debug("after IsValid, cirm=%+v", cirm)

	for i := 0; i < jal.BonusNeedStep; i++ {
		stepElems := cirm.GetCheckInScheduleElems(jal.Aid.CheckInPeriod, jal.JalId, t)
		t = t.Add(d)
		elemArr = append(elemArr, stepElems)
		//elemArr[i] = stepElems
	}

	return elemArr
}
