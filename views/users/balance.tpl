{{append . "HeadStyles" "/static/css/custom.css"}}
{{append . "HeadScripts" "/static/js/custom.js"}}

{{template "alert.tpl" .}}
<div class="container text-center">
    <div class="row">
        <ul class="nav nav-tabs">
            <li><a href="/user/index">基本资料</a></li>
            <li class="active"><a href="javascript:;">积分</a></li>
            <li><a href="/attendance/my">打卡</a></li>
            <li><a href="/user/bonus">收益</a></li>
            <li><a href="/user/invite">推广码</a></li>
        </ul>
        <h4>我的账户余额</h4>
        <div class="col-lg-2">
            <p>当前余额</p>
            <label class="label label-success">{{.balance.Balance}}积分</label>
        </div>
        <div class="col-lg-8">
            <table class="table table-bordered table-striped">
                <thead>
                <tr>
                    <td>交易号</td>
                    <td>金额(积分)</td>
                    <td>交易类型</td>
                    <td>来源</td>
                    <td>交易后余额</td>
                    <td>时间</td>
                    <td>状态</td>
                    <td>备注</td>
                </tr>
                </thead>
                <tbody>
                {{if eq 0 .total }}
                <tr>
                    <td colspan="12">没有记录</td>
                </tr>
                {{else}}
                {{range $utli, $utl := .utls}}
                <tr>
                    <td>{{.UtlId}}</td>
                    <td>{{.Amount}}</td>
                    <td>{{.TradeType}}</td>
                    <td>{{.SourceType}}</td>
                    <td>{{.Balance}}</td>
                    <td>{{timefmtm .Created}}</td>
                    <td>{{jalstatus .Status}}</td>
                    <td>{{.Remark}}</td>
                </tr>
                {{end}}
                {{end}}
                </tbody>
            </table>
        </div>
        <div class="col-lg-2"></div>
    </div>
    <div class="row">
        <div class="col-lg-2 col-lg-offset-2">
            共{{.total}}条
        </div>
    </div>
</div>
