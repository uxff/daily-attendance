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
        <div class="col-lg-3"></div>
        <div class="col-lg-6">
            <table class="table table-bordered table-striped">
                <thead>
                <tr>
                    <td>参与时间</td>
                    <td>活动名称</td>
                    <td>花费积分</td>
                    <td>完成度</td>
                    <td>状态</td>
                    <td>累计获得奖励</td>
                </tr>
                </thead>
                <tbody>

            {{if eq 0 .total }}
                <tr>
                    <td colspan="12">您还没有参与任何活动</td>
                </tr>
            {{else}}
                {{range $ji, $jal := .jals}}
                <tr>
                    <td>{{timefmtm .Created}}</td>
                    <td>{{.Aid.Name}}</td>
                    <td>{{.JoinPrice}}</td>
                    <td>{{.Step}}/{{.BonusNeedStep}}</td>
                    <td>{{jalstatus .Status}}</td>
                    <td>0</td>
                </tr>
                {{end}}
            {{end}}
                </tbody>
            </table>
        </div>
        <div class="col-lg-3"></div>
    </div>
    <div class="row">
        <div class="col-lg-3 col-lg-offset-2">
            共{{.total}}条
        </div>

    </div>
</div>
