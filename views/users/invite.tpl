{{append . "HeadStyles" "/static/css/custom.css"}}
{{append . "HeadScripts" "/static/js/custom.js"}}

{{template "alert.tpl" .}}
<div class="container text-center">
    <div class="row">
        <ul class="nav nav-tabs">
            <li><a href="/user/index">基本资料</a></li>
            <li><a href="/user/balance">积分</a></li>
            <li><a href="/attendance/my">打卡</a></li>
            <li><a href="/user/bonus">收益</a></li>
            <li class="active"><a href="javascript:;">推广码</a></li>
        </ul>
        <h4>我的邀请码</h4>
        <div class="well well-lg">
            邀请功能暂未上线
        </div>

    </div>
</div>
