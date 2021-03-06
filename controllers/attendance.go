package controllers

import (
	"encoding/json"
	"html/template"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/uxff/daily-attendance/lib/modules/attendance"
	"github.com/uxff/daily-attendance/lib/modules/attendance/models"
)

func init() {
	beego.AddFuncMap("jalstatus", func(jalStatus int8) string {
		return models.JalStatusMap[jalStatus]
	})
	beego.AddFuncMap("checkinperiod", func(v int8) string {
		return models.CheckInPeriodMap[v]
	})
	beego.AddFuncMap("tradetype", func(v int8) string {
		return models.TradeTypeMap[v]
	})
}

type AttendanceController struct {
	BaseController
}

func (c *AttendanceController) NestPrepare() {
}

// func (c *UsersController) NestFinish() {}

func (c *AttendanceController) Index() {

	activities := attendance.ListActivities(map[string]interface{}{"status": models.StatusNormal})
	c.Data["activities"] = activities

	c.TplName = "attendance/index.tpl"
}

func (c *AttendanceController) My() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.URLFor("UsersController.Index"))
		return
	}

	// 查看已经参与的活动
	jals := attendance.ListUserActivityLog(c.Userinfo.Uid, 0, nil)

	c.Data["jals"] = jals
	c.Data["total"] = len(jals)
	c.Data["jalStatusMap"] = models.JalStatusMap
	c.TplName = "attendance/my.tpl"
}

func (c *AttendanceController) Join() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.URLFor("UsersController.Index"))
		return
	}

	c.TplName = "attendance/join.tpl"

	aid, _ := c.GetInt("aid")

	flash := beego.NewFlash()
	if aid == 0 {
		flash.Warning("没有指定aid")
		flash.Store(&c.Controller)
		return
	}

	act := attendance.GetActivity(aid)
	if act == nil {
		flash.Warning("aid不存在")
		flash.Store(&c.Controller)
		return
	}

	// 查看已经参与的活动
	jals := attendance.ListUserActivityLog(c.Userinfo.Uid, aid, []interface{}{models.JalStatusAchieved, models.JalStatusInited})

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["act"] = act
	c.Data["jals"] = jals

	if !c.Ctx.Input.IsPost() {
		return
	}

	if !c.CheckXSRFCookie() {
		flash.Warning("页面过期，请刷新后重试")
		flash.Store(&c.Controller)
		return
	}

	if len(jals) > 0 {
		flash.Warning("您已经参与过活动%s(%d次)", act.Name, len(jals))
		flash.Store(&c.Controller)
		return
	}

	utlId := 0
	var err error
	if act.JoinPrice > 0 {
		actProduct := attendance.ActivityToProduct(act)
		utlId, err = attendance.Consume(c.Userinfo.Uid, actProduct, 1, "参与活动:"+actProduct.GetName())
		if err != nil {
			flash.Warning("参与失败：%v", err)
			flash.Store(&c.Controller)
			return
		}
		logs.Warn("交易成功:utlId:%d", utlId)
	}

	err = attendance.UserJoinActivity(aid, c.Userinfo.Uid, utlId)
	if err != nil {
		flash.Warning("参与活动%d失败：%v", aid, err)
		flash.Store(&c.Controller)
		return
	}

	flash.Success("参与活动%s(%d)成功", act.Name, aid)
	flash.Store(&c.Controller)
	c.Redirect(c.Ctx.Request.RequestURI, 303)
	//c.Redirect(c.URLFor("AttendanceController.Join", "aid", act.Aid), 303)
}

