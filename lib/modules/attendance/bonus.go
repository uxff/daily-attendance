package attendance

import (
	"github.com/uxff/daily-attendance/lib/modules/attendance/models"
	"github.com/astaxie/beego/orm"
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
func ShareMissedAttendance(misshedJal *models.JoinActivityLog) {
	// list all activity
	activities := ListActivities(map[string]interface{}{"status":models.StatusNormal })
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
	return nil
}

func ShareMissedJal(missedJal *models.JoinActivityLog, successJals[]*models.JoinActivityLog) error {
	//goldsWillShare := missedJal.JoinUtlId.Price
	//allAchievedFeederGoods := GetAllAchievedGolds()
	//
	// for _, sjal := range successJals {
	//    goods := goldsWillShare * (sjal.JoinUtlId.Price*sjal.Step/allSuccessFeederGoods)
	//    DispatchBonus(sjal.Aidd, sjal.Uid, goods)
	// }
	return nil
}

func ListAchievedJal(Aid int)[]*models.JoinActivityLog {
	return nil
}
