package attendance

import (
	"time"

	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/uxff/daily-attendance/lib/modules/attendance/models"
)

var accountingPeriodSec = 60

func ListUserBonusLog(Uid int) []*models.WastageShare {
	list := []*models.WastageShare{}

	ormObj := orm.NewOrm()
	_, err := ormObj.QueryTable(models.WastageShare{}).Filter("to_uid", Uid).RelatedSel("aid").All(&list)
	if err != nil {
		logs.Warn("query WastageShare(toUid:%d) error: %v", Uid, err)
	}

	return list
}

func ListUserWastageLog(Uid int) []*models.WastageShare {
	list := []*models.WastageShare{}

	ormObj := orm.NewOrm()
	ormObj.QueryTable(models.WastageShare{}).Filter("from_uid", Uid).All(&list)

	return list
}

// share all missed
func ShareMissedAttendance() {
	// list all activity
	ormObj := orm.NewOrm()
	activities := ListActivities(map[string]interface{}{"status": models.StatusNormal})
	for _, act := range activities {
		allJoinedAmount := AccoutingActivityByStatus(act.Aid, nil)

		if allJoinedAmount <= 0 {
			logs.Warn("no more joined amount of Aid:%d, ignore", act.Aid)
			continue
		}

		if allJoinedAmount != act.JoinedAmount {
			act.JoinedAmount = allJoinedAmount
			_, err := ormObj.Update(act, "joined_amount")
			if err != nil {
				logs.Error("update act(%d).joined_amount error:%v", act.Aid, err)
			}

		}

		missedJals := ListMissedJal(act.Aid)
		successJals := ListAchievedJal(act.Aid)

		if len(successJals) == 0 {
			logs.Warn("no more success jals of aid:%d, ignore", act.Aid)
			continue
		}

		if len(missedJals) == 0 {
			logs.Warn("no more missed jals of aid:%d, ignore", act.Aid)
			continue
		}

		for _, mjal := range missedJals {
			logs.Debug("%d successors will share amount %d from missed jal %d, act amount:%d", len(successJals), mjal.JoinPrice, mjal.JalId, act.JoinedAmount)
			ShareMissedJal(mjal, successJals, act.JoinedAmount, act)
		}

		act.UnsharedAmount = AccoutingActivityByStatus(act.Aid, []int8{models.JalStatusMissed, models.JalStatusStopped})
		act.SharedAmount = AccoutingActivityByStatus(act.Aid, []int8{models.JalStatusShared})
		act.AllMissedAmount = act.UnsharedAmount + act.SharedAmount
		act.MissedUserCount = len(missedJals)
		_, err := ormObj.Update(act, "shared_amount", "unshared_amount", "all_missed_amount", "missed_user_count")
		if err != nil {
			logs.Error("update act(%d).shared_amount error:%v", act.Aid, err)
		}

	}
}

func ListAchievedJal(Aid int) []*models.JoinActivityLog {
	list := []*models.JoinActivityLog{}

	ormObj := orm.NewOrm()
	ormObj.QueryTable(models.JoinActivityLog{}).Filter("aid_id", Aid).Filter("status", models.JalStatusAchieved).All(&list)

	return list
}

//
func ListMissedJal(Aid int) []*models.JoinActivityLog {
	list := []*models.JoinActivityLog{}

	ormObj := orm.NewOrm()
	ormObj.QueryTable(models.JoinActivityLog{}).Filter("aid_id", Aid).Filter("status", models.JalStatusMissed).All(&list)

	return list
}

func ShareMissedJal(missedJal *models.JoinActivityLog, successJals []*models.JoinActivityLog, allJoinedAmounts int, act *models.AttendanceActivity) error {
	ormObj := orm.NewOrm()

	if missedJal.JoinPrice <= 0 {
		logs.Warn("missedJal.JoinPrice is empty, no need to share")
		return nil
	}
	//allAchievedFeederGoods := GetAllAchievedGolds()
	missedJal.Status = models.JalStatusShared

	ormObj.Update(missedJal, "status")

	for _, sjal := range successJals {
		oneBonus := missedJal.JoinPrice * (sjal.JoinPrice * sjal.Step / sjal.BonusNeedStep / allJoinedAmounts)
		SharedToOne(missedJal, sjal, oneBonus, act)
		//DispatchBonus(oneBonus, sjal)
	}

	return nil
}

