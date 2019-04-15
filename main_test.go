package main

import (
	"fmt"
	"github.com/c-my/lottery_client_server/web/websockets"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	clientWs, err := websockets.NewWebsocketClient("ws://127.0.0.1:1923/ws")
	if err != nil {
		fmt.Println(err)
	}
	clientWs.SetHandler(nil)
	clientWs.Run()
	loginJson := `{"username":"666666", "password":"666666"}`
	actJson := `{
    "name": "东北大学才明洋表彰大会"
    [,"start_time": "%Y-%m-%d %H:%M:%S"]
    [,"end_time": "结束时间"]
}`
	http.Post("http://127.0.0.1:1923/signin", "application/json", strings.NewReader(loginJson))
	http.PostForm("http://127.0.0.1:1923/signin", url.Values{
		"username": {"666666"},
		"password": {"666666"},
	})

	http.Post("http://127.0.0.1:1923/append-activity", "application/json", strings.NewReader(actJson))

	m.Run()
}

func TestGetAct(t *testing.T) {
	res, _ := http.Get("http://127.0.0.1:1923/get-activities")
	byt, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(byt))
}
