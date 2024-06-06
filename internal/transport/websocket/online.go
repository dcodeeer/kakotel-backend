package websocket

import (
	"api/pkg/kuro"
	"encoding/json"
	"log"
	"strconv"
)

func (w *WebSocket) IsUserOnline(userId int) bool {
	err := w.server.RoomExists(userRoomPrefix + strconv.Itoa(userId))
	return err == nil
}

type onlineDto struct {
	UserId int `json:"user_id"`
}

func (w *WebSocket) handleOnlineListen(socket *kuro.Client, payload []byte) {
	var dto onlineDto
	if err := json.Unmarshal(payload, &dto); err != nil {
		log.Println(err)
		return
	}

	w.server.Join(userRoomPrefix+strconv.Itoa(dto.UserId), socket)
}
