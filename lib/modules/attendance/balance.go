package attendance

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/uxff/daily-attendance/lib/modules/attendance/models"
)

func GetUserBalance(Uid int) *models.UserBalance {
	ormObj := orm.NewOrm()

	ub := models.UserBalance{Uid: Uid}

	filter := ormObj.QueryTable(&models.UserBalance{}).Filter("uid", Uid)
	n, err := filter.Count()
	if err != nil {
		logs.Error("query user(%d) balance error:%v n:%d", Uid, err, n)
		//return nil
	}

	if n > 0 {
	}

	if n <= 0 {
		_, err = ormObj.Insert(&ub)
	} else {
		err = filter.One(&ub)
		//err = ormObj.Read(&ub)// must use pk
	}

	if err != nil {
		logs.Warn("load user(%d) balance error:%v n:%d", Uid, err, n)
	}

	return &ub
}

func ListUserTradeLog(Uid int) []*models.UserTradeLog {
	list := []*models.UserTradeLog{}

	ormObj := orm.NewOrm()
	ormObj.QueryTable(models.UserTradeLog{}).Filter("uid", Uid).All(&list)

	return list
}
