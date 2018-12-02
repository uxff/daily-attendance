{{append . "HeadStyles" "/static/bower/smalot-bootstrap-datetimepicker/css/bootstrap-datetimepicker.css"}}
{{append . "HeadScripts" "/static/bower/smalot-bootstrap-datetimepicker/js/bootstrap-datetimepicker.js"}}
{{append . "HeadStyles" "/static/css/custom.css"}}
{{append . "HeadScripts" "/static/js/custom.js"}}


<div class="container">
    <div class="row">
    {{template "alert.tpl" .}}
        <div class="col-lg-8 col-lg-offset-2">

        <div class="panel panel-info">
            <div class="panel-heading">创建打卡方案</div>
            <div class="panel-body">
                <form role="form" action="" method="post">
                    <div class="form-group row">
                        <div class="col-md-4">
                            <label for="name">名称</label>
                        </div>
                        <div class="col-md-4">
                            <input type="text" class="form-control" id="name" name="name" placeholder="请输入名称">
                        </div>
                    </div>
                    <div class="form-group row">
                        <div class="col-md-4">
                            <label for="startTime">开始时间</label>
                        </div>
                        <div class="col-md-4">
                            <input type="text" class="form-control form_datetime" id="startTime" name="startTime" placeholder="请输入开始时间" data-date-format="yyyy-mm-dd hh:ii" >
                        </div>
                    </div>
                    <div class="form-group row">
                        <div class="col-md-4">
                            <label for="endTime">结束时间</label>
                        </div>
                        <div class="col-md-4">
                            <input type="text" class="form-control form_datetime" id="endTime" name="endTime" placeholder="请输入结束时间" data-date-format="yyyy-mm-dd hh:ii">
                        </div>
                    </div>
                    <div class="form-group row">
                        <div class="col-md-4">
                            <label for="checkInRule">打卡规则</label>
                        </div>
                        <div class="col-md-4">
                            <input type="text" class="form-control" name="checkInRule" placeholder="请输入规则" value='{"HEALTH":{"timespan":"06:00-09:00"}}'>
                        </div>
                        <p class="help-block">打卡规则举例:{"HEALTH":{"timespan":"06:00-09:00"}}</p>
                    </div>
                    <div class="form-group row">
                        <div class="col-md-4">
                            <label for="needStep">连续多长时间可以分红</label>
                        </div>
                        <div class="col-md-4">
                            <input type="text" class="form-control" style="width: 80px;display: inline" name="needStep" placeholder="请输入天数" value="5">
                            <select class="form-control" name="checkInPeriod" style="width: 80px;display: inline">
                                <option value="1">小时</option>
                                <option value="2" selected>天</option>
                                <option value="3">月</option>
                            </select>
                        </div>
                    </div>
                    <div class="form-group row">
                        <div class="col-md-4">
                            <label for="joinPrice">参与需要的押金</label>
                        </div>
                        <div class="col-md-4">
                            <input type="text" class="form-control" style="width: 80px;display: inline" id="joinPrice" placeholder="请输入积分" value="50">积分
                        </div>
                    </div>
                    <div class="form-group row">
                        <div class="col-md-4">
                            <label for="wastagePercent">放弃坚持后，押金扣除比例</label>
                        </div>
                        <div class="col-md-4">
                            <input type="text" class="form-control" style="width: 80px;display: inline" id="wastagePercent" placeholder="请输入比例" value="50">%
                        </div>
                    </div>
                    {{.xsrfdata}}
                    <div class="form-group text-center">
                        <button type="submit" class="btn btn-primary">提交</button>
                    </div>
                </form>

            </div>
        </div>
        </div>
    </div>
</div>
