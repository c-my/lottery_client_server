/* global variable begin */
let login_url = "login";
let bg_img_url = window.location.origin + "/post" // for debug
let drawing = false;
let running = false;
let started = false;
let isOnline = true;
let $cur_item = new Object();
let userType = "debug";
let max_winner = 1;
let cur_draw_times = 0;
let max_draw_times = 1;
let cheat_data = new Array();
/* global variable end */

/** WebSocket initialization begin */
let ws = new Object();
let connected = false;


function recv_ws_message(message) {
    console.log(message.data);
    let msg = JSON.parse(message.data);
    console.log(msg.action);
    switch (msg.action) {
        case 'append-user':
            let username = msg.content.nickname;
            let uid = msg.content.uid;
            let avatarUrl = msg.content.avatar;
            auto_add_user(avatarUrl, username, uid);
            break;
        case 'initialize':
            isOnline = (msg.content.online == "true");
            userType = msg.content.userType;
            console.log(isOnline, userType);
            if (isOnline) {
                $("#online-symbol").text("ONLINE");
                $('.import-participants').attr("disabled", true);
            } else {
                $("#online-symbol").text("OFFLINE");
            }

            /*init*/
            $('#bg-img').val(msg.content.url);
            /*基本信息*/
            $('#activity-name').val(msg.content.activity_name);
            /*奖池设置*/
            $('#draw-mode').val(msg.content.draw_mode_chosen)
            for (var i = 0; i < msg.content.reward_items_names.length; i++) {
                add_prize1(msg.content.reward_items_names[i]);
            }

            for (var i = 0; i < msg.content.prize_names.length; i++) {
                add_item1(msg.content.prize_names[i][0], msg.content.prize_names[i][1], msg.content.prize_names[i][2]);
            }
            /*方式主题*/

            // $('#show_style').attr("value",msg.content.topic_color);
            $("#show-style").setColor(msg.content.topic_color);

            $("#draw-style option[value=" + msg.content.lottery_style + "]").attr("selected", "selected");
            $("#draw-music option[value=" + msg.content.lottery_music + "]").attr("selected", "selected");
            $("#lucky-music option[value=" + msg.content.win_prize_music + "]").attr("selected", "selected");
            $("#reward-music option[value=" + msg.content.get_prize_music + "]").attr("selected", "selected");
            $("#show-style option[value=" + msg.content.topic_color + "]").attr("selected", "selected");
            /*弹幕设置*/
            // $('#font-size').val(msg.content.bullet_font_size);
            $('#font-size').scrollTo(msg.content.bullet_font_size);
            // $('#opacity').val(msg.content.bullet_transparent_degree);
            $('#opacity').scrollTo(msg.content.bullet_transparent_degree);

            // $('#font-color').val(msg.content.bullet_color);
            $('#font-color').setColor(msg.content.bullet_color);

            $('#danmu-switch').turn(msg.content.danmu_switch);
            $('#danmu-check-switch').turn(msg.content.danmu_check_switch);


            $("#font-family option[value=" + msg.content.bullet_font + "]").attr("selected", "selected");
            $("#danmu-speed option[value=" + msg.content.bullet_velocity + "]").attr("selected", "selected");
            $("#danmu-position option[value=" + msg.content.bullet_location + "]").attr("selected", "selected");
            $('.form-control-chosen').chosen();


            /*抽奖状态*/
            for (var i = 0; i < msg.content.reward_users_list.length; i++)
                show_lucky_dog(msg.content.reward_users_list[i].uid, msg.content.reward_users_list[i].nickname, msg.content.reward_users_list[i].itemname);
            ///TODO:剩余抽奖次数逻辑没写

            break;
        case 'who-is-lucky-dog':
            show_lucky_dog(msg.content.uid, msg.content.nickname, msg.content.itemname);
            console.log(msg.content.uid + "-" + msg.content.nickname + " gets " + msg.content.itemname);
            break;
        default:
            console.log('unknown action:\n' + msg.action);
    }
};
/** WebSocket initialization end */

/* input style implement begin */
$(".input-field input").on('focus', function () {
    $(this).parent().find("label").addClass("active");
});
$(".input-field input").on('blur', function () {
    if ($(this).val() == "")
        $(this).parent().find("label").removeClass("active");
});

/* input style implement end */

function auto_add_user(img_url, username, uid) {
    $("#parts-list-head").html("<tr><td>头像</td><td>用户名</td><td>UID</td></tr>");
    $("#parts-list-body").append("<tr><td>" + "<img width='60px' height='60px' src='" + img_url + "'/></td><td>" + username + "</td><td>" + uid + "</td></tr>");
}

/* start draw begin */
function start_draw() {
    $("#draw-label").text("停止抽奖");
    $(".btn-draw-action").css("background-color", "#dc3545");
    $("#draw-action").attr("class", "mdi mdi-stop-circle-outline");
    // 禁用抽奖作弊功能
    if (userType != "normal")
        $cur_item.find(".cfg-cheat").attr("disabled", true);
    // 禁用抽奖数量设置
    $cur_item.find(".max_winner").attr("readonly", "readonly");
    // 禁用奖项列表中的奖品设置
    $cur_item.find(".item-prize").addClass("disable");
    runNotify({
        message: $cur_item.find(".item-name").val() + "第" + (cur_draw_times + 1)
            + "次抽奖，剩余" + (max_draw_times - cur_draw_times - 1) + "次",
        messageTitle: 'title',
        levelMessage: 'info',
        timer: '2500'
    });
    if (true) ws.send(JSON.stringify({
        action: 'start-drawing',
        content: {'dkind': $("#draw-style").val()}
    }));
}

