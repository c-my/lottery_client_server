## WsTool: 前端接口调试助手
***
### 版本信息
+ Version： 0.1 内测版
+ Author： 东北大学-16级物联网工程-刘尚育
+ Description：Postman用来调试后端接口，而wsTool用来调试前端

### 支持功能
+ 与网页建立websocket连接
	+ 接收和显示JSON格式的消息
	+ 编辑和发送JSON格式的消息
+ 接收和显示网页发出的POST请求
	+ 返回JSON：`{"success": "true"}`

### 使用方法
下载本级文件夹中的所有内容，打开wsTool.exe服务器。打开要调试的前端页面和接口测试工具页面`wsTool.html`。注意要保证两个页面通信之前wsTool.exe为打开状态，否则可以尝试刷新界面重新建立连接。
在接口测试工具`wsTool.html`中，左侧窗口为想要通过websocket推送到待测前端的数据编辑区，写入JSON格式的文本后，点击右上角的发送图标发送。发送成功将在左下角显示发送成功，网络错误将在左下角显示网络超时（由于服务器没有启动），暂时不提供JSON的语法检查。
在接口测试工具`wsTool.html`中，右侧窗口为从待测前端接收的websocket数据或post的数据。可以点击右上角的重置按钮清空面板。
#### websocket API
在待测试的前端的javascript中与工具建立websocket连接的格式为：
```js
let ws = new WebSocket('ws://127.0.0.1:1923/ws');
```
在待测试的前端的javascript中监听消息的格式为：
```js
ws.onmessage = function (message) {
    let msg = JSON.parse(message.data); // msg为解析得到的对象
    switch (msg.action) {
        case 'Your-Action-1':
            // do something
            break;
        case 'Your-Action-2':
            // do something
			break;
        default:
			// otherwise
            console.log('unknown action:\n' + msg.action);
    }
};
```
在待测试的前端的javascript中发送消息的格式为：
```js
ws.send(JSON.stringify({
	action: 'Your-action',
	content: { // JSON格式的内容
		'name_1': 'value_1',
		'name_2': 'value_2'
	}
}));
```
#### http.post API
你可以使用任何方法实现POST，但注意`url: "http://localhost:1923/post"`，内容需为JSON格式，且包含字段格式形如`{ "action": "",content: { } }`。下面是使用`jQuery`的`$.ajax()`方法实现的例子：
```js
$.ajax({
	type: "POST",
	url: "http://localhost:1923/post",
	data: JSON.stringify({
		action: "background-image",
		dataType: "json",
		content: {
            "base64Str": base64Str
        }
	}),
	success: function(result) {
		console.log("send image: " + result);
	}
});
```