func (c *AttendanceController) Add() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.URLFor("UsersController.Index"))
		return
	}

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())

	c.TplName = "attendance/add.tpl"
	if !c.Ctx.Input.IsPost() {
		// show tpl
		return
	}

	flash := beego.NewFlash()

	if !c.CheckXSRFCookie() {
		flash.Warning("页面过期，请刷新后重试")
		flash.Store(&c.Controller)
		return
	}

	name := c.GetString("activity_name")
	startTimeStr := c.GetString("startTime")
	endTimeStr := c.GetString("endTime")
	needStep, _ := c.GetInt("needStep")
	checkInRuleStr := c.GetString("checkInRule")
	joinPrice, _ := c.GetInt("joinPrice")
	awardPerCheckin, _ := c.GetInt("awardPerCheckin")
	checkInPeriod, _ := c.GetInt8("checkInPeriod")

	startTime, err := time.Parse("2006-01-02 15:04", startTimeStr)
	endTime, err := time.Parse("2006-01-02 15:04", endTimeStr)

	checkInRuleMap := new(attendance.CheckInRuleMap)
	err = json.Unmarshal([]byte(checkInRuleStr), checkInRuleMap)
	if err != nil {
		flash.Warning("规则格式不正确 " + err.Error() + " origin:" + checkInRuleStr)
		flash.Store(&c.Controller)
		return
	}
	logs.Warn("will create activity:%v %v %v %v %v %v", name, startTime, endTime, needStep, checkInRuleMap, awardPerCheckin)

	act, err := attendance.AddActivity(name, startTime, endTime, *checkInRuleMap, needStep, checkInPeriod, c.Userinfo.Uid, joinPrice, awardPerCheckin)
	if err != nil {
		flash.Warning("创建活动失败 " + err.Error())
		flash.Store(&c.Controller)
		return
	}
	flash.Warning("创建活动成功 ")
	flash.Store(&c.Controller)
	c.Redirect(c.URLFor("AttendanceController.Join", "aid", act.Aid), 303)
}

func (c *AttendanceController) Checkin() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.URLFor("UsersController.Index"))
		return
	}

	c.TplName = "attendance/checkin.tpl"

	jalId, _ := c.GetInt("jalid")

	flash := beego.NewFlash()
	if jalId == 0 {
		flash.Warning("没有指定aid")
		flash.Store(&c.Controller)
		return
	}

	jal := attendance.GetJoinActivityLog(jalId)
	if jal == nil || jal.Uid != c.Userinfo.Uid {
		flash.Warning("jal不存在")
		flash.Store(&c.Controller)
		return
	}

	// 查看已经参与的活动
	//jals := attendance.ListUserActivityLog(c.Userinfo.Uid, aid, []interface{}{models.JalStatusAchieved, models.JalStatusInited})

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["jal"] = jal
	//c.Data["jals"] = jals

	cils := attendance.ListUserCheckInLog(c.Userinfo.Uid, jal.JalId, 0)
	c.Data["cils"] = cils
	c.Data["cilsTotal"] = len(cils)

	if jal.Status == models.JalStatusMissed || jal.Status == models.JalStatusShared {
		flash.Error("活动方案[%s]参与后，未完成按计划打卡，请重新参与", jal.Aid.Name)
		flash.Store(&c.Controller)
		return
	}
	if jal.Status == models.JalStatusStopped {
		flash.Error("活动方案[%s]参与后，已中断，请重新参与", jal.Aid.Name)
		flash.Store(&c.Controller)
		return
	}

	if !c.Ctx.Input.IsPost() {
		return
	}

	if !c.CheckXSRFCookie() {
		flash.Warning("页面过期，请刷新后重试")
		flash.Store(&c.Controller)
		return
	}

	err := attendance.UserCheckIn(c.Userinfo.Uid, jal)
	if err != nil {
		logs.Error("user(%d) jal(%d) CheckIn failed:%v", jal.Uid, jal.JalId, err)
		flash.Warning("活动[%s]打卡失败：%v", jal.Aid.Name, err)
		flash.Store(&c.Controller)
		return
	}

	flash.Success("活动[%s](%d)打卡成功", jal.Aid.Name, jal.Aid.Aid)
	flash.Store(&c.Controller)
	c.Redirect(c.Ctx.Request.RequestURI, 303)
}

func (c *AttendanceController) MyCheckInLog() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.URLFor("UsersController.Index"))
		return
	}

	c.TplName = "attendance/mycheckinlog.tpl"

	jalId, _ := c.GetInt("jalid")
	flash := beego.NewFlash()

	if jalId == 0 {
		flash.Warning("没有指定jalid")
		flash.Store(&c.Controller)
		return
	}

	jal := attendance.GetJoinActivityLog(jalId)

	c.Data["jal"] = jal
	c.Data["schedules"] = attendance.Json2CheckInSchedules(jal.Schedule)
	c.Data["act"] = attendance.GetActivity(jal.Aid.Aid)

}
