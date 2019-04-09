package routers

import (
	"bytes"
	"encoding/json"
	"github.com/c-my/lottery_client_server/config"
	"github.com/c-my/lottery_client_server/web/controllers"
	"github.com/c-my/lottery_client_server/web/logger"
	"github.com/c-my/lottery_client_server/web/websockets"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
)

// SetSubRouter sets sub router for root
func SetSubRouter(parent string, r *mux.Router) {
	subRouter := r.Host(parent).Subrouter()
	setStatic(subRouter)
	setGet(subRouter)
	setPost(subRouter)
	setWebsocket(subRouter)
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
			logger.Error.Println("failed to get exist user:", err)
		}
	})

	r.HandleFunc("/get-activities", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write(controllers.ActivityControl.GetAll())
		if err != nil {
			logger.Error.Println("failed to get exist user:", err)
		}
	})

	r.HandleFunc("/get-participants", func(writer http.ResponseWriter, request *http.Request) {
		//queryActivity := request.Form["activity"]
		// get user list
		writer.Header().Set("Content-Type", "application/json")
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
	r.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {

		localJson := parsePostForm(r)
		logger.Info.Println("post to cloud:", string(localJson))

		loginMap, cookies := parsePostResponse(config.CloudLoginURL, localJson)

		var err error
		var result = make(map[string]string)

		switch loginMap["result"] {
		case "success":
			result["success"] = "true"

			session := "sessionid=" + getSessionFromCookie(cookies)
			logger.Info.Println("trying connecting to cloud with session:", session)
			websockets.Client, err = websockets.NewWebsocketClient(config.CloudWsServer, session)
			if err != nil {
				logger.Warning.Println("failed to connect to cloud", err)
			} else {
				websockets.Client.SetHandler(func(wsc *websockets.WebsocketClient, messageType int, p []byte) {
					websockets.HUB.ServerMsg <- websockets.ServerMsg{messageType, p}
				})
				websockets.Client.Run()
			}

		case "error":
			result["sucess"] = "false"
		}

		js, err := json.Marshal(result)

		w.Header().Set("Content-Type", "application/json")
		logger.Info.Println("response to login page:", string(js))
		_, err = w.Write(js)
		if err != nil {
			logger.Error.Println("fail to response login:", err)
		}
	}).Methods("POST")

	r.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {

		logger.Info.Println("get signup post request:", r.PostForm)
		localJson := parsePostForm(r)
		signupMap, _ := parsePostResponse(config.CloudSignupURL, localJson)

		var result = make(map[string]string)
		switch signupMap["result"] {
		case "success":
			result["success"] = "ok"
			//get session
		case "error":
			result["success"] = "false"
		}

		w.Header().Set("Content-Type", "application/json")
		js, _ := json.Marshal(result)
		w.Write(js)

	}).Methods("POST")

}

func setWebsocket(r *mux.Router) {
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("websocket established:", r.RemoteAddr)
		upgrader := websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
		}
		websockets.HUB.Register <- c
		defer closeWebsocket(c)
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				break
			}
			websockets.HUB.BroadMsg <- websockets.ClientMsg{c, mt, message}
			//logger.Info.Println("received from local:", string(message), "message type:", mt)
		}
	})
}

func closeWebsocket(conn *websocket.Conn) {
	conn.Close()
	websockets.HUB.Unregister <- conn
	logger.Info.Println("websocket disconnected:", conn.RemoteAddr())
}

func onWebsocketServerReceived(conn *websocket.Conn, data []byte) {
	msg, err := websockets.DecodeMsg(data)
	if err != nil {
		return
	}
	switch msg["action"] {
	//case "stop-drawing":
	//	conn.WriteMessage()
	case "start-draw":
		conn.WriteMessage(websocket.TextMessage, data)
	case "append-user":

	}
}

func parsePostForm(r *http.Request) []byte {
	err := r.ParseForm()
	if err != nil {
		logger.Warning.Println("can not parse form form the post:", err)
		return []byte{}
	}
	var message = make(map[string]string)
	for k, v := range r.PostForm {
		message[k] = v[0]
	}

	resJson, err := json.Marshal(message)
	if err != nil {
		logger.Warning.Println("cannot marshal postform to json:", err)
		logger.Warning.Println("the form is:", r.PostForm)
		return []byte{}
	}
	return resJson
}

func parsePostResponse(url string, jsonToPost []byte) (map[string]interface{}, []*http.Cookie) {
	res, err := http.Post(url, "application/json", bytes.NewReader(jsonToPost))
	if err != nil {
		logger.Warning.Println("failed to post cloud:", err)
		return make(map[string]interface{}), []*http.Cookie{}
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Warning.Println("failed to read post response body:", err)
		return make(map[string]interface{}), []*http.Cookie{}
	}
	resMap := make(map[string]interface{})
	err = json.Unmarshal(body, &resMap)
	if err != nil {
		logger.Warning.Println("failed to turn body to json:", err)
		logger.Warning.Println("the body is:", string(body))
		return make(map[string]interface{}), []*http.Cookie{}
	}
	return resMap, res.Cookies()
}

func getSessionFromCookie(cookies []*http.Cookie) string {
	for _, c := range cookies {
		if c.Name == "sessionid" {
			return c.Value
		}
	}
	return ""
}
