{{append . "HeadStyles" "/static/css/custom.css"}}
{{append . "HeadScripts" "/static/js/custom.js"}}

{{template "alert.tpl" .}}
<div class="container text-center">
    <div class="row">
        <ul class="nav nav-tabs">
            <li><a href="/user/index">基本资料</a></li>
            <li><a href="/user/balance">积分</a></li>
            <li><a href="/attendance/my">打卡</a></li>
            <li class="active"><a href="javascript:;">收益</a></li>
            <li><a href="/user/invite">推广码</a></li>
        </ul>
        <h4>我的收益</h4>
        <div class="col-lg-2">
            <p>累计收益</p>
            <label class="label label-success">{{.balance.Balance}} 积分</label>
        </div>
        <div class="col-lg-8">
            <table class="table table-bordered table-striped">
                <thead>
                <tr>
                    <td>交易号</td>
                    <td>金额(积分)</td>
                    <td>活动名</td>
                    <td>时间</td>
                    <td>操作</td>
                </tr>
                </thead>
                <tbody>
                {{if eq 0 .total }}
                <tr>
                    <td colspan="12">没有记录</td>
                </tr>
                {{else}}
                {{range $k, $ws := .bonusLog}}
                <tr>
                    <td>{{.WsId}}</td>
                    <td>{{.Amount}}</td>
                    <td>{{.Aid.Name}}</td>
                    <td>{{timefmtm .Created}}</td>
                    <td>
                        <a href="/attendance/join?aid={{.Aid.Aid}}" >参与活动</a>
                        <a href="/attendance/checkin?jalid={{.ToJalId}}" >继续打卡</a>
                        <a href="/trade/detail?utlid={{.UtlId}}" >交易详情</a>
                    </td>
                </tr>
                {{end}}
                {{end}}
                </tbody>
            </table>
        </div>
        <div class="col-lg-2"></div>
        <div class="col-lg-8 col-lg-offset-2 text-left">共{{.total}}条</div>
    </div>
</div>