/* start draw end */

/* stop draw begin */
function stop_draw() {
    $("#draw-label").text("开始抽奖");
    $(".btn-draw-action").css("background-color", "##28a745");
    $("#draw-action").attr("class", "mdi mdi-arrow-right-drop-circle-outline");
    // 启用抽奖作弊功能
    if (userType != "normal") $cur_item.find(".cfg-cheat").attr("disabled", false);
    // 启用抽奖数量设置
    $cur_item.find(".max_winner").removeAttr("readonly");
    // 启用奖项列表中的奖品设置
    $cur_item.find(".item-prize").removeClass("disable");
    cur_draw_times++;
    if (cur_draw_times == max_draw_times) {
        $(".btn-draw-action").prop("disabled", true);
    }
    if (true) ws.send(JSON.stringify({
        action: 'stop-drawing',
        content: ""
    }));
}

/* stop draw end */

/* draw button clicked begin */
function on_draw_btn_click() {
    if (!drawing) start_draw(); else stop_draw();
    drawing = !drawing;
}

/* draw button clicked end */

/* activity button clicked begin */
function on_activity_btn_click() {
    if (!running) {
        $("#activity-label").text("隐藏活动");
        $(".btn-activity-action").css("background-color", "#dc3545");
        $("#activity-action").attr("class", "mdi mdi-stop-circle-outline");
        max_draw_times = parseInt($cur_item.find(".max_winner").val()); // 当前奖项的最大可抽次数
        console.log("max_draw_times: ", max_draw_times);
        $cur_item.find("input.item-name").attr("readonly", "readonly"); // 禁用当前奖项在列表中改名
        $cur_item.find(".del-item").attr("disabled", true); // 禁用当前奖项在列表中的删除
        cur_draw_times = 0; // 初始化已抽奖次数为0
        $(".btn-draw-action").prop('disabled', false); // 启用开始抽奖按钮
        $("#cur-item-input").attr("disabled", true); // 禁用输入当前奖项
        if (!started) {
            var t = new Date();
            var now = t.getFullYear() + "-" + t.getMonth() + "-" + t.getDate() + "-" +
                t.getHours() + "-" + t.getMinutes() + "-" + t.getSeconds();
            if (true) ws.send(JSON.stringify({
                action: 'activity-start-time',
                content: now
            }));
            started = true;
        }
        if (true) ws.send(JSON.stringify({
            action: 'show-activity',
            content: {
                cur_item: $cur_item.find("input.item-name").val()
            }
        }));
    } else {
        $("#activity-label").text("展示活动");
        $(".btn-activity-action").css("background-color", "#006400");
        $("#activity-action").attr("class", "mdi mdi-arrow-right-drop-circle-outline");
        $(".btn-draw-action").prop('disabled', true);
        // 清空并启用输入当前奖项
        $("#cur-item-input").val("");
        $("#cur-item-input").blur();
        $("#cur-item-input").attr("disabled", false);
        // 启用奖项列表的删除对应项功能
        $cur_item.find(".del-item").attr("disabled", false);
        // 禁用抽奖作弊功能
        if (userType != "normal") $cur_item.find(".cfg-cheat").attr("disabled", true);
        // 禁用抽奖数量设置
        $cur_item.find(".max_winner").attr("readonly", "readonly");
        // 禁用奖项列表中的奖品设置
        $cur_item.find(".item-prize").addClass("disable");
        $(".btn-activity-action").prop("disabled", true);
        if (drawing) stop_draw();
        if (true) ws.send(JSON.stringify({
            action: 'hide-activity',
            content: ''
        }));
    }
    running = !running;
}

/* activity button clicked end */

/* csv parser begin */
$.extend({
    csv: function (record, f) {
        record = record.split(/\n/);
        var title = record[0].split(",");
        record.shift();
        var data = [];
        for (var y = 0; y < title.length; y++) {
            title[y] = title[y].trim();
        }
        for (var i = 0; i < record.length; i++) {
            var t = $.trim(record[i]).split(",");
            for (var y = 0; y < t.length; y++) {
                if (!data[i]) data[i] = {};
                data[i][title[y]] = t[y];
            }
        }
        f.call(this, data);
        data = null;
    }
});
/* csv parser end */

