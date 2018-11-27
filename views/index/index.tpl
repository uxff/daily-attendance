{{append . "HeadStyles" "/static/css/custom.css"}}
{{append . "HeadScripts" "/static/js/custom.js"}}


<div class="container">
    <div class="row" style="">

    {{template "alert.tpl" .}}
        <div class="btn-group btn-group-justified" role="group" aria-label="...">
            <a href="javascript:;" class="btn btn-default" role="button">广告位招商 A</a>
            <a href="javascript:;" class="btn btn-default" role="button">广告位招商 B</a>
            <a href="javascript:;" class="btn btn-default" role="button">广告位招商 C</a>
        </div>
        <p></p>
    </div>

    <div class="row">

            <div class="panel panel-success">
                <div class="panel-heading">
                    <h3 class="panel-title">Welcome to Daily Attendance</h3>
                </div>
                <div class="panel-body">
                    Subscribe a Daily Attendance plan for your health.
                </div>
            </div>

    </div>
</div>
