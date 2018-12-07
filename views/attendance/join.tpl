{{append . "HeadStyles" "/static/css/custom.css"}}
{{append . "HeadScripts" "/static/js/custom.js"}}

{{template "alert.tpl" .}}
<div class="container">
    <div class="row">
        <div class="col-md-4 col-md-offset-2">
            活动名称
        </div>
        <div class="col-md-4 col-md-offset-6">
            {{.act.Name}}
        </div>
    </div>
    <div class="row">
        <div class="col-md-4 col-md-offset-2">
            活动参与价格
        </div>
        <div class="col-md-4 col-md-offset-6">
            {{.act.JoinPrice}} 积分
        </div>
    </div>
    <div class="row">
        <div class="col-md-4 col-md-offset-2">
            连续多长时间可以分红
        </div>
        <div class="col-md-4 col-md-offset-6">
            {{.act.BonusNeedStep}} 天
        </div>
    </div>
    <div class="row text-center">
        <div class="col-md-4 col-md-offset-2">
            <form action="" method="post">
                {{.xsrfdata}}
                {{range $k, $jal := .jals}}
                    <p>您已经于{{$jal.Created}}参与了该活动({{$jal.Aid}})。</p>
                {{end}}
                您确定要参加该活动吗？
                <button type="submit" class="btn btn-primary">确认参与</button>
            </form>
        </div>
    </div>
</div>