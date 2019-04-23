let next_page = "start-menu";
let cloud_url = "https://sampling.alphamj.cn/signin";
let local_url = window.location.origin + "/post";

function inform_local(isOnline, uid) {
    $.ajax({
        type: "POST",
        url: local_url,
        dataType: "json",
        data: JSON.stringify({
            "action": "login",
            "content": {
                online: isOnline
            }
        })
    });
}

$(function() {
    history.pushState(null, null, document.URL);
    window.addEventListener('popstate', function () {
        history.pushState(null, null, document.URL);
    });
    var storage = window.localStorage;
    var username = $("#username"),
        password = $('#password'),
        check = $('#remember-me');
    var localUser = storage.getItem('user') || '';
        localPass = storage.getItem('pass') || '';
    if (localUser !== '' && localPass !== '') {
        username.val(localUser);
        password.val(localPass);
        check.attr('checked', 'true');
    }
    $("#login, #offline").click(function remember() {
        if (check.get(0).checked) {
            storage.setItem('user', username.val());
            storage.setItem('pass', password.val());
        } else {
            storage.setItem('user', '');
            storage.setItem('pass', '');
        }
    });
    $("#offline").click(function() {
        inform_local("false");
        self.location.href = next_page;
    });
    $("#login").click(function() {
        $.ajax({
            type: "POST",
            url: "https://sampling.alphamj.cn/signin",
            data: JSON.stringify({
                username: $("#username").val(),
                password: $("#password").val()
            })
        });
        $.ajax({
            method: "POST",
            url: window.location.origin + '/signin',
            timeout: 10000,
            dataType : "json",
            data: JSON.stringify({
                username: $("#username").val(),
                password: $("#password").val()
            }),
            success: function(data) {
                console.log(data);
                if (data.success === "true") {
                    console.log("uid: ", data.uid);
                    inform_local("true");
                    self.location.href = next_page;
                } else {
                    $(".log-bar").text("用户名或密码错误");
                    console.log("message: ", data.success);
                }
            },
            complete: function(XMLHttpRequest, status) {
                if (status == 'timeout') {
                    $(".log-bar").text("网络超时错误");
                }
        　　},
            error: function(xhr, state, errorThrown){
                if (!xhr.status)
                $(".log-bar").text("未连接网络");
            }
        });
        return false;
    });
});