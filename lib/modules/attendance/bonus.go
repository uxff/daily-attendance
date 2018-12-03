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

func GetUserBonus(Uid int) {

}

func ListUserWastageLog(Uid int) []*models.WastageShare {
	list := []*models.WastageShare{}

	ormObj := orm.NewOrm()
	ormObj.QueryTable(models.WastageShare{}).Filter("from_uid", Uid).All(&list)

	return list
}

// share all missed
func ShareMissedAttendance(jal *models.JoinActivityLog) {
	// list all activity
	activities := ListActivities()
	for _, act := range activities {
		missedJals := ListMissedJal(act.Aid)
		successJals := ListSuccessJal(act.Aid)
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
	//allSuccessFeederGoods := GetAllSuccessorGolds()
	//
	// for _, sjal := range successJals {
	//    goods := goldsWillShare * (sjal.JoinUtlId.Price*sjal.Step/allSuccessFeederGoods)
	//    DispatchBonus(sjal.Aid, sjal.Uid, goods)
	// }
	return nil
}

func ListSuccessJal(Aid int)[]*models.JoinActivityLog {
	return nil
}
