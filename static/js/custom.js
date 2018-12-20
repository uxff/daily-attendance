if ($('#startTime').length > 0 && $('#endTime').length > 0) {
    $("#startTime").datetimepicker({
        format: 'yyyy-mm-dd hh:ii',
        minView:'month',
        language: 'zh-CN',
        autoclose:true,
        startDate:new Date()
    }).on("click",function(){
        $("#startTime").datetimepicker("setEndDate",$("#endTime").val())
    });
    $("#endTime").datetimepicker({
        format: 'yyyy-mm-dd hh:ii',
        minView:'month',
        language: 'zh-CN',
        autoclose:true,
        startDate:new Date()
    }).on("click",function(){
        $("#endTime").datetimepicker("setStartDate",$("#startTime").val())
    });
}

$(function () {
    console.log('x da defined.');
    $('.da-btn-all').on('click', function (e) {
        $('.da-tr').show();
    });
    $('.da-btn-ok').on('click', function (e) {
        $('.da-tr').hide();
        $('.da-tr-ok').show();
    });
    $('.da-btn-fail').on('click', function (e) {
        $('.da-tr').hide();
        $('.da-tr-fail').show();
    });
});
