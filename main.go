package main

import (
	"github.com/c-my/lottery_client_server/config"
	"github.com/c-my/lottery_client_server/web/logger"
	"github.com/c-my/lottery_client_server/web/routers"
	"github.com/c-my/lottery_client_server/web/tools"
	"github.com/c-my/lottery_client_server/web/websockets"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func main() {
	r := mux.NewRouter()
	routers.SetSubRouter(config.LocalUrl, r)

	c, err := websockets.NewWebsocketClient(config.CloudWsServer)
	if err != nil {
		logger.Warning.Println("stop trying to connect")
	} else {
		c.SetHandler(func(wsc *websockets.WebsocketClient, messageType int, p []byte) {
			websockets.HUB.ServerMsg <- websockets.ServerMsg{messageType, p}
		})
		c.Run()
	}

	srv := &http.Server{
		Handler: r,
		Addr:    config.LocalUrl,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	tools.LaunchBrowser("http://" + config.InitialUrl)

	logger.Error.Fatal(srv.ListenAndServe())
}

//func addAction(action string, content interface{}) ([]byte, error) {
//	m := map[string]interface{}{
//		"action":  action,
//		"content": content,
//	}
//	return json.Marshal(m)
//}
