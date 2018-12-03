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

	c.TplName = "attendance/index.tpl"
}

func (c *AttendanceController) Project() {

}

func (c *AttendanceController) Join() {

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

