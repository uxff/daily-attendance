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


    $('.da-timespan').hide();
    $('.da-timespan-end-4  option:last').attr('selected','selected');
    $('.da-timespan-end-5  option:last').attr('selected','selected');
    $('.da-timespan-end-6  option:last').attr('selected','selected');
    $('.da-timespan-end-7  option:last').attr('selected','selected');
    $('.da-timespan-5').show();

    $('.da-checkin-period').on('change', function(e) {
        var period = $('.da-checkin-period').val();
        console.log('da-checkin-period='+period);
        $('.da-timespan').hide();
        $('.da-timespan-'+period).show();
    });

    $('.da-btn-add-act').on('click', function (e) {
        var period = $('.da-checkin-period').val();
        var str = '';
        switch (period) {
            case "3":
                str = '{"'+period+'":{"dayspan":"0-59"}}';
                break;
            case "4":
                str = '{"'+period+'":{"timespan":"'+$('.da-timespan-start-'+period).val()+'-'+$('.da-timespan-end-'+period).val()+'"}}';
                break;
            case "5":
                str = '{"'+period+'":{"timespan":"'+$('.da-timespan-start-'+period).val()+'-'+$('.da-timespan-end-'+period).val()+'"}}';
                break;
            case "6":
                str = '{"'+period+'":{"dayspan":"'+$('.da-timespan-start-'+period).val()+'-'+$('.da-timespan-end-'+period).val()+'"}}';
                break;
            case "7":
                str = '{"'+period+'":{"dayspan":"'+$('.da-timespan-start-'+period).val()+'-'+$('.da-timespan-end-'+period).val()+'"}}';
                break;
        }
        console.log('str='+str);
        $('.da-checkin-rule').val(str)
    });
});
