package wsroom

import (
	"errors"
	"fmt"
	"sync"

	"github.com/jmhammock/ereader/cmd/server/events"
)

var (
	ErrMemberAlreadyExists = errors.New("member already exists")
	ErrMemberDoesNotExist  = errors.New("member does not exist")
)

type WSRoom struct {
	Id          string
	Members     map[string]Client
	messageChan chan *events.Event
	mu          *sync.Mutex
}

type WSRooms []WSRoom

func NewWSRoom(id string) *WSRoom {
	return &WSRoom{
		Id:          id,
		Members:     make(map[string]Client),
		messageChan: make(chan *events.Event),
		mu:          &sync.Mutex{},
	}
}

func (r *WSRoom) Join(c Client) error {
	r.mu.Lock()
	if _, exists := r.Members[c.GetId()]; exists == true {
		return ErrMemberAlreadyExists
	}
	r.Members[c.GetId()] = c
	r.mu.Unlock()
	r.messageChan <- events.NewJoinEvent(c.GetId())
	return nil
}

func (r *WSRoom) Leave(id string) {
	r.mu.Lock()
	if m, exists := r.Members[id]; exists == true {
		m.Close()
		delete(r.Members, id)
	}
	r.mu.Unlock()
	r.messageChan <- events.NewLeaveEvent(id)
}

func (r *WSRoom) MembersLen() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.Members)
}

func (r *WSRoom) Close(id string) {
	r.messageChan <- events.NewWSRoomCloseEvent(id)
	close(r.messageChan)
	r.mu.Lock()
	for _, c := range r.Members {
		c.Close()
	}
	r.Members = make(map[string]Client)
	r.mu.Unlock()
}

func (r *WSRoom) Broadcast(e *events.Event) {
	r.messageChan <- e
}

func (r *WSRoom) Receiver() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for e := range r.messageChan {
			for id, client := range r.Members {
				if e.SenderId != id {
					err := client.Send(e)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}()
}
