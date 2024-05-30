package websocket

import (
	"api/pkg/kuro"
	"log"
	"net/http"
)

type WebSocket struct {
	server *kuro.Server
}

func New() *WebSocket {
	return &WebSocket{
		server: kuro.New(),
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

	w.server.Run()
}

func (ws *WebSocket) SetBeforeUpgrade(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (w *WebSocket) OnConnect(socket *kuro.Client) {
	log.Println("here")
	w.server.Join("room_1", socket)
}

func (w *WebSocket) OnDisconnect(socket *kuro.Client) {
	w.server.Leave("room_1", socket)
}
