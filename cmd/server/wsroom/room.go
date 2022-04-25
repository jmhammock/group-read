package wsroom

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jmhammock/ereader/cmd/server/events"
)

var (
	ErrMemberAlreadyExists = errors.New("member already exists")
)

type WSRoom struct {
	Id            string
	Members       map[string]*Client
	BroadcastChan chan events.Event
	mu            *sync.Mutex
}

type WSRooms []WSRoom

func NewWSRoom(id string) *WSRoom {
	return &WSRoom{
		Id:            id,
		Members:       make(map[string]*Client),
		BroadcastChan: make(chan events.Event),
	}
}

func (r *WSRoom) AddMember(c *Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.Members[c.Id]; exists == true {
		return ErrMemberAlreadyExists
	}
	r.Members[c.Id] = c
	return nil
}

func (r *WSRoom) RemoveMember(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Members, id)
}

func (r *WSRoom) MembersLen() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.Members)
}

func (r *WSRoom) Close(e events.Event) {
	r.BroadcastChan <- e
	close(r.BroadcastChan)
	r.mu.Lock()
	r.Members = make(map[string]*Client)
	r.mu.Unlock()
}

func (r *WSRoom) Broadcaster() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for event := range r.BroadcastChan {
			fmt.Println(event)
			for id, client := range r.Members {
				if event.SenderId != id {
					err := client.Conn.WriteJSON(event)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
		wg.Done()
	}()
}

type Client struct {
	Id   string
	Conn *websocket.Conn
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Id:   uuid.NewString(),
		Conn: conn,
	}
}
