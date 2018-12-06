package controllers

import (
	"html/template"
	"time"
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/uxff/daily-attendance/lib/modules/attendance"
)

type AttendanceController struct {
	BaseController
}

func (c *AttendanceController) NestPrepare() {
}

// func (c *UsersController) NestFinish() {}

func (c *AttendanceController) Index() {

	activities := attendance.ListActivities()
	c.Data["activities"] = activities

	c.TplName = "attendance/index.tpl"
}

func (c *AttendanceController) Project() {

}

func (c *AttendanceController) Join() {
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

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["act"] = act
	if !c.Ctx.Input.IsPost() {

		return
	}

	if !c.CheckXSRFCookie() {
	}

	utlId := 0
	var err error
	if act.JoinPrice>0 {
		utlId, err = attendance.Consume(c.Userinfo.Uid, attendance.ActivityToProduct(act), 1)
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

	flash.Success("参与活动%d成功", aid)
	flash.Store(&c.Controller)
}

func (c *AttendanceController) Add() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.URLFor("AttendanceController.Index"))
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

	name := c.GetString("name")
	startTimeStr := c.GetString("startTime")
	endTimeStr := c.GetString("endTime")
	needStep, _ := c.GetInt("needStep")
	checkInRuleStr := c.GetString("checkInRule")
	joinPrice, _ := c.GetInt("joinPrice")
	wastagePercent, _ := c.GetFloat("wastagePercent")
	checkInPeriod, _ := c.GetInt8("checkInPeriod")

	startTime, err := time.Parse("2006-01-02 15:04", startTimeStr)
	endTime, err := time.Parse("2006-01-02 15:04", endTimeStr)


	checkInRuleMap := new(attendance.CheckInRuleMap)
	err = json.Unmarshal([]byte(checkInRuleStr), checkInRuleMap)
	if err != nil {
		flash.Warning("规则格式不正确 "+ err.Error()+ " origin:"+checkInRuleStr)
		flash.Store(&c.Controller)
		return
	}
	logs.Warn("will create activity:%v %v %v %v %v %v", name, startTime, endTime, needStep, checkInRuleMap, wastagePercent)

	err = attendance.AddActivity(name, startTime, endTime, *checkInRuleMap, needStep, checkInPeriod, c.Userinfo.Uid, joinPrice, float32(wastagePercent))
	if err != nil {
		flash.Warning("创建活动失败 "+ err.Error())
		flash.Store(&c.Controller)
		return
	}
	flash.Warning("创建活动成功 ")
	flash.Store(&c.Controller)
}

