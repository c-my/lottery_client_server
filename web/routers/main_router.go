package routers

import (
	"encoding/json"
	"github.com/c-my/lottery_client_server/web/controllers"
	"github.com/c-my/lottery_client_server/web/logger"
	"github.com/gorilla/mux"
	"net/http"
)

// SetSubRouter sets sub router for root
func SetSubRouter(parent string, r *mux.Router) {
	subRouter := r.Host(parent).Subrouter()
	setStatic(subRouter)
	setGet(subRouter)
	setPost(subRouter)
}

func setStatic(r *mux.Router) {
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("assets/css"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("assets/js"))))
	r.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/", http.FileServer(http.Dir("assets/fonts"))))
	r.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("assets/img"))))
	r.PathPrefix("/avatars/").Handler(http.StripPrefix("/avatars/", http.FileServer(http.Dir("assets/avatars"))))
}

func setGet(r *mux.Router) {
	r.HandleFunc("/get-exist-user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(controllers.UserControl.Get())
		if err != nil {
			logger.Error.Println("failed to get exist user: ", err)
		}
	})

	r.HandleFunc("/console", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "console.html")
	})

	r.HandleFunc("/screen", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "PrizeDraw.html")
	})

	r.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "login.html")
	}).Methods("GET")
}

func setPost(r *mux.Router) {
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		js, err := json.Marshal(map[string]string{
			"success": "true",
		})
		if err != nil {
			logger.Error.Println("an impossible error happened: ", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(js)
		if err != nil {
			logger.Error.Println("fail to response login: ", err)
		}
	}).Methods("POST")
}
