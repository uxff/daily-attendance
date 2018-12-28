package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

const (
	RoleAdmin = 1
	//RoleManager = 2

	FlagSuperPaid      = 0x800
	FlagPaid           = 8
	FlagRealnamed      = 4
	FlagWechatVerified = 4
	FlagPhoneVerified  = 2
	FlagEmailVerified  = 1
)

type User struct {
	Uid            int       `orm:"pk;auto"`
	Email          string    `orm:"size(64)" form:"Email" valid:"Required;Email"` // unique if registered by emal
	Password       string    `orm:"size(32)" form:"Password" valid:"Required;MinSize(6)"`
	Repassword     string    `orm:"-" form:"Repassword" valid:"Required"`
	Lastlogintime  time.Time `orm:"type(datetime)" form:"-"`
	Created        time.Time `orm:"auto_now_add;type(datetime)"`
	Updated        time.Time `orm:"auto_now;type(datetime)"`
	Status         int       `orm:"type(tinyint);default(1)"`
	EmailActivated time.Time `orm:"type(datetime)"`
	Lastloginip    string    `orm:"size(16);default()"`
	Phone          string    `orm:"size(16);default()"`
	PhoneActivated time.Time `orm:"type(datetime)"`
	Nickname       string    `orm:"size(20);default()"`
	Role           int       `orm:"-"`
	UpstreamUid    *User     `orm:"rel(fk);default(0);null"`
	Openid         string    `orm:"size(32)"`
	WoaId          int       `orm:"type(int);default(0)"`
	WxNickname     string    `orm:"size(32)"`
	WxLogoUrl      string    `orm:"size(256)"`
	//SocialFlag 	int
}

//func (u *User) Id() int64 {
//	return  u.Uid
//}

func (u *User) Valid(v *validation.Validation) {
	if u.Password != u.Repassword {
		v.SetError("Repassword", "Does not matched password, repassword")
	}
}

func (m *User) Insert() error {
	m.EmailActivated = time.Unix(0, 1)
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *User) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) ReadOrCreate(field string, fields ...string) (bool, int64, error) {
	return orm.NewOrm().ReadOrCreate(m, field, fields...)
}

func (m *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (t *User) TableEngine() string {
	return "INNODB"
}

func Users() orm.QuerySeter {
	var table User
	return orm.NewOrm().QueryTable(table).OrderBy("-Uid")
}

func (m *User) IsAdmin() bool {
	return m.Role == RoleAdmin
}

func init() {
	orm.RegisterModelWithPrefix(
		beego.AppConfig.String("dbprefix"),
		new(User))
}
