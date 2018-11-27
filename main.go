package main

import (
	"flag"
	"github.com/astaxie/beego"
	_ "github.com/uxff/daily-attendance/conf/inits"
	_ "github.com/uxff/daily-attendance/routers"
	"github.com/astaxie/beego/logs"
	"github.com/uxff/daily-attendance/models"
)

func main() {
	logdeep := 3
	addr := ":"+ beego.AppConfig.String("httpport")

	flag.StringVar(&addr, "addr", addr, "beego run param addr, format as ip:port")
	flag.Parse()

	logs.SetLevel(logs.LevelDebug)
	logs.SetLogFuncCallDepth(logdeep)

	//models.LoadIndexLinksFromFile("./conf/friends.json")

	models.SetFriendlyLinksPath("./conf/friends.json")
	models.LoadFriendlyLinks()

	logs.Info("beego server will run. addr=%s", addr)

	beego.Run(addr)
}
