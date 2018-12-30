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
        <h4>打卡计划：{{.act.Name}}</h4>
        <div class="col-lg-2"></div>
        <div class="col-lg-8">
            <table class="table table-bordered table-striped">
                <thead>
                <tr>
                    <td>步骤</td>
                    <td>时间</td>
                    <td>签到</td>
                </tr>
                </thead>
                <tbody>
                {{range $step, $schedule := .schedules}}
                {{range $_, $checkInElem := $schedule}}
                <tr>
                    <td>{{$step}}</td>
                    <td>{{$checkInElem.From}} -> {{$checkInElem.To}}</td>
                    <td>
                    {{ if gt $checkInElem.CilId 0}}
                        <span class="label label-success "><i class="glyphicon glyphicon-check"></i></span>
                    {{else}}
                        <span class="label label-default "><i class="glyphicon glyphicon-check"></i></span>
                    {{end}}
                    </td>
                </tr>
                {{end}}
                {{end}}
                <tr>
                    <td colspan="12">{{.jal.Step}}/{{.jal.BonusNeedStep}} {{checkinperiod .jal.Aid.CheckInPeriod}}</td>
                </tr>
                </tbody>
            </table>
        </div>
        <div class="col-lg-2"></div>
    </div>
    <div class="row">
        <div class="col-lg-2 col-lg-offset-2">
        </div>
    </div>
</div>
