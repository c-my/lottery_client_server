package routers

import (
	"bytes"
	"encoding/json"
	"github.com/c-my/lottery_client_server/config"
	"github.com/c-my/lottery_client_server/datamodels"
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
	//subRouter := r.Host(parent).Subrouter()
	//subRouter := r.Host("0.0.0.0：1923").Subrouter()
	//subRouter := r.Path("").Subrouter()

	setStatic(r)
	setGet(r)
	setPost(r)
	setWebsocket(r)
	//setStatic(subRouter)
	//setGet(subRouter)
	//setPost(subRouter)
	//setWebsocket(subRouter)
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
	}).Methods("GET")

	r.HandleFunc("/get-activities", func(writer http.ResponseWriter, request *http.Request) {

		response, err := getWithSession(config.CloudGetACtivities, nil, "")
		if err != nil {
			logger.Error.Println("failed to get activities list from cloud:", err)
		}
		defer response.Body.Close()
		resBody, err := ioutil.ReadAll(response.Body)
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write(resBody)
		if err != nil {
			logger.Error.Println("failed to response exist user to local:", err)
		}
	}).Methods("GET")

	r.HandleFunc("/get-participants/{act-id}", func(writer http.ResponseWriter, request *http.Request) {
		//queryActivity := request.Form["activity"]
		// get user list
		//actID := mux.Vars(request)["act-id"]
		writer.Header().Set("Content-Type", "application/json")
	}).Methods("GET")

	r.HandleFunc("/console", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "views/console.html")
	}).Methods("GET")

	r.HandleFunc("/screen", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "PrizeDraw.html")
	}).Methods("GET")

	r.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "views/index.html")
	}).Methods("GET")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "views/login.html")
	}).Methods("GET")

	r.HandleFunc("/start-menu", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "views/start-menu.1.html")
	}).Methods("GET")

	r.HandleFunc("/rtmp", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "rtmp.html")
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
			websockets.SessionStr = session
			//logger.Info.Println("trying connecting to cloud with session:", session)
			//websockets.Client, err = websockets.NewWebsocketClient(config.CloudWsServer, session)
			if err != nil {
				logger.Warning.Println("failed to connect to cloud", err)
			} else {
				//websockets.Client.SetHandler(func(wsc *websockets.WebsocketClient, messageType int, p []byte) {
				//	websockets.HUB.ServerMsg <- websockets.ServerMsg{messageType, p}
				//})
				//websockets.Client.Run()
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

	r.HandleFunc("/append-activity", func(w http.ResponseWriter, r *http.Request) {
		//TODO: append in database
		//TODO: request json format

		localJson := parsePostForm(r)
		logger.Info.Println("receive post from local:", r.URL)
		logger.Info.Println("the postform is:", string(localJson))
		actName := r.PostForm["name"][0]

		response, err := postWithSession(config.CloudAppendActivities, localJson, "")
		if err != nil {
			logger.Error.Println("failed to get activities list from cloud:", err)
		}
		defer response.Body.Close()

		appendActRes := parseResponseBody(response)

		if len(appendActRes) == 0 {
			return
		}
		actID, _ := appendActRes["activity_id"].(int)
		newAct := datamodels.Activity{Id: actID, Name: actName}
		controllers.ActivityControl.Append(newAct)
	}).Methods("POST")

	r.HandleFunc("/bg-img", func(writer http.ResponseWriter, request *http.Request) {

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
				logger.Error.Println("failed to read from local:", err)
				break
			}
			websockets.HUB.BroadMsg <- websockets.ClientMsg{c, mt, message}
			//logger.Info.Println("received from local:", string(message), "message type:", mt)
		}
	})
}

func closeWebsocket(conn *websocket.Conn) {
	err := conn.Close()
	if err != nil {
		logger.Error.Println("failed to close websocket connection:", err)
	}
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

func postWithSession(url string, jsonToPost []byte, session string) (response *http.Response, err error) {
	return requestWithSession(url, "POST", jsonToPost, session)
}

func getWithSession(url string, jsonToGet []byte, session string) (response *http.Response, err error) {
	return requestWithSession(url, "GET", jsonToGet, session)
}

func requestWithSession(url string, method string, body []byte, session string) (response *http.Response, err error) {
	client := http.Client{}
	logger.Info.Print(method, ": send to cloud (with cookie):", string(body))
	logger.Info.Println("cookies is:", websockets.SessionStr)
	request, err := http.NewRequest(method, url, bytes.NewReader(body))

	cookies := session
	if cookies == "" {
		cookies = websockets.SessionStr
	}
	if request == nil {
		panic("null request")
	}
	request.Header.Add("Cookie", cookies)
	response, err = client.Do(request)

	return
}

func parsePostResponse(url string, jsonToPost []byte) (map[string]interface{}, []*http.Cookie) {
	res, err := http.Post(url, "application/json", bytes.NewReader(jsonToPost))
	if err != nil {
		logger.Warning.Println("failed to post cloud:", err)
		return make(map[string]interface{}), []*http.Cookie{}
	}
	resMap := parseResponseBody(res)
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

func parseResponseBody(response *http.Response) map[string]interface{} {
	body, err := ioutil.ReadAll(response.Body)
	logger.Info.Println("receive body from cloud:", string(body))
	if err != nil {
		logger.Warning.Println("failed to read post response body:", err)
		return make(map[string]interface{})
	}
	resMap := make(map[string]interface{})
	err = json.Unmarshal(body, &resMap)
	if err != nil {
		logger.Warning.Println("receive some shit from cloud:", err)
		logger.Warning.Println("the body is:", string(body))
		return make(map[string]interface{})
	}
	return resMap
}
