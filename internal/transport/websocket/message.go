package websocket

import (
	"api/internal/core"
	"api/pkg/kuro"
	"encoding/json"
	"log"
	"strconv"
)

type Message struct {
	ToId           int    `json:"to_id"`
	MessageType    string `json:"message_type"`
	MessageContent string `json:"message_content"`
}

func (w *WebSocket) handleMessage(socket *kuro.Client, payload []byte) {
	var message Message
	if err := json.Unmarshal(payload, &message); err != nil {
		log.Println(err)
		return
	}

	if _, err := w.users.GetOneById(message.ToId); err != nil {
		log.Println("user not found")
		return
	}

	userId := socket.Get("userId").(int)

	chatId, err := w.chats.GetChatIdByMembers(message.ToId, userId)
	if err != nil {
		if newChatId, err := w.chats.Add(userId, message.ToId); err != nil {
			log.Println("can't create chat: ", err)
			return
		} else {
			chatId = newChatId
		}
	}

	newMessage, err := w.chats.AddMessage(&core.Message{
		ChatId:   chatId,
		SenderId: userId,
		TypeId:   0,
		Content:  message.MessageContent,
	})
	if err != nil {
		log.Println(err)
		return
	}

	output, err := json.Marshal(newMessage)
	if err != nil {
		log.Println(err)
		return
	}

	w.server.Emit(userRoomPrefix+strconv.Itoa(userId), messageEventPrefix+strconv.Itoa(chatId), output)
	w.server.Emit(userRoomPrefix+strconv.Itoa(message.ToId), messageEventPrefix+strconv.Itoa(chatId), output)

	w.server.Emit(userRoomPrefix+strconv.Itoa(userId), newMessageEventPrefix, output)
	w.server.Emit(userRoomPrefix+strconv.Itoa(message.ToId), newMessageEventPrefix, output)
}
