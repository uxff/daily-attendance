package routers

import (
	"github.com/astaxie/beego"
	"github.com/uxff/daily-attendance/controllers"
)

func init() {
	beego.Router("/", &controllers.IndexController{}, "get:Index")
	beego.Router("/user", &controllers.UsersController{}, "get,post:Index")
	beego.Router("/login", &controllers.UsersController{}, "get,post:Login")
	beego.Router("/logout", &controllers.UsersController{}, "get:Logout")
	beego.Router("/signup", &controllers.UsersController{}, "get,post:Signup")
	beego.Router("/attendance", &controllers.AttendanceController{}, "get:Index")
	beego.Router("/attendance/add", &controllers.AttendanceController{}, "get,post:Add")
	beego.Router("/attendance/project", &controllers.AttendanceController{}, "get,post:Project")
}
