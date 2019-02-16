package inits

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/lib/pq"
	//_ "github.com/mattn/go-sqlite3"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/uxff/daily-attendance/lib/modules/attendance/models"
	_ "github.com/uxff/daily-attendance/models"
	"github.com/uxff/daily-attendance/lib/modules/wxoa/wxmodels"
	"github.com/uxff/daily-attendance/models"
)

func init() {

	runmode := beego.AppConfig.String("runmode")
	dbname := "default" //beego.AppConfig.String("dbname")
	datasource := beego.AppConfig.String("datasource")

	switch runmode {
	//case "prod":
	case "dev":
		orm.Debug = true
		fallthrough
	default:
		orm.RegisterDataBase(dbname, "mysql", datasource, 30)
	}

	orm.DefaultTimeLoc = time.FixedZone("Asia/Shanghai", 8*60*60)


	// 这一步必须要在 orm.RunSyncdb(dbname, force, verbose) 前
	orm.RegisterModelWithPrefix(
		beego.AppConfig.String("dbprefix"),
		new(models.User),
		new(wxmodels.WechatOfficalAccounts),
	)

	force, verbose := false, true
	// 必须要在 orm.RegisterModelWithPrefix() 后执行
	err := orm.RunSyncdb(dbname, force, verbose)
	if err != nil {
		panic(err)
	}

	// orm.RunCommand()
}
