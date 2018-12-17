package attendance

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/uxff/daily-attendance/lib/modules/attendance/models"
)

func GetUserBalance(Uid int) *models.UserBalance {
	ormObj := orm.NewOrm()

	ub := models.UserBalance{Uid: Uid}

	n, err := ormObj.QueryTable(models.UserBalance{}).Filter("uid", Uid).Count()
	if err != nil {
		logs.Error("query user(%d) balance error:%v", Uid, err)
		return nil
	}
	if n <= 0 {
		ormObj.Insert(&ub)
	} else {
		ormObj.Read(&ub)
	}

	return &ub
}

func ListUserTradeLog(Uid int) []*models.UserTradeLog {
	list := []*models.UserTradeLog{}

	ormObj := orm.NewOrm()
	ormObj.QueryTable(models.UserTradeLog{}).Filter("uid", Uid).All(&list)

	return list
}
