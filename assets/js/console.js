let ws = new WebSocket('ws://127.0.0.1:8000/ws');

ws.onmessage = function (message) {
    console.log(message.data);
    let msg = JSON.parse(message.data);
    switch (msg.action) {
        case 'who-is-lucky-dog':
            add_luck_dog(msg);
            break;
        default:
            console.log('unknown action:\n' + message);
    }
};

let running = false;

function on_pick_button_click() {
    if (!running) {
        $("#pick-btn").html("停！");
        ws.send(JSON.stringify({
            action: 'start-drawing'
        }));
    } else {
        $("#pick-btn").html("抽！");
        ws.send(JSON.stringify({
            action: 'stop-drawing'
        }));
    }
    running = !running;
}

function on_reset_button_click() {
    ws.send(JSON.stringify({
        action: 'reset-page'
    }));
}


// msg: JSON Object
function add_luck_dog(msg) {
    let user_id = msg['content']['uid'];
    let content = $("#user_" + user_id).html();
    let user_info_template = `<div id="won_${user_id}" class="valign-wrapper">${content}</div>`;
    $("#won-list").append('<li class="collection-item">' + user_info_template + '</li>')
}

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
                let user_info_template = `<div id="user_${user_id}" class="valign-wrapper"><img class="circle img" src="${user_avatar}" alt="avatar"><span>${user_name}</span></div>`;
                $("#user-list").append('<li class="collection-item">' + user_info_template + '</li>')
            }
        }
    })
}

$(document).ready(function () {
    $('select').material_select();
    get_participant_list();
    $("#pick-btn").click(on_pick_button_click);
    $("#reset-btn").click(on_reset_button_click);
});