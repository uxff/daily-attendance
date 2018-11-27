package controllers

import (
	"github.com/uxff/daily-attendance/models"
)

type IndexController struct {
	BaseController
}

func (this *IndexController) Index() {

	theLinks := models.LoadIndexLinks()

	this.Data["thelinks"] = theLinks

	this.TplName = "index/index.tpl"
}

