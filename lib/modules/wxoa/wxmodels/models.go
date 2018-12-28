package wxmodels

import "time"

type WechatOfficalAccounts struct {
	WoaId          int       `orm:"pk;auto"`
	Name           string    `orm:"size(32)"`
	Appid          string    `orm:"size(32)"`
	Appsecret      string    `orm:"size(64)"`
	Token          string    `orm:"size(32)"`
	EncodingAesKey string    `orm:"size(64)"`
	OriginId       string    `orm:"size(32)"`
	Created        time.Time `orm:"auto_now_add;type(datetime)"`
	Updated        time.Time `orm:"auto_now;type(datetime)"`
	Status         int       `orm:"type(tinyint);default(1)"`
}