/* import participants begin */
function import_participants() {
    var file = document.getElementById("participants").files[0];
    // var fileName = file.name.substring(0, file.name.lastIndexOf('.'));
    var reader = new FileReader();
    reader.onload = function (f) {
        $('#parts').val($('#participants').val());
        $('#participants').val("");
        $("#parts-list-head").empty();
        $("#parts-list-body").empty();
        $('.parts-label').addClass("active");
        var data = this.result;
        $(function () {
            $.csv(data, function (data) {
                if (!data) return;
                if (data && !data[0].hasOwnProperty('username')) {
                    for (var item in data[0]) console.log(item);
                    runNotify({
                        message: 'CSV必须包含username字段',
                        messageTitle: 'title',
                        levelMessage: 'warn',
                        timer: '2000'
                    });
                    $('#parts').val("");
                    $('#parts').blur();
                    return;
                }
                var head_str = "<tr>";
                for (item in data[0])
                    head_str += "<th>" + item + "</th>";
                $("#parts-list-head").append(head_str + "</tr>");
                for (var i = 0; i < data.length; ++i) {
                    var record_str = "<tr>";
                    var content = {};
                    for (item in data[0]) {
                        record_str += "<td>" + data[i][item] + "</td>";
                        content[item] = data[i][item];
                    }
                    if (true) ws.send(JSON.stringify({
                        action: 'manual-import',
                        content: content
                    }));
                    $("#parts-list-body").append(record_str + "</tr>");
                }
            });
        });
    };
    reader.readAsText(file);
}

/* import participants end */

/* import background image begin */
function getObjectURL(file) {
    var url = null;
    if (window.createObjectURL !== undefined) { // basic
        url = window.createObjectURL(file);
    } else if (window.URL !== undefined) { // mozilla(firefox)
        url = window.URL.createObjectURL(file);
    } else if (window.webkitURL !== undefined) { // webkit or chrome
        url = window.webkitURL.createObjectURL(file);
    }
    return url;
}

function converImgTobase64(url, callback, outputFormat) {
    var canvas = document.createElement('CANVAS'),
        ctx = canvas.getContext('2d'),
        img = new Image;
    img.crossOrigin = 'Anonymous';
    img.onload = function () {
        canvas.height = img.height;
        canvas.width = img.width;
        ctx.drawImage(img, 0, 0);
        var dataURL = canvas.toDataURL(outputFormat || 'image/png');
        callback.call(this, dataURL);
        canvas = null;
    };
    img.src = url;
}

function import_bg_image() {
    var input = $("#background-img")[0],
        src = getObjectURL(input.files[0]);
    // src = input.files[0];
    // var src = document.getElementById("background-img").files[0];
    console.log("src: ", src);
    converImgTobase64(src, function (base64Str) {
        $.ajax({
            type: "POST",
            url: bg_img_url,
            dataType: "json",
            data: JSON.stringify({
                action: "background-image",
                content: {
                    base64Str: base64Str
                }
            }),
            success: function (result) {
                console.log("send image: " + result);
            }
        });
    });
    $('.bg-img-label').addClass("active");
    $('#bg-img').val($('#background-img').val());
    $('#background-img').val("");
}

/* import background image end */

/* ecalendar plugin initialize begin */
var eCalendar = {
    type: "time",   //模式，time: 带时间选择; date: 不带时间选择;
    stamp: false,    //是否转成时间戳，默认true;
    offset: [0, 2],   //弹框手动偏移量;
    format: "yyyy年mm月dd日 hh:ii",   //时间格式 默认 yyyy-mm-dd hh:ii;
    skin: "#1B7690",   //皮肤颜色，默认随机，可选值：0-8,或者直接标注颜色值;
    step: 10,   //选择时间分钟的精确度;
    callback: function (v, e) {
    } //回调函数
};
/* ecalendar plugin initialize end */

/* prize pool table operation begin */

/* add a prize begin*/
function add_prize() {
    if ($('input:focus').length != 0) return;
    if ($('#prize-list-body').find("tr").length == 0 || $('#prize-list-body tr:last').find('input').val()) {
        $("#prize-list-body").append(`<tr>
            <td><input class="table-input" style="margin-bottom: 0px;" type="text" draggable="true" ondragstart="drag_prize(event)" onblur="check_content(this)"></input></td>
            <td class="icon-td"><div class="del-prize table-icon" onclick="del_row(this)"><i class="mdi mdi-trash-can-outline"></i></div></td>
            </tr>`);
        $(".table-input").attr("onkeydown", "input_keydown(this, event)");
    }
    $('#prize-list-body tr:last').find('input').focus();
}

/* add a prize end */

/* add an item begin */
function add_item() {
    if ($('input:focus').length != 0) return;
    if ($('#item-list-body').find("tr").length == 0 || $('#item-list-body tr:last').find('input.item-name').val()) {
        var cheat_functions = `<div class="cfg-cheat table-icon" onclick="cfg_cheat(this)"><i class="mdi mdi-settings"></i></div>`
        if (userType == "normal") cheat_functions = '';
        $("#item-list-body").append(`<tr>
            <td><input class="table-input item-name" type="text" onblur="check_content(this)"></input></td>
            <td><input class="item-prize" readonly="readonly" ondrop="drop_prize(event)" ondragover="allowDrop(event)"></input></td>
            <td><input class="table-input max_winner" type="number" min="1" value="1"></input></td>
            <td class="icon-td">
            ` + cheat_functions + `
            <div class="del-item table-icon" onclick="del_item(this)"><i class="mdi mdi-trash-can-outline"></i></div>
            </td>
            </tr>`);
        $(".table-input").attr("onkeydown", "input_keydown(this, event)");
        $(".cfg-cheat").attr("onmouseover", "show_cheat_info(this, event)");
        $(".cfg-cheat").attr("onmouseout", "hide_cheat_info(this, event)");
        $(".max_winner").attr("onchange", "update_max_draw_times(this, event)");
    }
    $('#item-list-body tr:last').find('input.item-name').focus();
}

