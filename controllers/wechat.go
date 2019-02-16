package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/message"

	"github.com/uxff/daily-attendance/lib/modules/wxoa/wxmodels"
	"github.com/uxff/daily-attendance/models"
)

type WechatController struct {
	beego.Controller

	Openid   string
	Userinfo *models.User
	Wxoa     *wxmodels.WechatOfficalAccounts
}

func (c *WechatController) Prepare() {

	//mpid := c.Input().Get("oa")
	oaId, _ := c.GetInt("oa")

	woa := wxmodels.GetWoa(oaId)

	if woa == nil {
		return
	}

	//config :=
	/*
		需要基于该类重新设计：
		运行中config可以设置多次，但是setMessageHandler不需要每次请求都调用。
	*/

	wc := wechat.NewWechat(&wechat.Config{
		AppID:          woa.Appid,
		AppSecret:      woa.Appsecret,
		Token:          woa.Token,
		EncodingAESKey: woa.EncodingAesKey,
	})

	// 传入request和responseWriter
	server := wc.GetServer(c.Ctx.Request, c.Ctx.ResponseWriter)

	server.GetOpenID()

	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg message.MixMessage) *message.Reply {

		//回复消息：演示回复用户发送的消息
		text := message.NewText(msg.Content)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}

	//发送回复的消息
	server.Send()
}

func (c *WechatController) CheckOauth() {

}

func (c *WechatController) Index() {

}
