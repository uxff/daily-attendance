{{append . "HeadStyles" "/static/css/custom.css"}}
{{append . "HeadScripts" "/static/js/custom.js"}}

{{template "alert.tpl" .}}
<div class="container text-center">
    <div class="row">
        <ul class="nav nav-tabs">
            <li><a href="/user/index">基本资料</a></li>
            <li><a href="/user/balance">积分</a></li>
            <li class="active"><a href="javascript:;">打卡</a></li>
            <li><a href="/user/bonus">收益</a></li>
            <li><a href="/user/invite">推广码</a></li>
        </ul>
        <h4>我的参与的打卡计划</h4>
        <div class="btn-group" role="group" aria-label="...">
            <button type="button" class="btn btn-default da-btn-all">全部</button>
            <button type="button" class="btn btn-default da-btn-ok">坚持中</button>
            <button type="button" class="btn btn-default da-btn-fail">错过</button>
        </div>
        <div class="col-lg-2"></div>
        <div class="col-lg-8">
            <table class="table table-bordered table-striped">
                <thead>
                <tr>
                    <td>参与时间</td>
                    <td>活动名称</td>
                    <td>活动描述</td>
                    <td>花费积分</td>
                    <td>完成度</td>
                    <td>状态</td>
                    <td>累计获得奖励</td>
                    <td>操作</td>
                </tr>
                </thead>
                <tbody>
            {{if eq 0 .total }}
                <tr>
                    <td colspan="12">您还没有参与任何活动</td>
                </tr>
            {{else}}
                {{range $ji, $jal := .jals}}
                <tr class="da-tr {{inarr .Status "1,2" "da-tr-ok" ""}} {{inarr .Status "3,4,5" "da-tr-fail" ""}}">
                    <td>{{timefmtm .Created}}</td>
                    <td>{{.Aid.Name}}</td>
                    <td>{{.Aid.Desc}}</td>
                    <td>{{.JoinPrice}}</td>
                    <td>{{.Step}}/{{.BonusNeedStep}}{{checkinperiod .Aid.CheckInPeriod}}</td>
                    <td>{{jalstatus .Status}}</td>
                    <td>0</td>
                    <td>
                        <a href="/attendance/checkin?jalid={{.JalId}}">打卡</a>
                        <a href="/attendance/mycheckinlog?jalid={{.JalId}}">详情</a>
                    </td>
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
