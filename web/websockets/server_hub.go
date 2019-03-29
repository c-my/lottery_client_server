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

type Hub struct {
	Clients    map[*websocket.Conn]bool
	BroadMsg   chan ClientMsg
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
}

type ClientMsg struct {
	Client  *websocket.Conn
	Message []byte
}

func newHub() *Hub {
	return &Hub{
		Clients:    make(map[*websocket.Conn]bool),
		BroadMsg:   make(chan ClientMsg),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
}

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
	for client, _ := range h.Clients {
		client.WriteMessage(messageType, data)
	}
}

func (h *Hub) Broadcast(conn *websocket.Conn, messageType int, data []byte) {
	for client, _ := range h.Clients {
		if client != conn {
			logger.Info.Println("message delivered to local: ", string(data))
			client.WriteMessage(messageType, data)
		}
	}
}

func (h *Hub) handleMessage(msg *ClientMsg) {
	data := msg.Message
	conn := msg.Client
	m, err := DecodeMsg(data)
	if err != nil {
		return
	}
	switch m["action"] {
	case "start-drawing":
		h.Broadcast(conn, websocket.TextMessage, data)
	case "stop-drawing":
		// relay the message
		h.Broadcast(conn, websocket.TextMessage, data)
	}
}
