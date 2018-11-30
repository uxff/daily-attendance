package controllers

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
