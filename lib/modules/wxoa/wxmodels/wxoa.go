package wxmodels

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

func GetWoa(woaId int) *WechatOfficalAccounts {

	woa := WechatOfficalAccounts{WoaId:woaId}

	err := orm.NewOrm().QueryTable(&woa).One(&woa)
	if err != nil {
		logs.Warn("load woa(%d) error:%v", err)
		return nil
	}

	return &woa
}