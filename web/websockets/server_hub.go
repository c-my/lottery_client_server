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
}

// ClientMsg is the content sent to channel when some client
// receiving message
type ClientMsg struct {
	Client      *websocket.Conn
	MessageType int
	Message     []byte
}

func newHub() *Hub {
	return &Hub{
		Clients:    make(map[*websocket.Conn]bool),
		BroadMsg:   make(chan ClientMsg),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
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
			logger.Info.Println("received from local: ", string(message.Message))
			h.handleMessage(&message)
		}
	}
}

func (h *Hub) SendAll(messageType int, data []byte) {
	logger.Info.Println("message delivered to all:", data, "message type:", messageType)
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

func (h *Hub) handleMessage(msg *ClientMsg) {
	data := msg.Message
	mt := msg.MessageType
	conn := msg.Client
	m, err := DecodeMsg(data)
	if err != nil {
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
		}
	case websocket.BinaryMessage:
		break
	case websocket.PingMessage:
		break
	case websocket.PongMessage:
		break
	}

}
