package websockets

import (
	"github.com/c-my/lottery_client_server/web/logger"
	"github.com/gorilla/websocket"
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

func (h *Hub) SendAll(messageType int, data []byte) {
	logger.Info.Println("message delivered to all:", string(data), "message type:", messageType)
	for client, _ := range h.Clients {
		client.WriteMessage(messageType, data)
	}
}

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
			// relay the message
			h.Broadcast(conn, websocket.TextMessage, data)
		case "manual-import":
		case "switch-page":
		case "activity-start-time":
		case "show-activity":
		case "hide-activity":
		case "disable-lucky":
		case "part-update":

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
		logger.Warning.Println("unrecognized message received from cloud: ", string(data))
		return
	}
	switch mt {
	case websocket.TextMessage:
		switch m["action"] {
		case "append-user":
			h.SendAll(mt, data)
		case "send-danmu":
			h.SendAll(mt, data)
		case "modify-activity":
		case "participants":
		case "activity-info":
		}
	case websocket.BinaryMessage:
		break
	case websocket.PingMessage:
		break
	case websocket.PongMessage:
		break
	}
}
