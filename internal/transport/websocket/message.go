package websocket

import (
	"api/pkg/kuro"
	"encoding/json"
	"log"
)

func (w *WebSocket) handleMessage(socket *kuro.Client, payload []byte) {
	var wasd struct {
		UserId  int    `json:"user_id"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(payload, &wasd); err != nil {
		log.Println(err)
		return
	}

	w.server.Emit("room_1", "message", []byte("I'm joined!"))
}