func SharedToOne(missedJal *models.JoinActivityLog, toJal *models.JoinActivityLog, amount int, act *models.AttendanceActivity) {
	utlId, err := Award(toJal.Uid, amount, models.TradeTypeCheckInBonus, act.Name)
	if err != nil {
		logs.Error("dispatch bonus to jal(%d) error:%v", toJal.JalId, err)
		return
	}

	ws := models.WastageShare{
		WastedJalId: missedJal.JalId,
		ToJalId:     toJal.JalId,
		FromUid:     missedJal.Uid,
		ToUid:       toJal.Uid,
		Amount:      amount,
		Aid:         act,
		UtlId:       utlId,
	}

	ormObj := orm.NewOrm()
	_, err = ormObj.Insert(&ws)
	if err != nil {
		logs.Error("insert wastage share error:%v", err)
		return
	}

	toJal.BonusTotal += amount

	_, err = ormObj.Update(toJal, "bonus_total")
	if err != nil {
		logs.Error("update successor jal(%d).bonus_total error:%v", toJal.JalId, err)
		return
	}

	logs.Info("dispatch bonus ok, uid:%d bonus:%d jalId:%d joined:%d allAmount:%d utlId:%d", toJal.Uid, amount, toJal.JalId, toJal.JoinPrice, act.JoinedAmount, utlId)
}

func StopAllUnachiedJal() {
	ormObj := orm.NewOrm()
	jals := make([]*models.JoinActivityLog, 0)

	ormObj.QueryTable(models.JoinActivityLog{}).RelatedSel("aid").Filter("status__in", []interface{}{models.JalStatusInited, models.JalStatusAchieved}).All(&jals)

	now := time.Now()

	for _, jal := range jals {
		schedules := Json2CheckInSchedules(jal.Schedule)
		min, _ := schedules.GetMinMax()
		stepIdx := schedules.EstimateStep(min, now.Format("2006-01-02 15:04:05"), checkInPeriodToDuration(jal.Aid.CheckInPeriod))
		if (stepIdx - jal.Step) > 1 {
			// no
			logs.Debug("jal:%d db step:%d estimate step:%d will be missed", jal.JalId, jal.Step, stepIdx)
			jal.Status = models.JalStatusMissed
			ormObj.Update(jal, "status")
		}
	}
}

// all joined: status(all)
// all unshared: status(missed,stopped)
// all shared: status(shared)
func AccoutingActivityByStatus(Aid int, status []int8) (sum int) {

	ormObj := orm.NewOrm()
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		logs.Error("orm.NewQueryBuilder error:%v", err)
		return
	}

	theSum := struct {
		Thesum int
		//TheCount int
	}{}

	sql := qb.Select("sum(join_price) as thesum").From("join_activity_log").
		Where("aid_id = ? ").String()

	sts := ""
	for _, s := range status {
		sts += fmt.Sprintf(",%d", s)
	}

	if len(sts) > 1 {
		sts = sts[1:]
		sql = sql + " and status in (" + sts + ")"
	}

	err = ormObj.Raw(sql, Aid).QueryRow(&theSum)
	if err != nil {
		logs.Warn("query(%s) error:%v", sql, err)
	}

	logs.Debug("--------Aid:%d sum:%d sts:%s", Aid, theSum.Thesum, sts)

	return theSum.Thesum //, theSum.TheCount
}

func AutoAccounting() {
	for {
		time.Sleep(time.Second * time.Duration(accountingPeriodSec))
		// stop unachieved jal
		StopAllUnachiedJal()
		// share the missed
		ShareMissedAttendance()
		// todo: accounting each activity all joined money

	}
}

func SetAccountingPeriod(nSec int) {
	accountingPeriodSec = nSec
}
