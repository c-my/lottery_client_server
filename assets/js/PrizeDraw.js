let ws = new WebSocket("ws://127.0.0.1:8000/ws");

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
                            <tr><td>
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

ws.onmessage = function(message) {
    let msg = JSON.parse(message.data);
    switch (msg.action) {
        case 'reset-page':
            $("#draw-area").html('');
            break;
        case 'start-drawing':
            drawing = true;
            break;
        case 'who-is-lucky-dog':
            set_lucky_dog(msg);
            break;
        default:
            console.log('unknown action:\n' + message);
    }
};


// msg: JSON Object
function set_lucky_dog (msg) {
    drawing = false;
    let user_id = msg['content']['uid'];
    let content = $("#user_" + user_id).html();
    $("#draw-area").html(content);
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