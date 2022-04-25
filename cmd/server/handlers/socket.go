package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jmhammock/ereader/cmd/server/application"
	"github.com/jmhammock/ereader/cmd/server/events"
	"github.com/jmhammock/ereader/cmd/server/wsroom"
	"github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SocketHandler(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		roomId := p.ByName("id")
		wsr := wsroom.NewWSRoom(roomId)
		app.WSRooms.AddRoom(wsr)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		client := wsroom.NewClient(conn)
		err = wsr.AddMember(client)
		if err != nil {
			return
		}

		defer func() {
			wsr.RemoveMember(client.Id)
			conn.Close()
			leaveEvent := events.Event{
				Type:     "client.leave",
				SenderId: client.Id,
				Data: map[string]interface{}{
					"client": client.Id,
				},
			}
			wsr.BroadcastChan <- leaveEvent
			if wsr.MembersLen() == 0 {
				closeRoomEvent := events.Event{
					Type:     "room.close",
					SenderId: client.Id,
					Data: map[string]interface{}{
						"room": wsr.Id,
					},
				}
				app.WSRooms.RemoveRoom(wsr.Id, closeRoomEvent)
			}
			log.Printf("number of rooms: %d", app.WSRooms.RoomsLen())
		}()

		log.Printf(
			"number of %s clients: %d",
			roomId,
			wsr.MembersLen(),
		)

		for {
			var event events.Event
			err := client.Conn.ReadJSON(&event)
			if err != nil {
				log.Print("error")
				log.Print(err)
				break
			}

			wsr.BroadcastChan <- event
		}
	}
}
