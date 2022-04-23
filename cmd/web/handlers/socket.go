package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jmhammock/ereader/cmd/web/application"
	"github.com/jmhammock/ereader/cmd/web/events"
	"github.com/jmhammock/ereader/cmd/web/room"
	"github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SocketHandler(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		roomId := p.ByName("id")
		if _, ok := app.Rooms[roomId]; !ok {
			app.Rooms[roomId] = room.NewRoom(roomId)
		}
		readingRoom := app.Rooms[roomId]
		readingRoom.Broadcaster()
		log.Printf("number of rooms: %d", len(app.Rooms))

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		client := room.NewClient(conn)
		readingRoom.Mutex.Lock()
		readingRoom.Members[client.Id] = client
		readingRoom.Mutex.Unlock()

		defer func() {
			delete(readingRoom.Members, client.Id)
			conn.Close()
			leaveEvent := events.Event{
				Type:     "client.leave",
				SenderId: client.Id,
				Data: map[string]interface{}{
					"client": client.Id,
				},
			}
			readingRoom.BroadcastChan <- leaveEvent
			if len(readingRoom.Members) == 0 {
				close(readingRoom.BroadcastChan)
				delete(app.Rooms, readingRoom.Id)
			}
			log.Printf("number of rooms: %d", len(app.Rooms))
		}()

		log.Printf(
			"number of %s clients: %d",
			roomId,
			len(app.Rooms[roomId].Members),
		)

		for {
			var event events.Event
			err := client.Conn.ReadJSON(&event)
			if err != nil {
				log.Print("error")
				log.Print(err)
				break
			}

			readingRoom.BroadcastChan <- event
		}
	}
}