/* add an item end */

/* update max_draw_times begin */
function update_max_draw_times(obj, event) {
    if (running && $(obj).parent().parent().find("input.item-name").val() ==
        $cur_item.find("input.item-name").val()) {
        if (max_draw_times == cur_draw_times && parseInt($(obj).val()) > max_draw_times) {
            $(".btn-draw-action").prop("disabled", false);
        } else if (parseInt($(obj).val()) < cur_draw_times) {
            $(obj).val(cur_draw_times);
            runNotify({
                message: '数量不得小于已进行抽奖次数',
                messageTitle: 'title',
                levelMessage: 'warn',
                timer: '2000'
            });
        }
        max_draw_times = $(obj).val();
        console.log("max_draw_times updated to: ", max_draw_times);
    }
}

/* update max_draw_times end */

/* check new item content begin */
function check_content(obj) {
    if (obj.value == "") {
        del_row(obj);
        return false;
    }
    var tr = $(obj).parent().parent();
    var obj_index = tr.prevAll().length;
    $(tr.parent()).find("tr").each(function (index, element) {
        var input = $(this).children("td:first").find("input").val();
        if (index != obj_index && obj.value == input) {
            runNotify({
                message: '内容重复',
                messageTitle: 'title',
                levelMessage: 'error',
                timer: '1000'
            });
            obj.focus();
            return false;
        }
    });
    return true;
}

/* check new item content end */

/* delete a row begin */
function del_row(obj) {
    var tr = obj.parentNode.parentNode;
    tr.parentNode.removeChild(tr);
}

/* delete a row end */

/* delete an item begin */
function del_item(obj) {
    if ($(obj).attr("disabled")) {
        console.log("delete disabled");
        return;
    }
    console.log("delete an item");
    var item_name = $(obj).parent().parent().find(".item-name").val();
    if (item_name != '' && item_name in cheat_data) {
        delete cheat_data[item_name];
        // console.log(cheat_data);
    }
    del_row(obj);
}

/* delete an item end */

/* prize-name draggable begin */
function drag_prize(event) {
    event.dataTransfer.setData("Text", event.target.value);
}

function drop_prize(event) {
    if ($(event.target).hasClass("disable")) return;
    event.preventDefault();
    var data = event.dataTransfer.getData("Text");
    event.target.value = data;
}

function allowDrop(event) {
    event.preventDefault();
}

/* prize-name draggable begin */
/* prize pool table operation end */

/* range input filling begin */
function range_input() {
    var range = parseInt(parseInt(this.max) - parseInt(this.min));
    var percent = String(parseInt(parseInt(this.value) - parseInt(this.min)) * 100 / range);
    $(this).css('background', 'linear-gradient(to right, #343A40, #ebeff4 ' + percent + '%, #ebeff4)');
    $(this).next().text(String(this.value));
}

$.fn.RangeSlider = function () {
    var $input = $(this);
    $input.attr('min', this.min).attr('max', this.max).attr('step', 1);
    $input.bind("input", range_input);
};
/* range input filling end */

/* get lucky draw items for choice begin*/
function get_items() {
    if ($(this).attr("disabled")) {
        console.log("current item input disabled");
        $(this).blur();
        return;
    }
    $(this).val("");
    var items = "";
    $("#item-list-body tr").each(function () {
        var $input = $(this).children("td:first").find("input");
        if ($input.val() && typeof ($input.attr("readonly")) == "undefined") {
            items += '<option value="' + $input.val() + '" />';
        }
    });
    if (items == "") {
        runNotify({
            message: '请先设置奖项',
            messageTitle: 'title',
            levelMessage: 'warn',
            timer: '1000'
        });
        $(this).blur();
    }
    $("#items").html(items);
    return false;
}

/* get lucky draw items for choice end*/

/* check current item begin */
function check_cur_item() {
    var cur_itemname = $(this).val();
    if (cur_itemname != "") {
        $('#item-list-body').find("tr").each(function (index, element) {
            if ($(this).find("input.item-name").val() == cur_itemname) {
                $cur_item = $(this);
                console.log("Get global variable $cur_item: ", $cur_item);
            }
        });
        $(".btn-activity-action").prop('disabled', false);
    } else {
        $(".btn-activity-action").prop('disabled', true);
    }
}

/* check current item end */

/* show lucky dog begin */
function show_lucky_dog(uid, username, itemname) {
    $("#lucky-list-body").append(`<tr><td width="45%" class="username" uid="` + uid +
        `">` + username + `</td><td width="40%" class="itemname">` + itemname +
        `<td class="icon-td" width="15%">
        <div class="disable-lucky table-icon" onclick="disable_lucky(this, event)">
            <i class="mdi mdi-account-remove"></i>
        </div>
    </td>`);
}

