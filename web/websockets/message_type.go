package websockets

import (
	"encoding/json"
	"github.com/c-my/lottery_client_server/web/logger"
)

type WsMessage map[string]interface{}

// DecodeMsg decodes ws message received from local
func DecodeMsg(data []byte) (WsMessage, error) {
	var msg WsMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		logger.Warning.Println("unrecognized message received from local: ", string(data))
	}
	return msg, err
}
