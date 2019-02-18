package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/message"

	"github.com/uxff/daily-attendance/lib/modules/wxoa/wxmodels"
	"github.com/uxff/daily-attendance/models"
	"github.com/astaxie/beego/logs"
	"strings"
	"github.com/uxff/daily-attendance/lib/utils/oncetoken"
	"time"
	"github.com/uxff/daily-attendance/lib"
)

const (
	OauthToken  = "fromda"
)

/*
	微信mp接口配置url应该填写：
		yourdomain.com/wechat/index?oa=1
*/
type WechatController struct {
	beego.Controller

	Openid   string
	Userinfo *models.User
	Wxoa     *wxmodels.WechatOfficalAccounts
	Wc 	     *wechat.Wechat
	OauthToken string
	WechatApiDomain string // 需要nginx配置 proxy_set_header X-Host $host;proxy_set_header X-Scheme $scheme;
}

func (c *WechatController) Prepare() {

	// 域名通过nginx配置获取 proxy_set_header X-Host $host;
	c.WechatApiDomain =c.Ctx.Request.Header.Get("X-Scheme")+"://"+ c.Ctx.Request.Header.Get("X-Host")
	//mpid := c.Input().Get("oa")
	oaId, _ := c.GetInt("oa")

	woa := wxmodels.GetWoa(oaId)

	if woa == nil {
		logs.Warn("wxoaId not found:%d", oaId)
		return
	}

	//config :=
	/*
		需要基于该类重新设计：
		运行中config可以设置多次，但是setMessageHandler不需要每次请求都调用。
	*/

	c.Wc = wechat.NewWechat(&wechat.Config{
		AppID:          woa.Appid,
		AppSecret:      woa.Appsecret,
		Token:          woa.Token,
		EncodingAESKey: woa.EncodingAesKey,
	})

	// 传入request和responseWriter
	server := c.Wc.GetServer(c.Ctx.Request, c.Ctx.ResponseWriter)

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

// open this on pc, show a qr code
func (c *WechatController) ShowQrForLogin() {

}

// on wechat client, open this, it will start oauth
func (c *WechatController) OauthLogin() {
	oaId := c.GetString("oa")
	oauth := c.Wc.GetOauth()
	logs.Info("a oauthlogin request, url:%s", c.WechatApiDomain+c.URLFor("WechatController.Index", "oa", oaId))

	//otoken := oncetoken.GenToken()

	oauthUrl, err := oauth.GetRedirectURL(c.WechatApiDomain+ c.URLFor("WechatController.OauthCallback", "oa", oaId), "snsapi_userinfo", OauthToken)
	logs.Info("oauth url:%s", oauthUrl)

	err = oauth.Redirect(c.Ctx.ResponseWriter, c.Ctx.Request, c.WechatApiDomain+ c.URLFor("WechatController.OauthCallback", "oa", oaId), "snsapi_userinfo", OauthToken)
	if err != nil {
		logs.Error("make oauth login failed:%v", err)
		c.Ctx.WriteString(fmt.Sprintf("make oauth login failed:%v", err))
		return
	}

	//
	logs.Info("start oauth ok: oaId:%s", oaId)
	//c.Redirect(c.URLFor("WechatController.Index"))
}


//
func (c *WechatController) OauthCallback() {
	//oaId := c.GetString("oa")
	oauth := c.Wc.GetOauth()
	code := c.GetString("code")
	resToken, err := oauth.GetUserAccessToken(code)
	if err != nil {
		logs.Error("get access token by code failed: code=%s", code)
		return
	}

	logs.Info("resToken is :%+v", resToken)
	userInfo, err := oauth.GetUserInfo(resToken.AccessToken, resToken.OpenID)
	if err != nil {
		logs.Error("get userinfo by openid failed: openid=%s", resToken.OpenID)
		return
	}

	logs.Info("userInfo from wx:%+v", userInfo)

	var uid int

	// 查看 openid 是否注册过 未注册则注册并登录 注册则登录
	existUser := models.GetByEmail(userInfo.OpenID)
	if existUser != nil && existUser.Uid> 0 {
		uid = existUser.Uid
		logs.Info("user exist when oauthcallback: uid:%d", uid)
	} else {
		u := &models.User{
			Email:resToken.OpenID,
			Openid:resToken.OpenID,
			WxLogoUrl:userInfo.HeadImgURL,
			//WxNickname:userInfo.Nickname,
			Nickname:userInfo.Nickname,
		}

		if c.Wxoa != nil {
			logs.Info("wxoa exist when register by wechat")
			u.WoaId = c.Wxoa.WoaId
		}

		u.Lastlogintime = time.Unix(0, 0)

		// 必须填写默认值，否则数据库报错
		u.EmailActivated = time.Time{}
		u.PhoneActivated = time.Unix(0, 0)
		u.WxUnsubscribed = time.Unix(0, 0)

		uid, err = lib.SignupUser(u)
		if err != nil || uid < 1 {
			logs.Error("register from wx failed:%v", err)
			return
		}

		logs.Info("register from wx ok:id:%d openid:%s", uid, u.Openid)
	}



	utoken := oncetoken.GenToken()

	c.Redirect(c.URLFor("UsersController.LoginByWechat", "token", utoken, "uid", uid), 303)

}

// show qr code
func (c *WechatController) Index() {
	oaId := c.GetString("oa")
	// 判断是否微信中打开
	ua := c.Ctx.Request.Header.Get("User-Agent")
	if strings.Index(strings.ToLower(ua), "micromessenger") > 0 {
		// is in wechat client
		logs.Info("this is in micromessenger, will redirect to oauth login")
		c.Redirect(c.URLFor("WechatController.OauthLogin", "oa", oaId), 303)
	}

	logs.Info("this is NOT in micromessenger. host:%s", c.Ctx.Request.Header.Get("X-Host"))
	c.Ctx.WriteString("please open this page in wechat client")
}
