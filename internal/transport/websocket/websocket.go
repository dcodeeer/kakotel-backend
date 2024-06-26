package websocket

import (
	"api/internal/application"
	"api/pkg/kuro"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type WebSocket struct {
	server *kuro.Server

	users application.IUsers
	chats application.IChats
}

func New(users application.IUsers, chats application.IChats) *WebSocket {
	return &WebSocket{
		server: kuro.New(),

		users: users,
		chats: chats,
	}
}

func (w *WebSocket) GetServer() *kuro.Server {
	return w.server
}

func (w *WebSocket) Run() {
	w.server.SetBeforeUpgrade(w.SetBeforeUpgrade)
	w.server.OnConnect(w.OnConnect)
	w.server.OnDisconnect(w.OnDisconnect)

	w.server.OnEvent("message", w.handleMessage)
	w.server.OnEvent("online", w.handleOnlineListen)

	w.server.Run()
}

func (ws *WebSocket) SetBeforeUpgrade(w http.ResponseWriter, r *http.Request) (map[string]any, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}

	user, err := ws.users.GetByToken(cookie.Value)
	if err != nil {
		return nil, err
	}

	return map[string]any{"userId": user.ID}, nil
}

type handleOnlineDto struct {
	Online   bool   `json:"online"`
	LastSeen string `json:"last_seen"`
}

func (w *WebSocket) OnConnect(socket *kuro.Client) {
	userId := socket.Get("userId").(int)

	w.server.Join(userRoomPrefix+strconv.Itoa(userId), socket)

	bytes, _ := json.Marshal(handleOnlineDto{Online: true, LastSeen: ""})
	w.server.Emit(userRoomPrefix+strconv.Itoa(userId), userRoomPrefix+strconv.Itoa(userId), bytes)
}

func (w *WebSocket) OnDisconnect(socket *kuro.Client) {
	userId := socket.Get("userId").(int)
	if err := w.users.UpdateLastSeen(userId); err != nil {
		log.Println(err)
	}
	w.server.Leave(userRoomPrefix+strconv.Itoa(userId), socket)

	bytes, _ := json.Marshal(handleOnlineDto{Online: false, LastSeen: time.Now().String()})
	w.server.Emit(userRoomPrefix+strconv.Itoa(userId), userRoomPrefix+strconv.Itoa(userId), bytes)
}
