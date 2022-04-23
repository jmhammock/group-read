package application

import (
	"github.com/jmhammock/ereader/cmd/web/book"
	"github.com/jmhammock/ereader/cmd/web/room"
	"github.com/jmhammock/ereader/cmd/web/user"
)

type Application struct {
	BookRepository *book.BookRepository
	RoomRepository *room.RoomRepository
	UserRepository *user.UserRepository
	Rooms          map[string]*room.Room
}
