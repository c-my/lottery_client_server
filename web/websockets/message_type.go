package websockets

import (
	"encoding/json"
)

type WsMessage map[string]interface{}

// DecodeMsg decodes ws message received from local
func DecodeMsg(data []byte) (WsMessage, error) {
	var msg WsMessage
	err := json.Unmarshal(data, &msg)
	return msg, err
}
