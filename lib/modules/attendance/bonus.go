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


