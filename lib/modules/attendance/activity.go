package attendance

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/uxff/daily-attendance/lib/modules/attendance/models"
)

//var ormObj orm.Ormer

func ListActivities(conditions map[string]interface{}) []*models.AttendanceActivity {
	var activities []*models.AttendanceActivity
	filter := orm.NewOrm().QueryTable(&models.AttendanceActivity{Status: models.StatusNormal})
	if conditions != nil {
		for condName, cond := range conditions {
			filter = filter.Filter(condName, cond)
		}
	}
	filter.All(&activities)
	return activities
}

func GetActivity(Aid int) *models.AttendanceActivity {
	var ormObj = orm.NewOrm()
	act := models.AttendanceActivity{Aid: Aid}
	err := ormObj.Read(&act)
	if err != nil {
		logs.Error("load activity(%d) error:%v", Aid, err)
		return nil
	}

	return &act
}

func AddActivity(name string, startTime, endTime time.Time, checkInRule CheckInRuleMap, needStep int, checkInPeriod int8,
	creatorUid int, joinPrice int, awardPerCheckin int) (*models.AttendanceActivity, error) {

	if name == "" {
		return nil, errors.New("name cannot be null")
	}

	if endTime.Unix() <= startTime.Unix() {
		return nil, errors.New("endTime should not smaller than startTime")
	}

	if needStep <= 0 {
		return nil, errors.New("needStep is illegal")
	}

	if !checkInRule.IsValid(checkInPeriod) {
		return nil, fmt.Errorf("checkInRule invalid:%v", checkInRule)
	}

	checkInRuleJson, err := json.Marshal(&checkInRule)
	if err != nil {
		return nil, fmt.Errorf("when json.marshal checkInRule:%v", err)
	}

	act := models.AttendanceActivity{
		Name:            name,
		ValidTimeStart:  startTime.Format("2006-01-02 15:04:05"),
		ValidTimeEnd:    endTime.Format("2006-01-02 15:04:05"),
		CheckInRule:     string(checkInRuleJson),
		CheckInPeriod:   checkInPeriod,
		BonusNeedStep:   needStep,
		JoinPrice:       joinPrice,
		CreatorUid:      creatorUid,
		AwardPerCheckIn: awardPerCheckin,
		Status:          models.StatusNormal,
	}

	var ormObj = orm.NewOrm()
	_, err = ormObj.Insert(&act)
	if err != nil {
		logs.Error("inset jal error:%v", err)
		return nil, err
	}

	return &act, nil
}

func StopActivity(Aid int) {

}

func UserJoinActivity(Aid, Uid, UtlId int) error {

	//if UtlId <= 0 {
	//	//logs.Error("")
	//	return errors.New("UtlId cannot be 0")
	//}

	if Aid <= 0 {
		return errors.New("Aidd cannot be 0")
	}

	var ormObj = orm.NewOrm()
	act := models.AttendanceActivity{Aid: Aid}
	err := ormObj.Read(&act)
	if err != nil {
		logs.Error("cannot find aid(%d) in db: %v", Aid, err)
		return err
	}

	// look up joined and achieving jal

	// add jal
	jal := models.JoinActivityLog{
		Aid:           &act,
		Uid:           Uid,
		IsFinish:      0,
		BonusNeedStep: act.BonusNeedStep,
		JoinUtlId:     UtlId,
		Status:        models.StatusNormal,
		JoinPrice:     act.JoinPrice,
		//Schedulemap:      map[string]string{"a": "b"},
	}
	_, err = ormObj.Insert(&jal)
	if err != nil {
		logs.Error("insert jal error:%v", err)
		return err
	}

	schedules := MakeJalSchedule(&jal)
	jal.StartDate, jal.LastStepDate = schedules.GetMinMax()

	jalSchedule, err := json.Marshal(schedules)
	if err != nil {
		logs.Error("marshal schedule error:%v", err)
		return err
	}

	jal.Schedule = string(jalSchedule)
	_, err = ormObj.Update(&jal, "schedule")
	if err != nil {
		logs.Error("update jal(%d).schedule schedule error:%v", jal.JalId, err)
		return err
	}

	act.JoinedAmount += act.JoinPrice
	act.JoinedUserCount = act.JoinedUserCount + 1

	_, err = ormObj.Update(&act, "joined_amount", "joined_user_count")
	if err != nil {
		logs.Warn("update act(%d).join_price,joined_user_count error:%v", act.Aid, err)
	}
	//

	return nil
}

func ListUserActivityLog(Uid int, Aid int, status []interface{}) []*models.JoinActivityLog {

	list := []*models.JoinActivityLog{}

	ormObj := orm.NewOrm()

	filter := ormObj.QueryTable(models.JoinActivityLog{}).Filter("uid", Uid)
	filter = filter.RelatedSel("aid")
	if Aid > 0 {
		filter = filter.Filter("aid", Aid)
	}
	if len(status) > 0 {
		filter = filter.Filter("status__in", status...)
	}
	filter.All(&list)

	return list
}

func GetJoinActivityLog(JalId int) *models.JoinActivityLog {
	jal := models.JoinActivityLog{JalId: JalId}
	ormObj := orm.NewOrm()
	err := ormObj.QueryTable(jal).Filter("jal_id", JalId).RelatedSel("aid").One(&jal)

	if err != nil {
		logs.Error("cannot find jalId(%d) in db: %v", JalId, err)
		return nil
	}

	return &jal
}

func checkInPeriodToDuration(checkInPeriodType int8) (d time.Duration) {
	switch checkInPeriodType {
	case models.CheckInPeriodSecondly:
		d = time.Second
	case models.CheckInPeriodMinutely:
		d = time.Minute
	case models.CheckInPeriodHourly:
		d = time.Hour
	case models.CheckInPeriodDaily:
		d = time.Hour * 24
	case models.CheckInPeriodWeekly:
		d = time.Hour * 24 * 7
	case models.CheckInPeriodMonthly:
		d = time.Hour * 24 * 30
	case models.CheckInPeriodYearly:
		d = time.Hour * 24 * 365
	}

	return d
}
