/* input style implement begin */
$(".input-field input").on('focus', function () {
    $(this).parent().find("label").addClass("active");
});
$(".input-field input").on('blur', function () {
    if ($(this).val() == "")
        $(this).parent().find("label").removeClass("active");
});
/* input style implement end */

/* panel switch begin */
function show_history() {
    $("#create-guide").fadeIn();
    $("#create-guide").css('display', 'flex');
    $("#history-guide").css('display', 'none');
    $("#create-form").css('display', 'none');
    $("#history-form").fadeIn();
    $("#history-form").css('display', 'block');
}

function show_create() {
    $("#create-guide").css('display', 'none');
    $("#history-guide").fadeIn();
    $("#history-guide").css('display', 'flex');
    $("#create-form").fadeIn();
    $("#create-form").css('display', 'block');
    $("#history-form").css('display', 'none');
}

/* panel switch end */

var history_id = new Array();
var history_name = new Array();
var history_number;
history_name[0] = '才明洋表彰大会';
history_name[1] = '才明洋表彰大会';
history_name[2] = '才明洋表彰大会';
history_name[3] = '才明洋表彰大会';
history_name[4] = '才明洋表彰大会';
history_id[0] = '000';
history_id[1] = '001';
history_id[2] = '002';
history_id[3] = '003';
history_id[4] = '004';


history_number = 5;
var tem = 0;


function creat_from_history(e) {
    let data = {
        "action": "creat_from_history",
        "content": {
            "history_id": e.history_id
        }
    };
    $.ajax({
        url: '127.0.0.1:1923/get-activities',//老才的服务器名
        method: 'post',
        contentType: 'json',
        dataType: 'json',
        data: JSON.stringify(data),
        // success: function(e) {
        //     history_number=e.content.length;
        //     for(var i=0;i<e.content.length;i++){
        //         history_id[i]=content[i].id;
        //         history_name[i]=content[i].name;
        //     }
        //     console.log('获取历史记录成功');
        // },
        // error: function (e) {
        //     console.log(e);
        //     console.log('获取历史记录失败');
        // }
    })
}

function creat_from_blank() {
    let data = {
        "name": $('#input-activity-name').val(),
        "max-population": $('#max-population').val(),
        "remark": $('#validationTooltip03').val()
    };
    console.log('create from black');
    console.log(data);
    $.ajax({
        url: window.location.origin + '/append-activity',//老才的服务器名
        method: 'post',
        contentType: 'json',
        dataType: 'json',
        data: JSON.stringify(data),
        success: e => {
            window.location.href = "/console";
        }
    })
}

function resume_activity() {
    let data = {
        "action": "resume_activity",
        "content": ""
    };
    $.ajax({
        url: window.location.origin + '/get-activities',//老才的服务器名
        method: 'post',
        contentType: 'json',
        dataType: 'json',
        data: JSON.stringify(data),
    });
    window.location.href = "console";
}

function is_activity_unfinished(tem) {
    if (tem === 0) {
        $(function () {
            $("#resume").css('display', 'none');
            $("#create-or-history").css('display', 'block');
        });
    } else if (tem === 1) {
        $(function () {
            $("#resume").css('display', 'block');
            $("#create-or-history").css('display', 'none');
        });
    }
}

$(document).ready(function () {
    $.ajax({
        url: window.location.origin + '/is-activity-unfinished',
        method: 'get',
        success: function (e) {
            console.log(e);
            let tem = Number.parseInt(e);
            is_activity_unfinished(tem);
        },
        error: function (e) {
            console.log(e);
            is_activity_unfinished(0);
        }
    })
});