/* disable lucky dog begin */
function disable_lucky(obj, event) {
    var $tr = $(obj).parent().parent();
    var $user = $tr.find("td.username");
    var username = $user.text();
    var uid = $user.attr("uid");
    var itemname = $tr.find("td.itemname").text();
    if (true) ws.send(JSON.stringify({
        action: "disable-lucky",
        content: {
            uid: uid,
            username: username,
            itemname: itemname
        }
    }));
    del_row(obj);
    // add user to cheat-loser table
    if (userType != "normal") {
        $('#item-list-body').find("tr").each(function (index, element) {
            if ($(this).find("input.item-name").val() == itemname) {
                if (!(itemname in cheat_data)) cheat_data[itemname] = {winner: [], loser: []};
                var hasAlready = false;
                for (var i = 0; i < cheat_data[itemname].loser.length; i++) {
                    if (cheat_data[itemname].loser[i] == username) {
                        hasAlready = true;
                        break;
                    }
                }
                if (!hasAlready) cheat_data[itemname].loser.push(username);
            }
        });
    }
}

/* disable lucky dog end */
/* show lucky dog end */

/* cheat function begin */
function cfg_cheat(obj) {
    if ($(obj).attr("disabled")) {
        console.log("cfg-cheat disabled");
        return;
    }
    var item_name = $(obj).parent().parent().find(".item-name").val();
    max_winner = $(obj).parent().parent().find(".max_winner").val();
    $("#cheat-kind").text(item_name);
    if (item_name in cheat_data) {
        cheat_data[item_name].winner.forEach(function (value, i) {
            $("#cheat-winner").append(`<tr>
            <td><input class="table-input" value="` + value + `"style="margin-bottom: 0px;" type="text" onblur="check_content(this)"></input></td>
            <td class="icon-td"><div class="del-winner table-icon" onclick="del_row(this)"><i class="mdi mdi-trash-can-outline"></i></div></td>
            </tr>`);
        });
        cheat_data[item_name].loser.forEach(function (value, i) {
            $("#cheat-loser").append(`<tr>
            <td><input class="table-input" value="` + value + `" style="margin-bottom: 0px;" type="text" onblur="check_content(this)"></input></td>
            <td class="icon-td"><div class="del-loser table-icon" onclick="del_row(this)"><i class="mdi mdi-trash-can-outline"></i></div></td>
            </tr>`);
        });
        $(".table-input").attr("onkeydown", "input_keydown(this, event)");
    }
    $("#cheat-modal").modal('show');
}

function show_cheat_info(obj, event) {
    if ($(obj).attr("disabled")) {
        console.log("show_cheat_info disabled");
        return;
    }
    // console.log("hover");
    var max_len = 3;
    var item_name = $(obj).parent().parent().find(".item-name").val();
    if (item_name in cheat_data) {
        // console.log("show hover");
        if (cheat_data[item_name].winner.length == 0) {
            $("#cheat-winner-info").append(`<tr><td>无</td></tr>`);
        } else {
            for (var i = 0; i < cheat_data[item_name].winner.length; i++) {
                $("#cheat-winner-info").append(`<tr><td>` + cheat_data[item_name].winner[i] + `</td></tr>`);
                if (i + 1 >= max_len) break;
            }
            if (cheat_data[item_name].winner.length >= max_len) {
                $("#cheat-winner-info").append(`<tr><td>. . .</td></tr>`);
            }
        }
        if (cheat_data[item_name].loser.length == 0) {
            $("#cheat-loser-info").append(`<tr><td>无</td></tr>`);
        } else {
            for (var i = 0; i < cheat_data[item_name].loser.length; i++) {
                $("#cheat-loser-info").append(`<tr><td>` + cheat_data[item_name].loser[i] + `</td></tr>`);
                if (i + 1 >= max_len) break;
            }
            if (cheat_data[item_name].loser.length >= max_len) {
                $("#cheat-loser-info").append(`<tr><td>. . .</td></tr>`);
            }
        }
        $("#cheat-info").show();
    }
}

function hide_cheat_info(obj, event) {
    $("#cheat-info").hide();
    $("#cheat-winner-info").empty();
    $("#cheat-loser-info").empty();
}

function echo_info(statebar, state, content) {
    var prefix = "";
    var _class = "alert alert-dismissible ";
    if (state == 0) {
        prefix = "成功";
        _class += "alert-success";
    } else if (state == 1) {
        prefix = "警告";
        _class += "alert-warning";
    } else if (state == 2) {
        prefix = "错误";
        _class += "alert-danger";
    }
    $(statebar).html("<strong>" + prefix + ": </strong>" + content);
    $(statebar).attr("class", _class);
    $(statebar).show();
    setTimeout(
        function () {
            $(statebar).hide();
        }, 2000
    );
}

function check_cheat_content(obj) {
    if (obj.value == "") {
        del_row(obj);
        return false;
    }
    var tr = $(obj).parent().parent();
    var obj_index = tr.prevAll().length;

    $(tr.parent()).find("tr").each(function (index, element) {
        var input = $(this).children("td:first").find("input").val();
        if (index != obj_index && obj.value == input) {
            echo_info("#cheat-statebar", 2, "用户ID存在冲突");
            obj.focus();
            return false;
        }
    });
    var other = tr.parent().attr('id') == "cheat-winner" ? "#cheat-loser" : "#cheat-winner";
    // console.log(other);
    $(other).find("tr").each(function (index, element) {
        var input = $(this).children("td:first").find("input").val();
        if (obj.value == input) {
            echo_info("#cheat-statebar", 2, "用户ID存在冲突");
            obj.focus();
            return false;
        }
    });
    return true;
}

