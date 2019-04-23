package websockets

import (
	"encoding/json"
	"github.com/c-my/lottery_client_server/config"
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_client_server/services"
	"github.com/c-my/lottery_client_server/web/controllers"
	"github.com/c-my/lottery_client_server/web/logger"
	"github.com/c-my/lottery_client_server/web/tools"
	"github.com/gorilla/websocket"
	"reflect"
)

var HUB *Hub

func init() {
	HUB = newHub()
	go HUB.Run()
}

// Hub is used to control all the client
type Hub struct {
	Clients    map[*websocket.Conn]bool
	BroadMsg   chan ClientMsg
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	ServerMsg  chan ServerMsg
}

// ClientMsg is the content sent to channel when some client
// receiving message
type ClientMsg struct {
	Client      *websocket.Conn
	MessageType int
	Message     []byte
}

type ServerMsg struct {
	MessageType int
	Message     []byte
}

func newHub() *Hub {
	return &Hub{
		Clients:    make(map[*websocket.Conn]bool),
		BroadMsg:   make(chan ClientMsg),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
		ServerMsg:  make(chan ServerMsg),
	}
}

// Run makes Hub begins monitoring channels
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			delete(h.Clients, client)
		case message := <-h.BroadMsg:
			logger.Info.Println("received from local: ", string(message.Message), "message type:", message.MessageType)
			h.handleClientMessage(&message)
		case message := <-h.ServerMsg:
			logger.Info.Println("received from cloud: ", string(message.Message), "message type:", message.MessageType)
			h.handleServerMessage(&message)
		}
	}
}

// SendAll sends the message to all clients connected to us
func (h *Hub) SendAll(messageType int, data []byte) {
	logger.Info.Println("message delivered to all:", string(data), "message type:", messageType)
	for client, _ := range h.Clients {
		client.WriteMessage(messageType, data)
	}
}

// Broadcase sends the message to all clients connected to us except the sender
func (h *Hub) Broadcast(conn *websocket.Conn, messageType int, data []byte) {
	logger.Info.Println("broadcast message:", string(data), "message type:", messageType)
	for client, _ := range h.Clients {
		if client != conn {
			client.WriteMessage(messageType, data)
		}
	}
}

func (h *Hub) handleClientMessage(msg *ClientMsg) {
	data := msg.Message
	mt := msg.MessageType
	conn := msg.Client
	m, err := DecodeMsg(data)
	if err != nil {
		logger.Warning.Println("unrecognized message received from local: ", string(data))
		return
	}
	switch mt {
	case websocket.TextMessage:
		switch m["action"] {
		case "start-drawing":
			h.Broadcast(conn, websocket.TextMessage, data)
		case "stop-drawing":
			json := generateLuckyDog(&m)
			h.SendAll(websocket.TextMessage, json)
			h.Broadcast(conn, websocket.TextMessage, data)
		case "start-activity":
			tools.ConsoleConfig.IsActivityUnfinished = true
			startActivity(&m)
		case "end-activity":
			tools.ConsoleConfig.IsActivityUnfinished = false
		case "manual-import":
		case "switch-page":
			h.Broadcast(conn, websocket.TextMessage, data)
		case "activity-start-time":
		case "show-activity":
			h.Broadcast(conn, websocket.TextMessage, data)
		case "hide-activity":
			h.Broadcast(conn, websocket.TextMessage, data)
		case "disable-lucky":
			h.Broadcast(conn, websocket.TextMessage, data)
		case "danmu-switch":
			h.Broadcast(conn, websocket.TextMessage, data)
		case "danmu-check-switch":
			h.Broadcast(conn, websocket.TextMessage, data)
		case "part-update":
			updateConfig(&m)

		}
	case websocket.BinaryMessage:
		break
	case websocket.PingMessage:
		break
	case websocket.PongMessage:
		break
	}
}

func (h *Hub) handleServerMessage(msg *ServerMsg) {
	data := msg.Message
	mt := msg.MessageType
	m, err := DecodeMsg(data)
	if err != nil {
		if string(data) == config.CloudConfirmCode {
			logger.Info.Println("connect confirm code:", string(data))
		} else {
			logger.Warning.Println("unrecognized message received from cloud: ", string(data))
		}
		return
	}
	switch mt {
	case websocket.TextMessage:
		switch m["action"] {
		case "append-user":
			appendUser(&m)
			h.SendAll(mt, data)
		case "send-danmu":
			appendDanmu(&m)
			h.SendAll(mt, data)
		case "modify-activity":
		case "participants":
			fillUsers(&m)
		case "activity-info":
			addActInfo(&m)
		}
	case websocket.BinaryMessage:
		break
	case websocket.PingMessage:
		break
	case websocket.PongMessage:
		break
	}
}

func appendUser(msg *WsMessage) {
	userJson, _ := json.Marshal((*msg)["content"])
	userToAdd := datamodels.User{}
	json.Unmarshal(userJson, &userToAdd)
	controllers.UserControl.Append(userToAdd)
}

func appendDanmu(msg *WsMessage) {
	danmuJson, _ := json.Marshal((*msg)["content"])
	danmuToAdd := datamodels.BulletComment{}
	json.Unmarshal(danmuJson, &danmuToAdd)
	//danmuToAdd := (*msg)["content"].(datamodels.BulletComment)
	controllers.DanmuControl.Append(danmuToAdd)
}

func fillUsers(msg *WsMessage) {
	usersToAdd := (*msg)["content"].([]datamodels.User)
	for _, u := range usersToAdd {
		controllers.UserControl.Append(u)
	}
}

func addActInfo(msg *WsMessage) {
	actInfo := (*msg)["content"].(datamodels.Activity)
	controllers.ActivityControl.Append(actInfo)
}

func generateLuckyDog(msg *WsMessage) []byte {
	_, ok := (*msg)["content"].(string)
	if !ok {
		user := controllers.UserControl.RandomlyGet()
		json, _ := AddAction("who-is-lucky-dog", user)
		return json
	}
	prizeID := (*msg)["content"].(string)
	userList := controllers.UserControl.RandomlyGetAll()
	for _, user := range userList {
		if !services.WinnerServicer.AlreadyWin(user.ID, prizeID) {
			// add the winner record
			services.WinnerServicer.AddWinner(user.ID, prizeID)
			// generate json bytes
			json, _ := AddAction("who-is-lucky-dog", user)
			return json
		}
	}
	// TODO:deal when no available lucky-dog
	user := datamodels.User{}
	json, _ := AddAction("no-available", user)
	return json
}

func startActivity(msg *WsMessage) {
	actID := (*msg)["content"].(string)
	Client, _ = NewWebsocketClient(config.CloudWsServer + "/" + actID)
	Client.SetHandler(func(wsc *WebsocketClient, messageType int, p []byte) {
		HUB.ServerMsg <- ServerMsg{messageType, p}
	})
	Client.Run()
}

func updateConfig(msg *WsMessage) {
	content := (*msg)["content"]
	bytes, _ := json.Marshal(content)
	configMsg, _ := DecodeMsg(bytes)

	v := reflect.ValueOf(&(tools.ConsoleConfig)).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag
		json := tag.Get("json")
		_, ok := configMsg[json].(string)
		if ok {
			logger.Info.Println("updated key:", string(json))
			v.Field(i).SetString(configMsg[json].(string))
			tools.SaveConfigure(config.ConfigureFile)
		}
	}

}
