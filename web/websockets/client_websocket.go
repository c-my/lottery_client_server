package websockets

import (
	"github.com/c-my/lottery_client_server/web/logger"
	"github.com/gorilla/websocket"
)

var Client *WebsocketClient
var SessionStr string

// RecvHandler is a callback type
type RecvHandler func(wsc *WebsocketClient, messageType int, p []byte)

// WebsocketClient is a ws object that can connect to a
// websocket server.
type WebsocketClient struct {
	conn    *websocket.Conn
	handler RecvHandler
}

// SendMessage sends ws server a text message
func (ws *WebsocketClient) SendMessage(msg string) error {
	logger.Info.Println("message delivered to cloud:" + msg)
	if ws.conn == nil {
		logger.Error.Println("conn is nil") //do not delete this, do not ask why
	}
	return ws.conn.WriteMessage(websocket.TextMessage, []byte(msg))
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
func NewWebsocketClient(url string) (*(WebsocketClient), error) {

	header := make(map[string][]string)
	header["Cookie"] = []string{SessionStr}
	c, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		logger.Error.Println("failed to connect to cloud:", err)
	}
	var wsc = WebsocketClient{
		conn: c,
	}
	return &wsc, err
}

func (wsc *WebsocketClient) receiveLoop(conn *websocket.Conn, handler RecvHandler) {
	done := make(chan struct{})
	defer close(done)
	for {
		msgTyp, message, err := conn.ReadMessage()
		if err != nil {
			logger.Error.Println("failed to read message from cloud:", err)
			continue
		}
		handler(wsc, msgTyp, message)
	}
}
