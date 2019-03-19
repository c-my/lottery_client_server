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
                                <img class="circle img" src="${user_avatar}" alt="avatar">
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

let drawing = true;

ws.onmessage = function(message) {
    let msg = JSON.parse(message);
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
    let user_id = msg['uid'];
    let content = $("#user_" + user_id).html();
    $("#draw-area").html(content);
}

function* draw() {
    let list = $("#user-list").children();
    for (let item in list) {
        $("#draw-area").html($(list[item]).html());
        yield list[item].id;
    }
}

let iterator = draw();

$(document).ready(function () {
    get_participant_list();
    setInterval(function () {
        if (drawing) {
            iterator.next();
        }
    }, 300);
});