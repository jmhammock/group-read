package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/jmhammock/ereader/cmd/web/application"
	"github.com/jmhammock/ereader/cmd/web/book"
	"github.com/jmhammock/ereader/cmd/web/handlers"
	"github.com/jmhammock/ereader/cmd/web/room"
	"github.com/jmhammock/ereader/cmd/web/user"
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "ereader")
	if err != nil {
		log.Fatal(err)
	}

	app := &application.Application{
		BookRepository: book.NewRepository(db),
		RoomRepository: room.NewRepository(db),
		UserRepository: user.NewRepository(db),
	}

	r := httprouter.New()
	r.GET("/ws", handlers.SocketHandler(app))
	r.GET("/", handlers.Home)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
