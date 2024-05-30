package kuro

import (
	"log"
	"net/http"
)

type Server struct {
	clients    map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client

	roomManager *roomManager

	events        map[string]EventCallback
	beforeUpgrade BeforeUpgradeFunc
	onDisconnect  OnDisconnectFunc
	onConnect     OnConnectFunc
}

type BeforeUpgradeFunc func(http.ResponseWriter, *http.Request) error
type OnDisconnectFunc func(c *Client)
type OnConnectFunc func(c *Client)

type Message struct {
	Client  *Client
	Message []byte
}

func New() *Server {
	return &Server{
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),

		roomManager: newRoomManager(),
		events:      make(map[string]EventCallback),
	}
}

func (h *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serveWs(h, w, r)
}

func (s *Server) Join(roomName string, c *Client) {
	log.Println("here")
	s.roomManager.Join(roomName, c)
}

func (s *Server) Leave(roomName string, c *Client) {
	s.roomManager.Leave(roomName, c)
}

func (s *Server) Emit(roomName, eventName string, payload []byte) {
	s.roomManager.Emit(roomName, eventName, payload)
}

//

func (s *Server) OnConnect(callback OnConnectFunc) {
	s.onConnect = callback
}

func (s *Server) OnDisconnect(callback OnDisconnectFunc) {
	s.onDisconnect = callback
}

func (h *Server) OnEvent(event string, callback EventCallback) {
	if _, ok := h.events[event]; !ok {
		h.events[event] = callback
	}
}

func (h *Server) eventHandler(message *Message) {
	event, err := bytesToEvent(message.Message)
	if err != nil {
		log.Println(err)
		return
	}

	if callback, ok := h.events[event.Event]; ok {
		go callback(message.Client, event.Payload)
	}
}

func (h *Server) Run() {
	for {
		select {
		case client := <-h.register:
			h.onConnect(client)
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.onDisconnect(client)
		case message := <-h.broadcast:
			h.eventHandler(message)
			// for client := range h.clients {
			// 	select {
			// 	case client.send <- message:
			// 	default:
			// 		close(client.send)
			// 		delete(h.clients, client)
			// 	}
			// }
		}
	}
}
