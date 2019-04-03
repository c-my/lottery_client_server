let ws = new WebSocket("ws://127.0.0.1:1923/ws");


function get_participant_list() {

    $.ajax({
        url: 'get-exist-user',
        dataType: 'json',
        method: 'GET',
        success: function (result) {
            for (let i in result) {
                let user_id = result[i]['uid'];
                let user_avatar = result[i]['avatar'];
                let user_name = result[i]['nickname'];
                let user_info_template =
                    `<div id="user_${user_id}">
                        <table>
                            <tr><td align="center">
                                <img class="circle img avatar-img" src="${user_avatar}" alt="">
                            </td></tr>
                            <tr><td>
                                <span>${user_name}</span>
                            </td></tr>
                        </table>
                    </div>`;
                $("#user-list").append(user_info_template);
            }
            // $("#draw-area").html($("#user_20160001").html())
        }
    })
}

let drawing = false;
let draw_kind;
ws.onmessage = function (message) {
    let msg = JSON.parse(message.data);
    switch (msg.action) {
        case 'reset-page':
            $("#draw-area").html('');
            break;
        case 'start-drawing':
            start_drawing(msg);
            break;
        case 'who-is-lucky-dog':
            if(draw_kind ==="load"){
                set_lucky_dog(msg);
            }
            break;
        case 'send-danmu':
            eg.Send(msg.content.danmu);
            break;
        default:
            console.log('unknown action:\n' + message);
    }
};
var pre = $("#load");

function start_drawing(msg) {
    $("#draw-area").html('');
    draw_kind = msg['content']['dkind'];
    switch (draw_kind) {
        case "load"://load
            pre.hide();
            pre = $("#load");
            pre.show();
            console.log("load\n");
            draw_load();
            break;
        case "swing"://swing
            pre.hide();
            pre = $("#swing");
            pre.css("display","flex");
            console.log("swing\n");
            init_swing();
            draw_swing();
            break;
        case "flash"://flash
            pre.hide();
            pre = $("#flash");
            pre.show();
            console.log("flash\n");
            draw_flash();
            break;
        case "cube"://cube
            pre.hide();
            pre = $("#cube_div");
            pre.show();
            console.log("cube\n");
            draw_cube();
            break;
        default:
            break;
    }
}

function init_swing() {
    //获取图片列表，之后后台在这给数据就行
    function pictureList() {
        var pictureList = [];
        for (var i = 1; i <= 12; i++) {
            pictureList.push("img/" + i + ".jpeg");//把得到的图片放到数组中
        }
        return pictureList;//返回图片数组
    };

    var list = pictureList();//图片数组
    var divList = [];//div数组
    for (var i in list) {
        //添加图片,
        var newDiv = document.createElement("div");
        var newImg = document.createElement("img");
        //var newSDiv=document.createElement("div");
        var user = document.getElementById("swing");
        newDiv.className = "userW";
        newImg.className = "pro";
        // newSDiv.className="userS";
        newDiv.appendChild(newImg);
        // newDiv.appendChild(newSDiv);
        newImg.src = list[i];
        user.appendChild(newDiv);
        divList.push(newImg);//获取div数组

    }
    console.log(divList);
    //var flag=false;
}

function pre_change(element) {
    pre.style.display = "none";
    pre = element;
    pre.style.display = "inline";
}

// msg: JSON Object
function set_lucky_dog(msg) {
    drawing = false;
    let user_id = msg['content']['uid'];
    let content = $("#user_" + user_id).html();
    $("#draw-area").html(content);
}

function draw_swing() {
    $.fn.draw = function (options) {
        var $_this = $(this);
        var $dp = options.prod;
        var $i = 0; //index
        var $r = 0; //round
        var $s = 150; //spead
        //快速闪动，会走两遍
        function dr() {
            //drawround
            if ($_this.res != '' && $_this.res != undefined && $i == 0) {
                dr2();
            } else {
                $dp.find(".pro").removeClass("select").eq($i).addClass("select");
                $i = $i >= 11 ? 0 : $i + 1;//用户数组的数量
                setTimeout(dr, $s);
            }
        }

        function dr2() {
            $dp.find(".pro").removeClass("select").eq($i).addClass("select");
            $i = $i >= 11 ? 0 : $i + 1;
            $s = $s + 30;//让他变慢
            $rand = Math.ceil(Math.random() * 12)
            if ($r < $_this.res + $rand) {
                $r++;
            } else {
                $i = 0;
                $r = 0;
                $s = 100;
                setTimeout(result, 1000);
                return;
            }
            setTimeout(dr2, $s);
        }

        function getRes() {
            $_this.res = '';
            //赋结果值
            setTimeout(function () {
                $_this.res = 1
            }, 2000);
        }

        //监听点击事件
        function click() {
            $_this.bind("click", function () {
                getRes();
                console.log(i);
                dr();
                $_this.unbind("click");
            });
        }

        click();

        //输出获奖者
        function result() {
            alert("XX得到奖项");
            click();
        }
    }
    $("#btn").draw({
        prod: $("#swing")
    });
}

function draw_load() {
    drawing = true;
}

function draw_flash() {
    var img = [];
    var Div = document.getElementsByClassName("pic");
    var set = setInterval(function () {
        for (var i in Div) {
            var ran = Math.ceil(Math.random() * 12);
            Div[i].src = "img/" + ran + ".jpeg";
        }
    }, 10);
    setTimeout(function () {
        clearInterval(set);
    }, 3000);
}

function draw_cube() {
    var img = [];
    var Div = document.getElementsByClassName("pic");
    var flag = false;
    var MFB = document.getElementById("cube");
    if (!flag) {
        MFB.classList.add("cube_hover");
        flag = true;
    } else {
        MFB.classList.remove("cube_hover");
        flag = false;
    }
    var set = setInterval(function () {
        for (var i in Div) {
            var ran = Math.ceil(Math.random() * 12);
            Div[i].src = "img/" + ran + ".jpeg";
        }
    }, 100);
    setTimeout(function () {
        clearInterval(set);
    }, 15000);

}

function* draw() {
    while (1) {
        let list = $("#user-list").children("div");
        for (let i = 0; i < list.length; i++) {
            $("#draw-area").html(list[i].innerHTML);
            yield list[i].id;
        }
    }
}

let iterator = draw();

$(document).ready(function () {

    get_participant_list();
    setInterval(function () {
        if (drawing) {
            let id = iterator.next();
        }
    }, 100);
});