$(function() {
    var storage = window.localStorage;
    var username = $("#username"),
        password = $('#password'),
        check = $('#remember-me');
    var localUser = storage.getItem('user') || '';
        localPass = storage.getItem('pass') || '';
    if (localUser !== '' && localPass !== '') {
        $(".user-icon, .user-label, .pwd-icon, .pwd-label").addClass("active");
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
    $("#login").click(function() {
        self.location.href = "console";
        // $("#login-form").ajaxSubmit({
        //     timeout : 1000,
        //     data: {
        //         'username': $("#username").val(),
        //         'password': $("#password").val()
        //     },
        //     success: function(data) {
        //         if (data.success === "true") {
        //             self.location.href = "console.html";
        //         } else {
        //             $(".log-bar").text("用户名或密码错误");
        //         }
        //     },
        //     complete: function(XMLHttpRequest, status) {
        //         if (status == 'timeout') {
        //             $(".log-bar").text("网络超时错误");
        //         }
        // 　　}
        // });
        // return false;
    });
});