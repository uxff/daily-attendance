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
	ormObj.QueryTable(models.WastageShare{}).Filter("to_uid", Uid).All(&list)

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
	activities := ListActivities(map[string]interface{}{"status": models.StatusNormal})
	for _, act := range activities {
		actAmount := AccoutingActivityJoined(act.Aid)

		if actAmount <= 0 {
			logs.Warn("no more actvity amount of Aid:%d, ignore", act.Aid)
			continue
		}

		missedJals := ListMissedJal(act.Aid)
		successJals := ListAchievedJal(act.Aid)
		for _, mjal := range missedJals {
			logs.Debug("%d successors will share amount %d from missed jal %d, act amount:%d", len(successJals), mjal.JoinPrice, mjal.JalId, actAmount)
			ShareMissedJal(mjal, successJals, actAmount)
		}
	}
}

//
func ListMissedJal(Aid int) []*models.JoinActivityLog {
	list := []*models.JoinActivityLog{}

	ormObj := orm.NewOrm()
	ormObj.QueryTable(models.JoinActivityLog{}).Filter("aid_id", Aid).Filter("status", models.JalStatusMissed).All(&list)

	return list
}

func ShareMissedJal(missedJal *models.JoinActivityLog, successJals []*models.JoinActivityLog, allJoinedAmounts int) error {
	ormObj := orm.NewOrm()

	moneyWillShare := missedJal.JoinPrice
	if missedJal.JoinPrice <= 0 {
		logs.Warn("missedJal.JoinPrice is empty, no need to share")
		return nil
	}
	//allAchievedFeederGoods := GetAllAchievedGolds()
	missedJal.Status = models.JalStatusShared

	// todo: 么有分享的时候，无需更新missed
	ormObj.Update(missedJal, "status")

	for _, sjal := range successJals {
		oneBonus := moneyWillShare * (sjal.JoinPrice * sjal.Step / sjal.BonusNeedStep / allJoinedAmounts)
		LogShared(missedJal, sjal, oneBonus)
		DispatchBonus(oneBonus, sjal)
	}
	return nil
}

func LogShared(missedJal *models.JoinActivityLog, toJal *models.JoinActivityLog, amount int) {
	ws := models.WastageShare{
		WastedJalId: missedJal.JalId,
		ToJalId:     toJal.JalId,
		FromUid:     missedJal.Uid,
		ToUid:       toJal.Uid,
		Amount:      amount,
	}

	ormObj := orm.NewOrm()
	_, err := ormObj.Insert(&ws)
	if err != nil {
		logs.Error("insert wastage share error:%v", err)
	}
}

func ListAchievedJal(Aid int) []*models.JoinActivityLog {
	list := []*models.JoinActivityLog{}

	ormObj := orm.NewOrm()
	ormObj.QueryTable(models.JoinActivityLog{}).Filter("aid_id", Aid).Filter("status", models.JalStatusAchieved).All(&list)

	return list
}

func GetAchivedAmounts(Aid int) int {
	ormObj := orm.NewOrm()
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return 0
	}

	all := struct {
		JoinPriceAll int
	}{}

	qb.Select("sum(join_price) as join_price_all").From("join_activity_log").
		Where("aid = ?").Limit(1)

	sql := qb.String()

	//logs.Debug("sql=%s", sql)

	err = ormObj.Raw(sql, Aid).QueryRow(&all)
	if err != nil {
		logs.Debug("query error:%v", err)
	}

	return all.JoinPriceAll

}

func DispatchBonus(amount int, successorJal *models.JoinActivityLog) {
	utlId, err := Award(successorJal.Uid, amount, models.TradeTypeCheckInBonus, successorJal.Aid.Name)
	if err != nil {
		logs.Error("dispatch bonus error:%v", err)
		return
	}

	successorJal.BonusTotal += amount

	ormObj := orm.NewOrm()
	ormObj.Update(successorJal, "bonus_total")

	logs.Info("dispatch bonus ok, uid:%d amount:%d jalId:%d utlId:%d", successorJal.Uid, amount, successorJal.JalId, utlId)
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

// all remain = inited+achieved - (missed+shared+stopped)
func AccoutingActivityJoined(Aid int) int {

	ormObj := orm.NewOrm()
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return 0
	}

	allJoined := struct {
		JoinPriceAll int
	}{}

	allMissed := struct {
		JoinPriceAll int
	}{}

	qb.Select("sum(join_price) as join_price_all").From("join_activity_log").
		Where("aid = ? and status in ")

	sql := qb.String()

	//logs.Debug("sql=%s", sql)

	err = ormObj.Raw(sql+fmt.Sprintf("(%d,%d)", models.JalStatusInited, models.JalStatusAchieved), Aid).QueryRow(&allJoined)
	if err != nil {
		logs.Debug("query error:%v", err)
	}

	err = ormObj.Raw(sql+fmt.Sprintf("(%d,%d,%d)", models.JalStatusStopped, models.JalStatusMissed, models.JalStatusShared), Aid).QueryRow(&allMissed)
	if err != nil {
		logs.Debug("query error:%v", err)
	}

	//joined[0]["join_price_all"]
	logs.Debug("--------Aid:%d allJoined:%v allMissed:%v remain:%d", Aid, allJoined, allMissed, allJoined.JoinPriceAll-allMissed.JoinPriceAll)

	return allJoined.JoinPriceAll // - allMissed.JoinPriceAll
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
