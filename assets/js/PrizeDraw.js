let ws = new WebSocket("ws://"+window.location.host+"/ws");
let drawing = false;
let draw_kind;
var img;
var Div;
var pre = $("#load");
let iterator = draw();



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

ws.onmessage = function (message) {
    let msg = JSON.parse(message.data);
    switch (msg.action) {
        case 'reset-page':
            $("#draw-area").html('');
            break;
        case 'start-drawing':
            console.log("start")
            start_drawing(msg);
            break;
        case 'stop-drawing':
            set_lucky_dog(msg);
            break;
        case 'send-danmu':
            eg.Send(msg.content.danmu);
            break;
        case 'set-config':
            set_config(msg);
            break;
        default:
            console.log('unknown action:\n' + message);
    }
};

function set_config(msg) {
    let theme = msg['content']['config']['theme'];
    let font = msg['content']['config']['font'];
    let music = msg['content']['config']['music'];
    let danmu = msg['content']['config']['danmu'];
}

function start_drawing(msg) {
    $("#draw-area").html('');
    draw_kind = msg['content']['dkind'];
    console.log("start_drawing\n");
    drawing = true;
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
            pre.css("display", "flex");
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

// msg: JSON Object
function set_lucky_dog(msg) {
    drawing = false;
    let user_id = msg['content']['uid'];
    let content = $("#user_" + user_id).html();
    $("#draw-area").html(content);
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

//draw_swing_init
function draw_swing() {

    console.log("swing swingswing");
    var $_this = $(this);
    var $dp = $("#swing");
    var $i = 0; //index
    var $r = 0; //round
    var $s = 150; //spead
    //快速闪动

    function dr() {
        //drawround
        if (drawing)
            $_this.res = '';
        else
            $_this.res = 1;
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

            return;
        }
        setTimeout(dr2, $s);
    }

    dr();


}

//draw_load_init
function draw_load() {

}

//draw_flash_init
function draw_flash() {
    img = [];
    Div = document.getElementsByClassName("pic");
    /*var set = setInterval(function () {
        for (var i in Div) {
            var ran = Math.ceil(Math.random() * 12);
            Div[i].src = "img/" + ran + ".jpeg";
        }
    }, 10);*/
    /*setTimeout(function () {
        clearInterval(set);
    }, 3000);*/
}

//draw_cube_init
function draw_cube() {
    img = [];
    Div = document.getElementsByClassName("pic");
    var flag = false;
    var MFB = document.getElementById("cube");
    if (!flag) {
        MFB.classList.add("cube_hover");
        flag = true;
    } else {
        MFB.classList.remove("cube_hover");
        flag = false;
    }
    /*var set = setInterval(function () {
        for (var i in Div) {
            var ran = Math.ceil(Math.random() * 12);
            Div[i].src = "img/" + ran + ".jpeg";
        }
    }, 100);
    setTimeout(function () {
        clearInterval(set);
    }, 15000);*/

}

//method of load
function* draw() {
    while (1) {
        let list = $("#user-list").children("div");
        for (let i = 0; i < list.length; i++) {
            $("#draw-area").html(list[i].innerHTML);
            yield list[i].id;
        }
    }
}


$(document).ready(function () {

    let video = document.getElementById('my-video');
    let source = document.createElement('source');
    console.log("wosdasas");
    source.setAttribute('src', 'rtmp://'+window.location.host.split(':')[0]+':1935/live/movie/123');
    source.setAttribute('type','rtmp/flv');

    video.appendChild(source);



    get_participant_list();
    setInterval(function () {
        if (drawing) {
            switch (draw_kind) {
                case "load":
                    let id = iterator.next();
                    break;
                case "swing":
                    break;
                case "flash":
                    for (var i in Div) {
                        var ran = Math.ceil(Math.random() * 12);
                        Div[i].src = "img/" + ran + ".jpeg";
                    }
                    break;
                case "cube":
                    for (var i in Div) {
                        var ran = Math.ceil(Math.random() * 12);
                        Div[i].src = "img/" + ran + ".jpeg";
                    }
                    break;
                default:
                    console.log('unknown drawing:\n' + drawing);
                    break;
            }
        }


    }, 100);
});
