{{append . "HeadStyles" "/static/css/custom.css"}}
{{append . "HeadScripts" "/static/js/custom.js"}}


<div class="container">
    <div class="row">
    {{template "alert.tpl" .}}
        <h2><b>打卡方案介绍</b></h2>
        <p><span class="label label-success">为了健康，为了自己，坚持打卡，告别慵懒</span>
            <span class="label label-info"><a href="/attendance/add">创建活动</a></span>
        </p>

        <div class="panel panel-success">
            <div class="panel-heading">打卡方案1-早操锻炼 </div>
            <div class="panel-body container">
                <div class="row">
                    <div class="col-lg-6 col-sm-10">
                        <p>花50积分参与本活动，作为押金积分</p>
                        <p>参与活动后，坚持每天早晨6-9点打卡</p>
                        <p>如果有1天中断，押金积分的50%将进入“公共奖金池”</p>
                        <p>连续坚持5天以上，将可以参与瓜分“公共奖金池”</p>
                        <p>中断后，可以补交后，次日继续打卡参与活动，重新计算累计，历史累计清零</p>
                        <p>中断后，如果没补交，每多中断一天，剩余押金将继续按50%比例扣除进入“公共奖金池”</p>
                        <p>押金扣到0为止</p>
                        <p>活动最终解释权归我公司所有</p>
                    </div>
                    <div class="col-lg-6 col-sm-10">
                        <div class="well well-sm">
                            <p>参与人数：<span class="label label-warning">1104 人</span></p>
                            <p>公共奖池：<span class="label label-danger">32991 积分</span></p>
                            <p>参与价格：<span class="label label-info">50 积分</span></p>
                            <p>预计分红：<span class="label label-success">3.5 积分/日</span></p>
                            <p>
                                &nbsp;<a class="btn btn-primary" type="button" href="/attendance/add">&nbsp;立即参与&nbsp;</a>
                                &nbsp;<a class="btn btn-primary" type="button" href="javascript:;">&nbsp;查看排行&nbsp;</a>
                            </p>
                            <p></p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="panel panel-info">
            <div class="panel-heading">打卡方案2-上下班阅读打卡</div>
            <div class="panel-body container">
                <div class="row">
                    <div class="col-lg-6 col-sm-10">
                        <p>花50积分参与本活动，作为押金积分</p>
                        <p>参与活动后，坚持每天早晨0-9点打卡，每天晚上18-24点打卡</p>
                        <p>打卡期间需要阅读一篇鸡汤</p>
                        <p>如果有1天中断，押金积分的50%将进入“公共奖金池”</p>
                        <p>连续坚持5天以上，将可以参与瓜分“公共奖金池”</p>
                        <p>中断后，可以补交后，次日继续打卡参与活动，重新计算累计，历史累计清零</p>
                        <p>中断后，如果没补交，每多中断一天，剩余押金将继续按50%比例扣除进入“公共奖金池”</p>
                        <p>押金扣到0为止</p>
                        <p>活动最终解释权归我公司所有</p>
                    </div>
                    <div class="col-lg-6 col-sm-10">
                        <div class="well well-sm">
                            <p>参与人数：<span class="label label-warning">1104人</span></p>
                            <p>公共奖池：<span class="label label-danger">32991积分</span></p>
                            <p>参与价格：<span class="label label-info">50积分</span></p>
                            <p>预计分红：<span class="label label-success">3.5积分/日</span></p>
                            <p>&nbsp;&nbsp;<a class="btn btn-primary col-lg-3 col-sm-6" type="button" href="/attendance/add">立即参与</a></p>
                            <p>&nbsp;</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="panel panel-success">
            <div class="panel-heading">打卡方案3-小时锻炼</div>
            <div class="panel-body container">
                <div class="row">
                    <div class="col-lg-6 col-sm-10">
                <p>花50积分参与本活动，作为押金积分</p>
                <p>参与活动后，坚持每小时1次打卡</p>
                <p>如果有1次中断，押金积分的50%将进入“公共奖金池”</p>
                <p>连续坚持5次以上，将可以参与瓜分“公共奖金池”</p>
                <p>中断后，可以补交后，下一小时继续打卡参与活动，重新计算累计，历史累计清零</p>
                <p>中断后，如果没补交，没多中断一天，剩余押金将继续按50%比例扣除进入“公共奖金池”</p>
                <p>押金扣到0为止</p>
                <p>活动最终解释权归我公司所有</p>
                    </div>
                    <div class="col-lg-6 col-sm-10">
                        <div class="well well-sm">
                            <p>参与人数：<span class="label label-warning">1104人</span></p>
                            <p>公共奖池：<span class="label label-danger">32991积分</span></p>
                            <p>参与价格：<span class="label label-info">50积分</span></p>
                            <p>预计分红：<span class="label label-success">3.5积分/日</span></p>
                            <p>&nbsp;&nbsp;<a class="btn btn-primary col-lg-3 col-sm-6" type="button" href="javascript:;">立即参与</a></p>
                            <p>&nbsp;</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>


    </div>
</div>