function add_winner() {
    if ($('input:focus').length != 0) return;
    if ($("#cheat-winner").find("tr").length >= max_winner) {
        echo_info("#cheat-statebar", 1, "中奖人选超出上限");
        return;
    }
    $("#cheat-winner").append(`
        <tr draggable="true" ondragstart="drag_cheat_item(this, event)">
            <td><input class="table-input" style="margin-bottom: 0px;" type="text" onblur="check_cheat_content(this)"></input></td>
            <td class="icon-td"><div class="del-winner table-icon" onclick="del_row(this)"><i class="mdi mdi-trash-can-outline"></i></div></td>
        </tr>`);
    $(".table-input").attr("onkeydown", "input_keydown(this, event)");
    $('#cheat-winner tr:last').find('input').focus();
}

function add_loser() {
    if ($('input:focus').length != 0) return;
    $("#cheat-loser").append(`
        <tr draggable="true" ondragstart="drag_cheat_item(this, event)">
            <td><input class="table-input" style="margin-bottom: 0px;" type="text" onblur="check_cheat_content(this)"></input></td>
            <td class="icon-td"><div class="del-loser table-icon" onclick="del_row(this)"><i class="mdi mdi-trash-can-outline"></i></div></td>
        </tr>`);
    $(".table-input").attr("onkeydown", "input_keydown(this, event)");
    $('#cheat-loser tr:last').find('input').focus();
}

let drag_element = new Object();

function drag_cheat_item(obj, event) {
    $('input:focus').blur();
    drag_element = obj;
}

function drop_cheat_item(obj, event) {
    var $tbody = $(obj).find("table tbody");
    if ($tbody.attr("id") == "cheat-winner" && $tbody.find("tr").length >= max_winner) {
        echo_info("#cheat-statebar", 1, "中奖人选超出上限");
        return;
    }
    event.preventDefault();
    $tbody.append($(drag_element).clone());
    drag_element.parentNode.removeChild(drag_element);
}

function allowDrop(event) {
    event.preventDefault();
}

function save_cheat_info() {
    var item_name = $("#cheat-kind").text();
    cheat_data[item_name] = {winner: [], loser: []};
    $("#cheat-winner tr").each(function () {
        var text = $(this).children("td:first").find("input").val();
        cheat_data[item_name].winner.push(text);
    });
    $("#cheat-loser tr").each(function () {
        var text = $(this).children("td:first").find("input").val();
        cheat_data[item_name].loser.push(text);
    });
    if (cheat_data[item_name].winner.length == 0 && cheat_data[item_name].loser.length == 0) {
        delete cheat_data[item_name];
    }
    // console.log(cheat_data[item_name]);
    quit_cheat_cfg();
}

function quit_cheat_cfg() {
    $("#cheat-modal").modal('hide');
    $("#cheat-winner").empty();
    $("#cheat-loser").empty();
}

/* cheat function end */

/* switch stage page begin */
function switch_page() {
    $("#page-name").text($(this).text());
    if (true) ws.send(JSON.stringify({
        action: "switch-page",
        content: $("#page-name").text()
    }));
}

/* switch stage page end */

/* update settings function begin */
function update_setting(obj, event) {
    // $("#page-swichbar li .active").text()
    // 获得当前选项卡的另一种方法
    var href = $("#page-swichbar").find(".active").attr("href");
    var content = {}; // {"part": href}; 才哥说不要part
    if (href == "#basic-info") {
        content['activity-name'] = $("#activity-name").val();
        console.log("update #basic-info");
    } else if (href == "#prize-pool") {
        content['draw-mode'] = $("#draw-mode").val();
        console.log("update #prize-pool");
    } else if (href == "#style-theme") {
        content['draw-style'] = $("#draw-style").val();
        content['theme-color'] = $("#show-style").val();
        content['draw-music'] = $("#draw-music").val();
        if ($("#bg-img").val() == "") {
            content['background-img'] = "none";
        } else if ($("#bg-img").val().indexOf("fakepath") != -1) {
            content['background-img'] = "base64Str";
        } else {
            content['background-img'] = $("#bg-img").val();
        }
        content['lucky-music'] = $("#lucky-music").val();
        content['reward-music'] = $("#reward-music").val();
        console.log("update #style-theme");
    } else if (href == "#danmu-sets") {
        content['font-size'] = $("#font-size").val();
        content['opacity'] = $("#opacity").val();
        content['font-family'] = $("#font-family").val();
        content['font-color'] = $("#font-color").val();
        content['danmu-speed'] = $("#danmu-speed").val();
        content['danmu-position'] = $("#danmu-position").val();
        console.log("update #danmu-sets");
    }
    console.log(content);
    if (true) ws.send(JSON.stringify({
        action: "part-update",
        content: content
    }));
    runNotify({
        message: '已成功更新同步',
        messageTitle: 'title',
        levelMessage: 'info',
        timer: '2000'
    });
}

/* update settings function end */

