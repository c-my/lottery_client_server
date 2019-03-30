package main

import (
	"github.com/c-my/lottery_client_server/web/logger"
	"github.com/c-my/lottery_client_server/web/routers"
	"github.com/c-my/lottery_client_server/web/websockets"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const cloudWsServer string = "wss://sampling.alphamj.cn/ws"

func main() {
	r := mux.NewRouter()
	routers.SetSubRouter("127.0.0.1:1923", r)

	c := websockets.NewWebsocketClient(cloudWsServer)
	c.SetHandler(func(wsc *websockets.WebsocketClient, messageType int, p []byte) {
		websockets.HUB.ServerMsg <- websockets.ServerMsg{messageType, p}
	})
	c.Run()

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:1923",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Error.Fatal(srv.ListenAndServe())
}

//func addAction(action string, content interface{}) ([]byte, error) {
//	m := map[string]interface{}{
//		"action":  action,
//		"content": content,
//	}
//	return json.Marshal(m)
//}
