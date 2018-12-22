{{append . "HeadStyles" "/static/css/custom.css"}}
{{append . "HeadScripts" "/static/js/custom.js"}}


<div class="container">
    <div class="row">
    {{template "alert.tpl" .}}
        <h2><b>打卡方案介绍</b></h2>
        <p><span class="label label-success">为了健康，为了自己，坚持打卡，告别慵懒</span>
            <span class="label label-info"><a href="/attendance/add">创建活动</a></span>
        </p>

        {{range $i, $act := .activities}}
        <div class="panel {{mod $i 2 "panel-success" "panel-info"}} ">
            <div class="panel-heading">打卡方案{{$act.Aid}}-{{$act.Name}} </div>
            <div class="panel-body container">
                <div class="row">
                    <div class="col-lg-6 col-sm-10">
                        <p>花{{$act.JoinPrice}}积分参与本活动，作为押金积分</p>
                        <p>参与活动后，坚持每天早晨6-9点打卡</p>
                        <p>如果有1天中断，押金积分的100%将进入“公共奖金池”</p>
                        <p>连续坚持{{$act.BonusNeedStep}}天以上，将可以参与瓜分“公共奖金池”</p>
                        <p>中断后，可以补交后，次日继续打卡参与活动，重新计算累计，历史累计清零</p>
                        <p>中断后，如果没补交，每多中断一天，剩余押金将继续按100%比例扣除进入“公共奖金池”</p>
                        <p>押金扣到0为止</p>
                        <p>活动最终解释权归我公司所有</p>
                    </div>
                    <div class="col-lg-6 col-sm-10">
                        <div class="well well-sm">
                            <p>参与数量：<span class="label label-warning">{{$act.JoinedUserCount}}人 / {{$act.JoinedAmount}} 积分</span></p>
                            <p>参与价格：<span class="label label-info">{{$act.JoinPrice}} 积分</span></p>
                            <p>奖池金额：<span class="label label-success">{{$act.UnsharedAmount}} 积分</span></p>
                            <p>已瓜分：<span class="label label-danger">{{$act.SharedAmount}} 积分</span></p>
                            <p>要求：<span class="label label-success">连续 {{$act.BonusNeedStep}} {{checkinperiod $act.CheckInPeriod}}</span></p>
                            <p>
                                &nbsp;<a class="btn btn-primary" type="button" href="/attendance/join?aid={{$act.Aid}}">&nbsp;立即参与&nbsp;</a>
                                &nbsp;<a class="btn btn-primary" type="button" href="javascript:;">&nbsp;查看排行&nbsp;</a>
                            </p>
                            <p></p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        {{end}}


    </div>
</div>