/* finish activity function begin */
function confirm_finish_confirm() {
    var t = new Date();
    var now = t.getFullYear() + "-" + t.getMonth() + "-" + t.getDate() + "-" +
        t.getHours() + "-" + t.getMinutes() + "-" + t.getSeconds();
    if (true) ws.send(JSON.stringify({
        action: 'activity-end-time',
        content: now
    }));
    if (navigator.userAgent.indexOf("Firefox") != -1 ||
        navigator.userAgent.indexOf("Chrome") != -1) {
        window.opener = null;
        window.open(login_url, "_self");
        window.close();
    } else {
        window.opener = null;
        window.open("", "_self");
        window.close();
    }
}

function save_configurations() {
    ws.send(JSON.stringify({
        action: "save_as_history",
        content: ""
    }));


    echo_info("#save-statebar", 0, "设置信息已保存");
}

function quit_finish_confirm() {
    $("#finish-confirm").modal('hide');
}

function finish_activity() {
    $("#finish-confirm").modal('show');
}

/* finish activity function end */

/* danmu switch function begin */
function danmu_switch(event, state) {
    if (true) ws.send(JSON.stringify({
        action: "danmu-switch",
        content: state
    }));
    console.log("danmu switch: ", state);
}

function danmu_check_switch(event, state) {
    if (true) ws.send(JSON.stringify({
        action: "danmu-check-switch",
        content: state
    }));
    console.log("danmu check switch: ", state);
}

/* danmu switch function end */

/// 小工具
String.prototype.format = function () {
    var args = Array.prototype.slice.call(arguments);
    var count = 0;
    return this.replace(/%s/g, function (s, i) {
        return args[count++];
    });

}

/// chosen插件的选择功能
jQuery.fn.choose = function (option) {
    var chosen_id = $(this).attr("id").replace('-', '_') + "_chosen";
    var chosen_index = -1;
    $(this).find("option").each(function (index, element) {
        if ($(element).val() == option) chosen_index = index;
    });
    if (chosen_index == -1) return;
    $(("#%s a.chosen-single").format(chosen_id)).mousedown();
    $(("#%s li[data-option-array-index=%s]").format(chosen_id, chosen_index)).mouseup();
}


/// Range Input设置值
jQuery.fn.scrollTo = function (value) {
    $(this).val(value);
    $(this).trigger('input');
}

/// switch 开关
jQuery.fn.turn = function (state) {
    console.log("turn " + state);
    var new_state = state == "on";
    var cur_state = $(this).bootstrapSwitch("state");
    if ((new_state && !cur_state) || (!new_state && cur_state)) {
        $(this).bootstrapSwitch("toggleState");
    }
}

/// minicolor
jQuery.fn.setColor = function (color) {
    $(this).val(color);
    $(this).trigger('keyup');
}


/* shotcut key functions begin */
function global_keydown(e) { // global

}

function input_keydown(obj, event) { // input
    if (event.key == 'Enter') {
        $(obj).blur();
        // console.log("blured");
    }
}

/* shotcut key functions end */

//-------------------------------bjz-begin---------------------------------
//-------------------------------bjz-begin---------------------------------
//-------------------------------bjz-begin---------------------------------
//-------------------------------bjz-begin---------------------------------

function add_item1(a, b, c) {
    if ($('input:focus').length != 0) return;
    if ($('#item-list-body').find("tr").length == 0 || $('#item-list-body tr:last').find('input.item-name').val()) {
        var cheat_functions = `<div class="cfg-cheat table-icon" onclick="cfg_cheat(this)"><i class="mdi mdi-settings"></i></div>`
        if (userType == "normal") cheat_functions = '';
        $("#item-list-body").append(`<tr>
            <td><input class="table-input item-name" type="text" onblur="check_content(this)"></input></td>
            <td><input class="item-prize" readonly="readonly" ondrop="drop_prize(event)" ondragover="allowDrop(event)"></input></td>
            <td><input class="table-input max_winner" type="number" min="1" value="1"></input></td>
            <td class="icon-td">
            ` + cheat_functions + `
            <div class="del-item table-icon" onclick="del_item(this)"><i class="mdi mdi-trash-can-outline"></i></div>
            </td>
            </tr>`);
        $(".table-input").attr("onkeydown", "input_keydown(this, event)");
        $(".cfg-cheat").attr("onmouseover", "show_cheat_info(this, event)");
        $(".cfg-cheat").attr("onmouseout", "hide_cheat_info(this, event)");
        $(".max_winner").attr("onchange", "update_max_draw_times(this, event)");
    }
    $('#item-list-body tr:last').find('input.item-name').val(a);
    $('#item-list-body tr:last').find('input.item-prize').val(b);
    $('#item-list-body tr:last').find('input.max_winner').val(c);

}

function add_prize1(str) {
    if ($('input:focus').length != 0) return;
    if ($('#prize-list-body').find("tr").length == 0 || $('#prize-list-body tr:last').find('input').val()) {
        $("#prize-list-body").append(`<tr>
            <td><input class="table-input" style="margin-bottom: 0px;" type="text" draggable="true" ondragstart="drag_prize(event)" onblur="check_content(this)"></input></td>
            <td class="icon-td"><div class="del-prize table-icon" onclick="del_row(this)"><i class="mdi mdi-trash-can-outline"></i></div></td>
            </tr>`);
        $(".table-input").attr("onkeydown", "input_keydown(this, event)");
    }
    $('#prize-list-body tr:last').find('input').val(str);
}

