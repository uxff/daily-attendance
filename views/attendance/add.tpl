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
                            <input type="text" class="form-control" id="name" name="activity_name" placeholder="请输入名称">
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
                            <label for="needStep">连续多长时间可以分红</label>
                        </div>
                        <div class="col-md-4">
                            <input type="text" class="form-control" style="width: 80px;display: inline" name="needStep" placeholder="请输入天数" value="5">x
                            <select class="form-control da-checkin-period da-input-sm" name="checkInPeriod">
                                <option value="3">分钟</option>
                                <option value="4">小时</option>
                                <option value="5" selected>天</option>
                                <option value="6">周</option>
                                <option value="7">月</option>
                            </select>
                            <p>设置5天表示要求连续5天每天打一次卡才能享受分红</p>
                        </div>
                    </div>
                    <div class="form-group row">
                        <div class="col-md-4">
                            <label for="checkInRule">打卡规则</label>
                        </div>
                        <div class="col-md-4">
                            <input type="text" class="form-control da-checkin-rule" name="checkInRule" placeholder="请输入规则" value='{"HEALTH":{"timespan":"06:00-09:00"}}'>
                        </div>

                    </div>
                    <div class="form-group row">
                        <div class="col-md-4">
                            <label for="checkInRule">时间段</label>
                        </div>
                        <div class="col-md-4">
                            <span class="da-timespan da-timespan-3">任意</span>
                            <span class="da-timespan da-timespan-4">

                            <select class="form-control da-timespan-start-4 da-input-sm" name="minspan-start">
                            {{range $k, $v := genlist 60 0 1}}
                                <option value="{{$v}}:00">{{$v}}分</option>
                            {{end}}
                            </select>-
                            <select class="form-control da-timespan-end-4 da-input-sm" name="minspan-end">
                            {{range $k, $v := genlist 60 0 1}}
                                <option value="{{$v}}:59">{{$v}}分</option>
                            {{end}}
                            </select>
                            </span>
                            <span class="da-timespan da-timespan-5">
                            <select class="form-control da-timespan-start-5 da-input-sm" name="hourspan-start">
                            {{range $kh, $vh := genlist 24 0 1}}
                            {{range $km, $vm := genlist 12 0 5}}
                                <option value="{{$vh}}:{{$vm}}">{{$vh}}点{{$vm}}分</option>
                            {{end}}
                            {{end}}
                            </select>-
                            <select class="form-control da-timespan-end-5 da-input-sm" name="hourspan-end">
                            {{range $kh, $vh := genlist 24 0 1}}
                            {{range $km, $vm := genlist 12 0 5}}
                                <option value="{{$vh}}:{{$vm}}">{{$vh}}点{{$vm}}分</option>
                            {{end}}
                            {{end}}
                                <option value="23:59">23点59分</option>
                            </select>
                            </span>
                            <span class="da-timespan da-timespan-6">
                            <select class="form-control da-timespan-start-6 da-input-sm" name="weekdayspan-start">
                            {{range $k, $v := genlist 7 1 1}}
                                <option value="{{$k}}">{{$v}}</option>
                            {{end}}
                            </select>-
                            <select class="form-control da-timespan-end-6 da-input-sm" name="weekdayspan-end">
                            {{range $k, $v := genlist 7 1 1}}
                                <option value="{{$k}}">{{$v}}</option>
                            {{end}}
                            </select>
                            </span>
                            <span class="da-timespan da-timespan-7">
                            <select class="form-control da-timespan-start-7 da-input-sm" name="dayspan-start">
                            {{range $k, $v := genlist 31 1 1}}
                                <option value="{{$k}}">{{$v}}日</option>
                            {{end}}
                            </select>-
                            <select class="form-control da-timespan-end-7 da-input-sm" name="dayspan-end">
                            {{range $k, $v := genlist 31 1 1}}
                                <option value="{{$k}}">{{$v}}日</option>
                            {{end}}
                            </select>
                            </span>
                        </div>

                    </div>
                    <div class="form-group row hide">
                        <div class="col-md-12 col-md-offset-0">
                        <p >健康打卡规则举例:{"HEALTH":{"timespan":"06:00-09:00"}}</p>
                        <p >上班打卡规则举例:{"WORKUP":{"timespan":"00:00-10:00"},"WORKOFF":{"timespan":"18:00-23:59"}}</p>
                        <p >月报打卡规则举例:{"REPROT":{"dayspan":"01-05"}}</p>
                        </div>
                    </div>
                    <div class="form-group row">
                        <div class="col-md-4">
                            <label for="joinPrice">参与需要的押金</label>
                        </div>
                        <div class="col-md-4">
                            <input type="text" class="form-control" style="width: 80px;display: inline" name="joinPrice" placeholder="请输入积分" value="50"> 积分
                        </div>
                    </div>
                    <div class="form-group row">
                        <div class="col-md-4">
                            <label for="awardPerCheckin">每次打卡奖励</label>
                        </div>
                        <div class="col-md-4">
                            <input type="text" class="form-control" style="width: 80px;display: inline" name="awardPerCheckin" placeholder="请输入积分" value="10"> 积分
                        </div>
                    </div>
                    {{.xsrfdata}}
                    <div class="form-group text-center">
                        <button type="submit" class="btn btn-primary da-btn-add-act">提交</button>
                    </div>
                </form>

            </div>
        </div>
        </div>
    </div>
</div>
