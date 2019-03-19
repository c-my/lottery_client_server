let ws = new Ws("ws://127.0.0.1:8000/ws");

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

ws.On("reset-page", function(msg) {
    $("#draw-area").html('');
});

let drawing = false;
ws.On("start-drawing", function (msg) {
    drawing = true;
});

ws.On("who-is-lucky-dog", function (msg) {
    drawing = false;
    let result = JSON.parse(msg);
    let user_id = result['uid'];
    let content = $("#user_" + user_id).html();
    $("#draw-area").html(content);
});

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