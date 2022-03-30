package room

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/google/uuid"
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
	id := uuid.NewString()
	q := `INSERT INTO rooms (id, book_id)
	VALUES ($1, $2);`

	c, cancel := r.context(6)
	defer cancel()

	_, err := r.db.ExecContext(c, q, id, room.Book)
	if err != nil {
		return err
	}

	q = "INSERT INTO room_users (room_id, user_id) VALUES"
	n := len(room.UsersIds)
	for i := 0; i < n; i++ {
		next := strconv.Itoa(i + 2)
		q += "($1, $" + next + ")"
		if i != n-1 {
			q += ", "
		} else {
			q += ";"
		}
	}

	args := []interface{}{id}
	for _, id := range room.UsersIds {
		args = append(args, id)
	}

	_, err = r.db.ExecContext(c, q, args...)
	return err
}

func (r RoomRepository) GetById(id string) (*Room, error) {
	q := `SELECT r.id, b.book, ru.user_id
	JOIN books b ON r.book_id = b.id
	JOIN room_users ru ON r.id = ru.room_id
	WHERE r.id = $1;`

	c, cancel := r.context(3)
	defer cancel()

	rows, err := r.db.QueryContext(c, q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var room *Room
	room.UsersIds = make([]int64, 0)
	for rows.Next() {
		var user_id int64
		rows.Scan(
			&room.Id,
			&room.Book,
			&user_id,
		)

		room.UsersIds = append(room.UsersIds, user_id)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return room, nil
}
