let ws = new WebSocket("ws://localhost:1923/recv");

ws.onmessage = function (message) {
    $("#msg-recv-content").append(message.data + "\r\n");
};

function echo(statebar, state, content) {
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

$("#msg-send-send").click(function() {
    $("#msg-send-form").ajaxSubmit({
        method: "POST",
        url: "http://localhost:1923/send",
        timeout : 1000,
        data: {
            'content': $("#msg-send-content").val()
        },
        success: function(text) {
            var data = JSON.parse(text);
            if (data.success === "true") {
                echo("#send-statebar", 0, "消息已发送");
            }
        },
        complete: function(XMLHttpRequest, status) {
            if (status == 'timeout') {
                echo("#send-statebar", 2, "网络超时");
            }
    　　}
    });
    return false;
});
$("#msg-recv-reset").click(function(event) {
    $("#msg-recv-content").empty();
});