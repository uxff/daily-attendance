package main

import (
	"flag"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/uxff/daily-attendance/conf/inits"
	_ "github.com/uxff/daily-attendance/routers"

	"github.com/uxff/daily-attendance/lib/modules/attendance"
	"github.com/uxff/daily-attendance/models"
)

func main() {
	logdeep := 3
	addr := ":" + beego.AppConfig.String("httpport")

	flag.StringVar(&addr, "addr", addr, "beego run param addr, format as ip:port")
	flag.Parse()

	logs.SetLevel(logs.LevelDebug)
	logs.SetLogFuncCallDepth(logdeep)

	//models.LoadIndexLinksFromFile("./conf/friends.json")

	models.SetFriendlyLinksPath("./conf/friends.json")
	models.LoadFriendlyLinks()

	aperiod, _ := beego.AppConfig.Int("accounting_period")
	attendance.SetAccountingPeriod(aperiod)
	go attendance.AutoAccounting()

	logs.Info("beego server will run. addr=%s", addr)

	beego.Run(addr)
}
