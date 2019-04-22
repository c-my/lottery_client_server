function getJsonLength(jsonData){
    var jsonLength = 0;
    for(var item in jsonData){
        jsonLength++;
    }
    return jsonLength;
}

function getQueryString(name) {//获取name参数的值
    var result = window.location.search.match(new RegExp("[\?\&]" + name + "=([^\&]+)", "i"));
    if (result == null || result.length < 1) {
        return "";
    }
    return result[1];
}



if(getQueryString('set-from-history')!="")
    ws.send(JSON.stringify({
        action: "set-from-history",
        content: getQueryString('set-from-history'),
    }));

// var fso=new ActiveXObject(Scripting.FileSystemObject);
// var f=fso.opentextfile(path,1,true);
// while (!f.AtEndOfStream) {
//     console.log(f.Readline());
// }
// f.close();


// jsReadFiles(path);
