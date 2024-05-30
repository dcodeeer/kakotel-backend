package kuro

import (
	"log"
	"sync"
)

type roomManager struct {
	mutex sync.RWMutex
	rooms map[string]map[string]*Client
}

func newRoomManager() *roomManager {
	return &roomManager{
		rooms: make(map[string]map[string]*Client),
	}
}

func (r *roomManager) Join(room string, c *Client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	log.Println("join to room")

	if _, ok := r.rooms[room]; !ok {
		r.rooms[room] = make(map[string]*Client)
	}

	r.rooms[room][c.ID()] = c
}

func (r *roomManager) Leave(room string, c *Client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if clients, ok := r.rooms[room]; ok {
		delete(clients, c.ID())
	}
}

func (r *roomManager) Emit(room, event string, payload []byte) {
	for roomName, _ := range r.rooms {
		log.Println(roomName)
	}
	if clients, ok := r.rooms[room]; ok {
		for _, client := range clients {
			var dto Event
			dto.Event = event
			dto.Payload = payload
			message := dto.toBytes()
			client.send <- message
		}
	}
}
