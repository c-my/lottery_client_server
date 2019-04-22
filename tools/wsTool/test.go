package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/websocket"
)

type TestReturn struct {
	Success string `json:"success"`
}

type Message struct {
	Action  string                 `json:"action"`
	Content map[string]interface{} `json:"content"`
}

var console_ws, recv_ws *websocket.Conn

func console(ws *websocket.Conn) {
	recv_data := map[string]interface{}{}
	console_ws = ws
	for {
		err := websocket.JSON.Receive(ws, &recv_data)
		if err != nil {
			fmt.Printf("err %v", err)
			ws.Close()
			break
		}
		fmt.Printf("Server receive : %v\n", recv_data)
		websocket.JSON.Send(recv_ws, recv_data)
	}
}

func recv(ws *websocket.Conn) {
	recv_data := map[string]interface{}{}
	recv_ws = ws
	for {
		err := websocket.JSON.Receive(ws, &recv_data)
		if err != nil {
			fmt.Printf("err %v", err)
			ws.Close()
			break
		}
		fmt.Printf("Server receive : %v\n", recv_data)
	}
}

func send(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.WriteHeader(200)
	if r.Method == "POST" {
		str := r.PostFormValue("content")
		res_struct := TestReturn{
			Success: "true",
		}
		res, _ := json.Marshal(res_struct)
		fmt.Fprint(w, string(res))
		fmt.Println(str)
		send_data := map[string]interface{}{}
		json.Unmarshal([]byte(str), &send_data)
		websocket.JSON.Send(console_ws, send_data)
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.WriteHeader(200)
	if r.Method == "POST" {
		res_struct := TestReturn{
			Success: "true",
		}
		res, _ := json.Marshal(res_struct)
		fmt.Fprint(w, string(res))
		body, _ := ioutil.ReadAll(r.Body)
		recv_data := Message{}
		json.Unmarshal(body, &recv_data)
		fmt.Println(recv_data)
		websocket.JSON.Send(recv_ws, recv_data)
	}
}

func main() {
	http.HandleFunc("/send", send)
	http.HandleFunc("/post", post)
	http.Handle("/ws", websocket.Handler(console))
	http.Handle("/recv", websocket.Handler(recv))
	err := http.ListenAndServe(":1923", nil)
	fmt.Printf("http %s\n", err)
}
