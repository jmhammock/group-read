package wsroom

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jmhammock/ereader/cmd/server/events"
)

type WSRoom struct {
	Id            string `json:"id"`
	Members       map[string]*Client
	BroadcastChan chan events.Event
	*sync.Mutex
}

func NewWSRoom(id string) *WSRoom {
	return &WSRoom{
		Id:            id,
		Members:       make(map[string]*Client),
		BroadcastChan: make(chan events.Event),
		Mutex:         &sync.Mutex{},
	}
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