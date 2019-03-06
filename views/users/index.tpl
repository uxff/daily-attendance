{{append . "HeadStyles" "/static/css/custom.css"}}
{{append . "HeadScripts" "/static/js/custom.js"}}

{{template "alert.tpl" .}}
<div class="container text-center">
    <div class="row">
        <ul class="nav nav-tabs">
            <li class="active"><a href="javascript:;">基本资料</a></li>
            <li><a href="/user/balance">积分</a></li>
            <li><a href="/attendance/my">打卡</a></li>
            <li><a href="/user/bonus">收益</a></li>
            <li><a href="/user/invite">推广码</a></li>
        </ul>
        <h4>我的个人资料</h4>
        <div class="col-lg-3"></div>
        <div class="col-lg-6">
            <table class="table table-bordered table-striped">
                <tr>
                    <td>账号：</td>
                    <td>  {{.Userinfo.Email}} </td>
                </tr>
                <tr>
                    <td>昵称：</td>
                    <td> {{.Userinfo.Nickname}} </td>
                    
                </tr>
                <tr>
                    <td>手机号：</td>
                    <td>  {{.Userinfo.Phone}}</td>
                </tr>
                <tr>
                    <td>注册时间：</td>
                    <td>  {{.Userinfo.Created}}</td>
                </tr>
                <tr>
                    <td>激活认证：</td>
                    <td>
                    {{if .Userinfo.Openid}}<span class="label label-default">微信</span>{{end}}
                    {{if .Userinfo.Email}}{{if not .Userinfo.Openid}}<span class="label label-default">Email</span>{{end}}{{end}}
                    {{if .Userinfo.Phone}}<span class="label label-default">手机号</span>{{end}}
                    </td>
                </tr>
                <tr>
                    <td>最后登录时间：</td>
                    <td>  {{.Userinfo.Lastlogintime}}</td>
                </tr>
                <tr>
                    <td>最后登录IP：</td>
                    <td>  {{.Userinfo.Lastloginip}}</td>
                </tr>
                <tr>
                    <td>头像：</td>
                    <td>
                        {{if .Userinfo.WxLogoUrl}}
                            <img src="{{.Userinfo.WxLogoUrl}}">
                        {{end}}
                    </td>
                </tr>
            </table>
        </div>
        <div class="col-lg-3"></div>
    </div>
</div>
