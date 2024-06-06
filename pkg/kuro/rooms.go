package kuro

import (
	"errors"
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

func (r *roomManager) RoomExists(room string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.rooms[room]; !ok {
		return errors.New("room not found")
	}

	return nil
}

func (r *roomManager) Join(room string, c *Client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

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

		if len(clients) == 0 {
			delete(r.rooms, room)
		}
	}
}

func (r *roomManager) Emit(room, event string, payload []byte) {
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
