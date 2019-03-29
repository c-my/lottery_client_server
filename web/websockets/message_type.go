package websockets

import (
	"encoding/json"
	"github.com/c-my/lottery_client_server/web/logger"
)

type WsMessage map[string]interface{}

func DecodeMsg(data []byte) (WsMessage, error) {
	var msg WsMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		logger.Warning.Println("cannot encode command received from local: ", string(data))
	}
	return msg, err
}
