package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jmhammock/ereader/cmd/web/application"
	"github.com/jmhammock/ereader/cmd/web/events"
	"github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type currentPage struct {
	sync.Mutex
	pageNumber int
}

func (c *currentPage) setPage(p int) {
	c.Lock()
	defer c.Unlock()
	c.pageNumber = p
}

func (c *currentPage) getPage() int {
	c.Lock()
	defer c.Unlock()
	return c.pageNumber
}

func SocketHandler(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		cp := currentPage{
			pageNumber: 0,
		}
		for {
			var event *events.Event
			err := conn.ReadJSON(event)
			if err != nil {
				log.Print(err)
				break
			}

			switch event.Type {
			case "toPage":
				pn, ok := event.Data["pageNumber"].(int)
				if !ok {
					log.Println("page number must be an integer")
					break
				}
				cp.setPage(pn)
				conn.WriteJSON(event)
			case "join":
				err := conn.WriteJSON(events.NewToPageEvent(cp.getPage()))
				if err != nil {
					log.Println(err)
				}
			default:
				conn.WriteJSON(event)
			}
		}
	}
}
