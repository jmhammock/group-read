package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/jmhammock/ereader/cmd/server/application"
	"github.com/jmhammock/ereader/cmd/server/handlers"
	"github.com/jmhammock/ereader/cmd/server/wsroom"
	"github.com/jmhammock/ereader/internal/models"
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "ereader")
	if err != nil {
		log.Fatal(err)
	}

	wdb := models.NewDB(db)

	app := &application.Application{
		FamilyModel:   models.NewFamilyModel(wdb),
		UserModel:     models.NewUserModel(wdb),
		UserRoleModel: models.NewUserRoleModel(wdb),
		BookModel:     models.NewBookModel(wdb),
		RoomModel:     models.NewRoomModel(wdb),
		WSRooms:       wsroom.NewWSRoomManager(),
	}

	r := httprouter.New()

	r.GET("/api/rooms", handlers.GetRooms(app))
	r.GET("/api/rooms/:id", handlers.GetRoom(app))
	r.PUT("/api/rooms/:id", handlers.UpdateRoom(app))
	r.DELETE("/api/rooms/:id", handlers.DeleteRoom(app))
	r.GET("/ws/:id", handlers.SocketHandler(app))

	log.Fatal(http.ListenAndServe(":8080", r))
}
