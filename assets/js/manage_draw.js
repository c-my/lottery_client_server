var forbidden_words_list=new Array("优秀","太强了","天下第一");
var forbidden_times=new Array(0,0,0);
var black_users_list=new Array("李金牙","李金嘴","黄一凡");
var user_map = new Object();
let $cur_line = new Object();
var bullet_num=0;
let cursor_used = false;
update_forbidden_times();
let ws = new WebSocket('ws://127.0.0.1:1923/ws');
ws.onmessage = function (message) {
    console.log(message.data);
    let msg = JSON.parse(message.data);

    switch (msg.action) {
        case "bullet":
            var forbidden = false;
            for(var i=0;i<forbidden_words_list.length;i++) {
                var tem=msg.content.bullet.indexOf(forbidden_words_list[i]);
                if(tem>-1 || distinctive_check(black_users_list,msg.content.uid)==false){
                    console.log("bullet forbidden");
                    forbidden_times[i]++;
                    update_forbidden_times();
                    forbidden = true; break;
                }
            }
            if (!forbidden) {
                user_map[msg.uid] = msg.user_nickname;
                add_bullet(
                    msg.content.time,
                    msg.content.uid,
                    msg.content.nickname,
                    msg.content.bullet,
                    msg.content.id
                );
            }
            break;
        case "add-forbidden-time":
            for(var i=0;i<forbidden_words_list.length;i++) {
                if (forbidden_words_list[i] == msg.content) {
                    forbidden_times[i]++;
                    update_forbidden_times();
                    break;
                }
            }
            break;
    }
}
/// 小工具
String.prototype.format = function() {
    var args = Array.prototype.slice.call(arguments);
    var count = 0;
    return this.replace(/%s/g, function(s, i){
        return args[count ++];
    });
}

/// 快捷键
function pressDown() {
    if ($("#bullets_list_tbody tr").length == 0) return;
    $cur_line.removeClass("active");
    if ($cur_line.nextAll().length == 0) {
        $cur_line = $("#bullets_list_tbody tr:first");
    } else {
        $cur_line = $cur_line.next();
    }
    $cur_line.addClass("active");    
}
document.onkeydown = function(event) {
    switch (event.keyCode) {
        case 38: { // UP
            cursor_used = false;
            if ($("#bullets_list_tbody tr").length == 0) return;
            $cur_line.removeClass("active");
            if ($cur_line.prevAll().length == 0) {
                $cur_line = $("#bullets_list_tbody tr:last");
            } else {
                $cur_line = $cur_line.prev();
            }
            $cur_line.addClass("active");
            $("#bullets_list").animate(
                {scrollTop:$cur_line.prevAll().length * $cur_line.height()}, 100);
            break;
        }
        case 40: { // Down
            cursor_used = false;
            pressDown();
            $("#bullets_list").animate(
                {scrollTop:$cur_line.prevAll().length * $cur_line.height()}, 100);
            break;
        }
        case 37: { // Left
            allow_this_draw();
            break;
        }
        case 39: { // Right
            forbidden_this_draw();
            break;
        }
    }
};

/// 点击选中弹幕
function select_bullet(obj) {
    cursor_used = true;
    $cur_line.removeClass("active");
    $cur_line = $(obj); 
    $cur_line.addClass("active");
}

///添加屏蔽词
function add_forbidden() {
    var str=$('#add_a_forbidden_word').val();
    if(distinctive_check(forbidden_words_list,str)==false){
        toastr.error('敏感词汇:'+str+'已经被屏蔽！');
        return;
    }
    if(str=="")
        return;
    console.log(str);
    forbidden_words_list.push(str);
    forbidden_times.push(0);
    var tr=document.createElement('tr');
    var th=document.createElement('th');
    var td1=document.createElement('td');
    var td2=document.createElement('td');
    th.scope="row";
    th.innerText=forbidden_words_list.length;
    td1.innerText=str;
    td2.innerHTML='0'+'<a onclick="delete_this_word(this)" style="float:right" href="#!">删除</a>';
    tr.appendChild(th);
    tr.appendChild(td1);
    tr.appendChild(td2);
    $('#forbidden_words_body').append(tr);
    document.getElementById('add_a_forbidden_word').value="";
    ws.send(JSON.stringify({
        action: "add-forbidden-word",
        content: str
    }));
    toastr.warning('已添加屏蔽词'+str);
}
// 添加黑名单
function add_black_user(str) {
    if(distinctive_check(black_users_list,str)==false){
        toastr.error('用户:'+str+'已经在黑名单中！');
        return;
    }
    if(str=="")
        return;
    console.log(str);
    black_users_list.push(str);
    var tr = "<tr>%s</tr>";
    var td = "<td>%s</td>";
    var th = "<th scope='row'>%s</th>";
    var line_number = th.format(black_users_list.length);
    var user_id = td.format(str);
    var nickname = (str in user_map) ? user_map[str]: "--";
    var user_nickname = td.format(nickname + `
    <a onclick="delete_this_people(this)" style="float:right" href="#!">删除</a>`);
    $('#black_users_list_body').append(
        tr.format(line_number + user_id + user_nickname));
    document.getElementById('add_a_black_user').value = "";
    ws.send(JSON.stringify({
        action: "add-forbidden-user",
        content: str
    }));
    toastr.warning('已将用户'+str+'添加到黑名单');
}
var tem=0;

