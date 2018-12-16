package routers

import (
	"github.com/astaxie/beego"
	"github.com/uxff/daily-attendance/controllers"
)

func init() {
	beego.Router("/", &controllers.IndexController{}, "get:Index")
	beego.Router("/user", &controllers.UsersController{}, "get,post:Index")
	beego.Router("/user/index", &controllers.UsersController{}, "get:Index")
	beego.Router("/login", &controllers.UsersController{}, "get,post:Login")
	beego.Router("/logout", &controllers.UsersController{}, "get:Logout")
	beego.Router("/signup", &controllers.UsersController{}, "get,post:Signup")
	beego.Router("/attendance", &controllers.AttendanceController{}, "get:Index")
	beego.Router("/attendance/add", &controllers.AttendanceController{}, "get,post:Add")
	beego.Router("/attendance/join", &controllers.AttendanceController{}, "get,post:Join")
	beego.Router("/attendance/my", &controllers.AttendanceController{}, "get,post:My")
	beego.Router("/attendance/checkin", &controllers.AttendanceController{}, "get,post:Checkin")
	beego.Router("/attendance/mycheckinlog", &controllers.AttendanceController{}, "get:MyCheckInLog")

	beego.Router("/user/balance", &controllers.UsersController{}, "get:Balance")
	beego.Router("/user/invite", &controllers.UsersController{}, "get:Invite")
	beego.Router("/user/bonus", &controllers.UsersController{}, "get:Bonus")
}
