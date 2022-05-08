package wsroom

import (
	"errors"
	"sync"
)

var (
	ErrRoomDoesntExist = errors.New("room does not exist")
)

type WSRoomManager struct {
	WSRooms map[string]*WSRoom
	mu      *sync.Mutex
}

func NewWSRoomManager() *WSRoomManager {
	return &WSRoomManager{
		WSRooms: make(map[string]*WSRoom),
		mu:      &sync.Mutex{},
	}
}

func (w *WSRoomManager) GetRoom(id string) (*WSRoom, error) {
	w.mu.Lock()
	room, exists := w.WSRooms[id]
	w.mu.Unlock()
	if exists == false {
		return nil, ErrRoomDoesntExist
	}
	return room, nil
}

func (w *WSRoomManager) AddRoom(wr *WSRoom) {
	w.mu.Lock()
	if _, exists := w.WSRooms[wr.Id]; exists == true {
		return
	}
	w.WSRooms[wr.Id] = wr
	w.mu.Unlock()
	wr.Receiver()
}

func (w *WSRoomManager) RemoveRoom(roomId, clientId string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if room, exists := w.WSRooms[roomId]; exists == true {
		room.Close(clientId)
		delete(w.WSRooms, roomId)
	}
}

func (w *WSRoomManager) RoomsLen() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.WSRooms)
}
