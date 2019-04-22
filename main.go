package main

import (
	"github.com/c-my/lottery_client_server/config"
	"github.com/c-my/lottery_client_server/web/logger"
	"github.com/c-my/lottery_client_server/web/routers"
	"github.com/c-my/lottery_client_server/web/tools"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func main() {
	r := mux.NewRouter()
	routers.SetSubRouter(config.LocalUrl, r)

	srv := &http.Server{
		Handler: r,
		Addr:    config.LocalUrl,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	if config.LaunchBrowser {
		tools.LaunchBrowser("http://" + config.InitialUrl)
	}
	tools.Run()
	logger.Error.Fatal(srv.ListenAndServe())
}
