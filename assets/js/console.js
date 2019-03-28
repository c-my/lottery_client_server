let drawing = false;
let running = false;

function start_draw() {
    $("#draw-label").text("停止抽奖");
    $(".btn-draw-action").css("background-color", "orangered");
    $("#draw-action").attr("class", "mdi mdi-stop-circle-outline right");
}

function stop_draw() {
    $("#draw-label").text("开始抽奖");
    $(".btn-draw-action").css("background-color", "forestgreen");
    $("#draw-action").attr("class", "mdi mdi-arrow-right-drop-circle-outline right");
}

function on_draw_btn_click() {
    if (!drawing) start_draw(); else stop_draw();
    drawing = !drawing;
}

function on_activity_btn_click() {
    if (!running) {
        $("#activity-label").text("结束活动");
        $(".btn-activity-action").css("background-color", "orangered");
        $("#activity-action").attr("class", "mdi mdi-stop-circle-outline right");
        $(".btn-draw-action").removeClass("disabled");
    } else {
        $("#activity-label").text("活动结束");
        $(".btn-activity-action").css("background-color", "forestgreen");
        $("#activity-action").attr("class", "mdi mdi-arrow-right-drop-circle-outline right");
        $(".btn-activity-action").addClass("disabled");
        $(".btn-draw-action").addClass("disabled");
        stop_draw();
    }
    running = true;
}
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
            var t = record[i].split(",");
            for (var y = 0; y < t.length; y++) {
                if (!data[i]) data[i] = {};
                data[i][title[y]] = t[y];
            }
        }
        f.call(this, data);
        data = null;
    }
});
function import_participants() {
    var file = document.getElementById("participants").files[0];
    var fileName = file.name.substring(0, file.name.lastIndexOf('.'));
    var reader = new FileReader();
    reader.onload = function(f) {
        $('#parts').val($('#participants').val());
        $('.parts-label').addClass("active");
        var data = this.result;
        $(function() {
            $.csv(data, function(data) {
                if (!data) return;
                if (data && !data[0].hasOwnProperty('username')) {
                    for (var item in data[0]) console.log(item);
                    return;
                }
                var head_str = "<tr>";
                for (item in data[0])
                    head_str += "<th>" + item + "</th>";
                $("#parts-list-head").append(head_str + "</tr>");
                for (var i = 0; i < data.length; ++i) {
                    var record_str = "<tr>";
                    for (item in data[0])
                        record_str += "<td>" + data[i][item] + "</td>"
                    $("#parts-list-body").append(record_str + "</tr>");
                }
            });
        });
    };
    reader.readAsText(file);
}
var eCalendar = {
    type  : "time",   //模式，time: 带时间选择; date: 不带时间选择;
    stamp : false,    //是否转成时间戳，默认true;
    offset: [0, 2],   //弹框手动偏移量;
    format: "yyyy年mm月dd日 hh:ii",   //时间格式 默认 yyyy-mm-dd hh:ii;
    skin  : "#1B7690",   //皮肤颜色，默认随机，可选值：0-8,或者直接标注颜色值;
    step  : 10,   //选择时间分钟的精确度;
    callback:function(v,e) {} //回调函数
};
function del_row(obj) {
    var tr = obj.parentNode.parentNode;
    tr.parentNode.removeChild(tr);
}
function drag_prize(event) {
    event.dataTransfer.setData("Text", event.target.value);
}
function drop_prize(event) {
    event.preventDefault();
    var data = event.dataTransfer.getData("Text");
    event.target.value = data;
}
function allowDrop(event) {
    event.preventDefault();
}
function check_content(obj) {
    if (obj.value == "") {
        del_row(obj);
        return false;
    }
    var tr = $(obj).parent().parent();
    var obj_index = tr.prevAll().length;
    $(tr.parent()).find("tr").each(function(index, element) {
        var input = $(this).children("td:first").find("input").val();
        if (index != obj_index && obj.value == input) {
            Materialize.toast('名称重复', 1000);
            obj.focus();
            return false;
        }
    });
    return true;
}
function show_lucky_dog(username, prizename) {
    $("#lucky-list-body").append("<tr><td>" + username + "</td><td>" + prizename + "</td></tr>");
}
$(document).ready(function () {
    $(".btn-draw-action").click(on_draw_btn_click);
    $(".btn-activity-action").click(on_activity_btn_click);
    $("#participants").change(import_participants);
    $("#start-time").ECalendar(eCalendar);
    $("#end-time").ECalendar(eCalendar);
    $('select').material_select();
    $('INPUT.minicolors').minicolors();
    $(".btn-add-prize").click(function() {
        if ($('#prize-list-body').find("tr").length == 0 || $('#prize-list-body tr:last').find('input').val()) {
            $("#prize-list-body").append(`<tr>
            <td><input class="table-input" style="margin-bottom: 0px;" type="text" draggable="true" ondragstart="drag_prize(event)" onblur="check_content(this)"></input></td>
            <td><div class="del-prize" onclick="del_row(this)"><i class="mdi mdi-close-circle"></i></div></td>
            </tr>`);
        }
        $('#prize-list-body tr:last').find('input').focus();
    });
    $(".btn-add-item").click(function() {
        if ($('#item-list-body').find("tr").length == 0 || $('#item-list-body tr:last').find('input.item-name').val()) {
            $("#item-list-body").append(`<tr>
            <td><input class="table-input item-name" style="margin-bottom: 0px;" type="text" onblur="check_content(this)"></input></td>
            <td><input readonly="readonly" style="margin-bottom: 0px; color: rgb(64, 77, 91);" ondrop="drop_prize(event)" ondragover="allowDrop(event)"></input></td>
            <td><input class="table-input" style="margin-bottom: 0px;" type="number" min="1"></input></td>
            <td><div class="del-item" onclick="del_row(this)"><i class="mdi mdi-close-circle"></i></div></td>
            </tr>`);
        }
        $('#item-list-body tr:last').find('input.item-name').focus();
    });
    $("#cur-item-input").on('focus', function() {
        $(this).val("");
        var items = "";
        $("#item-list-body tr").each(function() {
            var text = $(this).children("td:first").find("input").val();
            if (text) {
                items += '<option value="' + text + '" />';
            }
        });
        $("#items").html(items);
        return false;
    });
});