// 添加弹幕
function add_bullet(time, user_id, nickname, bullet, id) {
    tem ++;
    var tr = "<tr bullet_id='%s' onclick='select_bullet(this)'>%s</tr>";
    var td = "<td class='%s'>%s</td>";
    var item_time = td.format("bullet-time", time);
    var item_user_id = td.format("bullet-userid", user_id);
    var item_nickname = td.format("bullet-nickname", nickname);
    var item_bullet = td.format("bullet-content", bullet);
    $('#bullets_list_tbody').append(
        tr.format(id, item_time + item_user_id + item_nickname + item_bullet));
    if ($("#bullets_list_tbody tr").length == 1) {
        $cur_line = $("#bullets_list_tbody tr:first");
        $cur_line.addClass("active");
    }
    bullet_num ++;
}
for(var i=1;i<=100;i++){
    if(i&1)
        add_bullet("13213" + String(i), "314159", "frog", "不要总想着搞什么大新闻", i);
    else
        add_bullet("13213" + String(i), "265358", "frog", "安慕安格瑞", i);
}

// 屏蔽该条弹幕
function forbidden_this_draw() {
    if ($("#bullets_list_tbody tr").length == 0) return;
    toastr.clear();
    var str = $cur_line.find(".bullet-content").text();
    $temp = $cur_line;
    pressDown();
    $temp.remove();
    if (!cursor_used) $("#bullets_list").animate({
        scrollTop:$cur_line.prevAll().length * $cur_line.height()}, 0);
    toastr.warning('已屏蔽弹幕: “' + str + '”');
}

// 允许该条弹幕
function allow_this_draw() {
    if ($("#bullets_list_tbody tr").length == 0) return;
    toastr.clear();
    var id = $cur_line.attr("bullet_id");
    var str = $cur_line.find(".bullet-content").text();
    $temp = $cur_line;
    pressDown();
    $temp.remove();
    if (!cursor_used) $("#bullets_list").animate({
        scrollTop:$cur_line.prevAll().length * $cur_line.height()}, 0);
    toastr.success('已允许弹幕: “' + str + '”');
    ws.send(JSON.stringify({
        action: 'bullet_allow',
        content: id
    }));
}

// 屏蔽该用户
function forbidden_this_user() {
    if ($("#bullets_list_tbody tr").length == 0) return;
    toastr.clear();
    var str = $cur_line.find(".bullet-userid").text();
    add_black_user(str);
    var first = true;
    $cur_line.removeClass("active");
    $("#bullets_list_tbody tr").each(function() {
        if ($(this).find(".bullet-userid").text() == str) {
            $(this).remove();
        } else if (first) {
            $cur_line = $(this);
            $cur_line.addClass("active");
            first = false;
        }
    });
}

// 在屏蔽词表中删除词
function delete_this_word(e) {
    var str=e.parentNode.previousSibling.innerText;
    console.log('str'+str);
    var loc=removearr(forbidden_words_list,str)+1;
    console.log(loc);
    var tem=e.parentNode.parentNode;
    while($(tem).next().length!=0){
        tem=tem.nextSibling;
        tem.firstElementChild.innerText=loc;
        loc++;
    }
    for(var i=loc;i<forbidden_times.length;i++){
        forbidden_times[i-1]=forbidden_times[i];
    }
    forbidden_times.pop();
    $(e).parents("tr").remove();
    ws.send(JSON.stringify({
        action: "remove-forbidden-word",
        content: str
    }));
}
// 在黑名单中删除用户
function delete_this_people(e) {
    var str=e.parentNode.value;
    console.log('str'+str);
    var loc=removearr(black_users_list,str)+1;
    console.log(loc);
    var tem=e.parentNode.parentNode;
    while($(tem).next().length!=0){
        tem=tem.nextSibling;
        tem.firstElementChild.innerText=loc;
        loc++;
    }
    $(e).parents("tr").remove();
    ws.send(JSON.stringify({
        action: "remove-forbidden-user",
        content: str
    }));
}

function removearr(a,b) {
    var loc=0;
    var flag=0;
    for(var i=0;i<a.length;i++){
        if(flag==1){
            a[i-1]=a[i];
        }else{
            console.log(a[i]+','+b);
            if(a[i]==b){
                flag=1;
                loc=i;
            }
        }
    }
    a.pop();
    return loc;
}

// 重复检查
function distinctive_check(a,b) {
    for(var i=0;i<a.length;i++){
        if(a[i]==b)
            return false;
    }
    return true;
}

// 更新屏蔽次数
function update_forbidden_times() {
    var tem=document.getElementById('forbidden_words_body').firstElementChild;
    for(var i=0;i<forbidden_times.length;i++){
        tem.lastElementChild.innerText=forbidden_times[i].toString();
        tem=tem.nextElementSibling;
    }
}