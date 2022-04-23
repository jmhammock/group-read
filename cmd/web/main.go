package main

import (
	"database/sql"
	"embed"
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

//go:embed templates/*
var templateFS embed.FS

//go:embed public/*
var publicFS embed.FS

func main() {
	db, err := sql.Open("sqlite3", "ereader")
	if err != nil {
		log.Fatal(err)
	}

	app := &application.Application{
		BookRepository: book.NewRepository(db),
		RoomRepository: room.NewRepository(db),
		UserRepository: user.NewRepository(db),
		Rooms:          make(map[string]*room.Room),
	}

	r := httprouter.New()

	r.GET("/public/*filepath", handlers.Public(publicFS))
	r.GET("/room/:id", handlers.Room(templateFS, app))
	r.GET("/ws/:id", handlers.SocketHandler(app))

	log.Fatal(http.ListenAndServe(":8080", r))
}
