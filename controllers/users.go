package controllers

import (
	"html/template"
	"time"

	"github.com/astaxie/beego"
	"github.com/uxff/daily-attendance/lib"
	"github.com/uxff/daily-attendance/lib/modules/attendance"
	models2 "github.com/uxff/daily-attendance/lib/modules/attendance/models"
	"github.com/uxff/daily-attendance/models"
)

type UsersController struct {
	BaseController
}

func (c *UsersController) NestPrepare() {
	// below code will in dead loop
	//if !c.IsLogin {
	//	c.Ctx.Redirect(302, c.LoginPath())
	//	return
	//}
}

// func (c *UsersController) NestFinish() {}

func (c *UsersController) Index() {
	beego.ReadFromRequest(&c.Controller)
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	c.TplName = "users/index.tpl"
}

func (c *UsersController) Login() {

	if c.IsLogin {
		//logs.Debug("is login ?>?????")
		c.Ctx.Redirect(302, c.URLFor("UsersController.Index"))
		return
	}

	c.TplName = "login/login.tpl"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())

	if !c.Ctx.Input.IsPost() {
		// show tpl
		return
	}

	flash := beego.NewFlash()

	if !TheCaptcha.VerifyReq(c.Ctx.Request) {
		flash.Warning("验证码错误")
		flash.Store(&c.Controller)
		return
	}

	if !c.CheckXSRFCookie() {
		flash.Warning("页面过期，请刷新后再试")
		flash.Store(&c.Controller)
		return
	}

	email := c.GetString("Email")
	password := c.GetString("Password")

	user, err := lib.Authenticate(email, password)
	if err != nil || user.Uid < 1 {
		flash.Warning("登录失败，不正确的用户或密码 " + err.Error())
		flash.Store(&c.Controller)
		return
	}

	flash.Success("登录成功")
	flash.Store(&c.Controller)

	c.SetLogin(user)

	c.Redirect(c.URLFor("UsersController.Index"), 303)
}

func (c *UsersController) Logout() {
	c.DelLogin()
	//flash := beego.NewFlash()
	//flash.Success("已成功退出")
	//flash.Store(&c.Controller)

	c.Ctx.Redirect(302, c.URLFor("UsersController.Login"))
}

func (c *UsersController) Signup() {
	c.TplName = "login/signup.tpl"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())

	if !c.Ctx.Input.IsPost() {
		return
	}

	var err error
	flash := beego.NewFlash()

	if !TheCaptcha.VerifyReq(c.Ctx.Request) {
		flash.Warning("验证码错误")
		flash.Store(&c.Controller)
		return
	}

	u := &models.User{}
	if err = c.ParseForm(u); err != nil {
		flash.Error("不合法注册!")
		flash.Store(&c.Controller)
		return
	}
	if err = models.IsValid(u); err != nil {
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}

	u.Lastlogintime = time.Unix(0, 0)
	u.EmailActivated = time.Time{}
	id, err := lib.SignupUser(u)
	if err != nil || id < 1 {
		flash.Warning(err.Error())
		flash.Store(&c.Controller)
		return
	}

	flash.Success("注册成功")
	flash.Store(&c.Controller)

	c.SetLogin(u)

	go func() {
		awardAmount, _ := beego.AppConfig.Int("register_award")
		if awardAmount > 0 {
			attendance.Award(u.Uid, awardAmount, models2.TradeTypeRegisterAward, "")
		}
	}()

	c.Redirect(c.URLFor("UsersController.Index"), 303)
}

func (c *UsersController) Balance() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	balance := attendance.GetUserBalance(c.Userinfo.Uid)
	utls := attendance.ListUserTradeLog(c.Userinfo.Uid)

	c.Data["balance"] = balance
	c.Data["utls"] = utls
	c.Data["total"] = len(utls)

	c.TplName = "users/balance.tpl"
}

func (c *UsersController) Bonus() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	bonusLog := attendance.ListUserBonusLog(c.Userinfo.Uid)
	c.Data["bonusLog"] = bonusLog
	c.Data["total"] = len(bonusLog)

	c.TplName = "users/bonus.tpl"
}

func (c *UsersController) Wastage() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}

	bonusLog := attendance.ListUserWastageLog(c.Userinfo.Uid)
	c.Data["bonusLog"] = bonusLog
	c.Data["total"] = len(bonusLog)

	c.TplName = "users/bonus.tpl"
}

func (c *UsersController) Invite() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.TplName = "users/invite.tpl"
}
