{{append . "HeadStyles" "/static/css/custom.css"}}
{{append . "HeadScripts" "/static/js/custom.js"}}

{{template "alert.tpl" .}}
<div class="container">
    <div class="row">
        <div class="col-lg-8 col-lg-offset-2">
            <div class="panel panel-info">
                <div class="panel-heading">参与活动{{.act.Name}}</div>
                <div class="panel-body">
                    <form role="form" action="" method="post">
                        <div class="form-group row">

                        </div>
                        <div class="form-group row">
                            <div class="col-md-4 col-md-offset-2">
                                <label>活动名称</label>
                            </div>
                            <div class="col-md-4">
                                <label>{{.act.Name}}</label>
                            </div>
                        </div>
                        <div class="form-group row">
                            <div class="col-md-4 col-md-offset-2">
                                活动参与价格
                            </div>
                            <div class="col-md-4">
                                <label class="label label-danger">{{.act.JoinPrice}} 积分</label>
                            </div>
                        </div>
                        <div class="form-group row">
                            <div class="col-md-4 col-md-offset-2">
                                连续多长时间可以分红
                            </div>
                            <div class="col-md-4">
                                <label class="label label-warning">{{.act.BonusNeedStep}} 天</label>
                            </div>
                        </div>
                        <div class="form-group row text-center">
                            <div class="col-md-6 col-md-offset-2">
                                {{.xsrfdata}}
                                {{range $k, $jal := .jals}}
                                    <p>您已经于{{timefmtm $jal.Created}}参与了该活动。<a href="/attendance/checkin?jalid={{.JalId}}">打卡</a></p>
                                {{end}}
                                您确定要参加该活动吗？
                                <button type="submit" class="btn btn-primary">确认参与</button>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>
</div>