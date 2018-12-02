package attendance

import (
	"github.com/uxff/daily-attendance/lib/modules/attendance/models"
	"github.com/astaxie/beego/orm"
)

func GetUserBalance(Uid int) *models.UserBalance {
	ormObj := orm.NewOrm()

	ub := models.UserBalance{Uid:Uid}
	ormObj.Read(&ub)

	return &ub
}

func ListUserTradeLog(Uid int) []*models.UserTradeLog {
	list := []*models.UserTradeLog{}

	ormObj := orm.NewOrm()
	ormObj.QueryTable(models.UserTradeLog{}).Filter("uid", Uid).All(&list)

	return list
}

