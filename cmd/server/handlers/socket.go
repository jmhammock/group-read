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

		member := wsroom.NewMember(conn)
		err = wsr.Join(member)
		if err != nil {
			return
		}

		defer func() {
			wsr.Leave(member.GetId())
			if wsr.MembersLen() == 0 {
				app.WSRooms.RemoveRoom(wsr.Id, member.GetId())
			}
			log.Printf("number of rooms: %d", app.WSRooms.RoomsLen())
		}()

		log.Printf(
			"number of %s clients: %d",
			roomId,
			wsr.MembersLen(),
		)

		for {
			var event *events.Event
			err := member.Read(event)
			if err != nil {
				log.Print("error")
				log.Print(err)
				break
			}
			wsr.Broadcast(event)
		}
	}
}
