package attendance

import (
	"errors"
	"fmt"
	"time"
	"encoding/json"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"github.com/uxff/daily-attendance/lib/modules/attendance/models"
)

//var ormObj orm.Ormer

func ListActivities(conditions map[string]interface{}) []*models.AttendanceActivity {
	var activities []*models.AttendanceActivity
	filter := orm.NewOrm().QueryTable(&models.AttendanceActivity{Status:models.StatusNormal})
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
	act := models.AttendanceActivity{Aid:Aid}
	err := ormObj.Read(&act)
	if err != nil {
		logs.Error("load activity(%d) error:%v", Aid, err)
		return nil
	}

	return &act
}

func AddActivity(name string, startTime, endTime time.Time, checkInRule CheckInRuleMap, needStep int, checkInPeriod int8,
	creatorUid int, joinPrice int, loserWastagePercent float32) error {

	if name == "" {
		return errors.New("name cannot be null")
	}

	if endTime.Unix() <= startTime.Unix() {
		return errors.New("endTime should not smaller than startTime")
	}

	if needStep <= 0 {
		return errors.New("needStep is illegal")
	}

	if !checkInRule.IsValid(checkInPeriod) {
		return fmt.Errorf("checkInRule invalid:%v", checkInRule)
	}

	checkInRuleJson, err := json.Marshal(&checkInRule)
	if err != nil {
		return fmt.Errorf("when json.marshal checkInRule:%v", err)
	}

	act := models.AttendanceActivity{
		Name:                name,
		ValidTimeStart:      startTime.Format("2006-01-02 15:04:05"),
		ValidTimeEnd:        endTime.Format("2006-01-02 15:04:05"),
		CheckInRule:         string(checkInRuleJson),
		CheckInPeriod:       checkInPeriod,
		BonusNeedStep:       needStep,
		JoinPrice:           joinPrice,
		CreatorUid:          creatorUid,
		LoserWastagePercent: loserWastagePercent,
		Status:              models.StatusNormal,
	}

	var ormObj = orm.NewOrm()
	_, err = ormObj.Insert(&act)
	if err != nil {
		logs.Error("inset jal error:%v", err)
		return err
	}

	return nil
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
	act := models.AttendanceActivity{Aid:Aid}
	err := ormObj.Read(&act)
	if err != nil {
		logs.Error("cannot find aid(%d) in db: %v", Aid, err)
		return err
	}

	// add jal
	jal := models.JoinActivityLog{
		Aid:             &act,
		Uid:              Uid,
		IsFinish:         0,
		RewardDispatched: 0,
		BonusNeedStep:    act.BonusNeedStep,
		JoinUtlId:        UtlId,
		Status:           models.StatusNormal,
	}
	_, err = ormObj.Insert(&jal)
	if err != nil {
		logs.Error("inset jal error:%v", err)
		return err
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

