package attendance

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/uxff/daily-attendance/lib/modules/attendance/models"
)

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
		missedJals := ListMissedJal(act.Aid)
		successJals := ListAchievedJal(act.Aid)
		for _, mjal := range missedJals {
			ShareMissedJal(mjal, successJals)
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

func ShareMissedJal(missedJal *models.JoinActivityLog, successJals []*models.JoinActivityLog) error {
	ormObj := orm.NewOrm()

	moneyWillShare := missedJal.JoinPrice
	// todo: 放在循环外
	allJoinedAmounts := GetAchivedAmounts(missedJal.Aid.Aid)
	//allAchievedFeederGoods := GetAllAchievedGolds()
	missedJal.Status = models.JalStatusShared

	// todo: 么有分享的时候，无需更新missed
	ormObj.Update(missedJal, "status")

	for _, sjal := range successJals {
		// todo: log in wastage share
		oneBonus := moneyWillShare * (sjal.JoinPrice * sjal.Step / allJoinedAmounts)
		DispatchBonus(sjal.Uid, oneBonus, sjal)
	}
	return nil
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

	logs.Debug("sql=%s", sql)

	err = ormObj.Raw(sql, Aid).QueryRow(&all)
	if err != nil {
		logs.Debug("query error:%v", err)
	}

	return all.JoinPriceAll

}

func DispatchBonus(Uid int, amount int, jal *models.JoinActivityLog) {
	utlId, err := Charge(Uid, amount, jal.Aid.Name)
	if err != nil {
		logs.Error("dispatch bonus error:%v", err)
		return
	}

	jal.BonusTotal += amount

	ormObj := orm.NewOrm()
	ormObj.Update(jal, "bonus_total")

	logs.Info("dispatch bonus ok, uid:%d amount:%d jalId:%d utlId:%d", Uid, amount, jal.JalId, utlId)
}

func AutoAccounting() {
	for {
		time.Sleep(time.Minute)
		ShareMissedAttendance()
	}
}