function getQueryString(name) {//获取name参数的值
    var result = window.location.search.match(new RegExp("[\?\&]" + name + "=([^\&]+)", "i"));
    if (result == null || result.length < 1) {
        return "";
    }
    return result[1];
}

var path = getQueryString('path');
// console.log(path);
// var fso=new ActiveXObject(Scripting.FileSystemObject);
// var f=fso.opentextfile(path,1,true);
// while (!f.AtEndOfStream) {
//     console.log(f.Readline());
// }
// f.close();
function save_as_history() {///保存为历史设置
    if (true) ws.send(JSON.stringify({
        action: "save_as_history",
        content: ""
    }));
}

function pause() {///临时暂停活动
    if (true) ws.send(JSON.stringify({
        action: "pause_this_activity",
        content: ""
    }));
}

function save_information_as_history() {///保存活动信息
    let data = {
        activity_id: 'tem_activity_id',
        activity_name: 'tem_activity_name',
        activity_maxnum: 'tem_activity_maxnum',
        reward_information: 'tem_reward_information',
        // password: $('#inputPassword').val(),
    };
    $.ajax({
        url: 'signin',//老才的服务器名
        method: 'POST',
        contentType: 'json',
        dataType: 'json',
        data: JSON.stringify(data),
        success: function (e) {
            console.log('保存成功');
        },
        error: function (e) {
            console.log(e);
        }
    })

}

//-------------------------------bjz-end---------------------------------
//-------------------------------bjz-end---------------------------------
//-------------------------------bjz-end---------------------------------
//-------------------------------bjz-end---------------------------------

/* initialization and bindings */
$(document).ready(function () {

    /* disable browser backward */
    // history.pushState(null, null, document.URL);
    // window.addEventListener('popstate', function () {
    //     history.pushState(null, null, document.URL);
    // });

    /* Import participants binding */
    $("#participants").change(import_participants);

    /* Date input binding */
    // $("#start-time").ECalendar(eCalendar);
    // $("#end-time").ECalendar(eCalendar);

    /* font-family list initialize */ // must before select input init
    $("#font-family option").each(function () {
        $(this).css("font-family", $(this).val());
    });

    /* select input plugin initialization */

    /* add prize binding */
    $(".btn-add-prize").click(add_prize);
    $(".btn-add-item").click(add_item);

    /* color picker initialization */
    $('INPUT.minicolors').minicolors({letterCase: 'uppercase'});

    /* range input initialization */
    $('#font-size').RangeSlider();
    $('#opacity').RangeSlider();

    /* automatically get item when click cur-item-input */
    $("#cur-item-input").on('focus', get_items);
    $("#cur-item-input").on('blur', check_cur_item);

    /* draw and activity button clicked */
    $(".btn-draw-action").click(on_draw_btn_click);
    $(".btn-activity-action").click(on_activity_btn_click);

    /* cheat binding */
    $("#btn-add-winner").click(add_winner);
    $("#btn-add-loser").click(add_loser);

    /* switch page binding */
    $(".dropdown-item.page").click(switch_page);

    /* import background image binding */
    $("#background-img").change(import_bg_image);

    /* swich initialization */
    $("[name='switch']").bootstrapSwitch();
    $("#danmu-switch").on('switchChange.bootstrapSwitch', danmu_switch);
    $("#danmu-check-switch").on('switchChange.bootstrapSwitch', danmu_check_switch);

    /* shotcut key map */
    document.onkeydown = global_keydown;
    $("#cur-item-input").attr("onkeydown", "input_keydown(this, event)");

    ws = new WebSocket('ws://'+window.location.host+'/ws');

    let test_data = {data: JSON.stringify({
        action: 'initialize',
        content: {
            reward_items_names: ['一等奖', '二等奖'],
            prize_names: [
                ['一等奖', '奔驰汽车', 10],
                ['二等奖', '漏油的奔驰汽车', 10]
            ]
        }
    })};
    console.log(test_data);
    ws.onmessage = recv_ws_message;
    recv_ws_message(test_data);

    if (!$.isEmptyObject(ws) && ws.readyState == 1) {
        connected = true;
        ws.onmessage = recv_ws_message;
    }

});

String.prototype.format = function () {
    var args = Array.prototype.slice.call(arguments);
    var count = 0;
    return this.replace(/%s/g, function (s, i) {
        return args[count++];
    });
}
/// minicolor
jQuery.fn.setColor = function (color) {
    $(this).val(color);
    $(this).trigger('keyup');
}


/// Range Input设置值
jQuery.fn.scrollTo = function (value) {
    $(this).val(value);
    $(this).trigger('input');
}

/// switch 开关
jQuery.fn.turn = function (state) {
    console.log("turn " + state);
    var new_state = state == "on";
    var cur_state = $(this).bootstrapSwitch("state");
    if ((new_state && !cur_state) || (!new_state && cur_state)) {
        $(this).bootstrapSwitch("toggleState");
    }
}