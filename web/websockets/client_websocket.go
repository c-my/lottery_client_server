package websockets

import "log"
import (
	gwebsocket "github.com/gorilla/websocket"
)

// RecvHandler is a callback type
type RecvHandler func(wsc *WebsocketClient, messageType int, p []byte)

// WebsocketClient is a ws object that can connect to a
// websocket server.
type WebsocketClient struct {
	conn    *gwebsocket.Conn
	handler RecvHandler
}

// SendMessage sends ws server a text message
func (ws *WebsocketClient) SendMessage(msg string) error {
	log.Println("message delivered: " + msg)
	return ws.conn.WriteMessage(gwebsocket.TextMessage, []byte(msg))
}

// SetHandler sets a handler that will be called when message received
func (ws *WebsocketClient) SetHandler(handler RecvHandler) {
	ws.handler = handler
}

// Run launches the websocket receive loop
func (ws *WebsocketClient) Run() {
	go ws.receiveLoop(ws.conn, ws.handler)
}

// NewClientWs returns a
func NewWebsocketClient(url string) *(WebsocketClient) {

	c, _, err := gwebsocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println(err)
	}
	var wsc = WebsocketClient{
		conn: c,
	}
	return &wsc
}

func (wsc *WebsocketClient) receiveLoop(conn *gwebsocket.Conn, handler RecvHandler) {
	done := make(chan struct{})
	defer close(done)
	for {
		msgTyp, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			continue
		}
		handler(wsc, msgTyp, message)
	}
}
