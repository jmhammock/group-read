package room

import (
	"context"
	"database/sql"
	"time"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{
		db: db,
	}
}

func (r RoomRepository) context(seconds int64) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
}

func (r RoomRepository) Create(room *Room) error {
	// q := `INSERT INTO rooms (name)
	// VALUES ($1, $2);`

	// c, cancel := r.context(6)
	// defer cancel()

	// // _, err := r.db.ExecContext(c, q, room.Name)
	// return err
	return nil
}

func (r RoomRepository) GetById(id string) (*Room, error) {
	// q := `SELECT id, name
	// FROM rooms
	// WHERE id = $1;`

	// c, cancel := r.context(3)
	// defer cancel()

	// row := r.db.QueryRowContext(c, q, id)
	// var room *Room
	// err := row.Scan(&room.Id, &room.Name)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
