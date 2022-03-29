package room

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type roomRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *roomRepository {
	return &roomRepository{
		db: db,
	}
}

func (r roomRepository) context(seconds int64) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
}

func (r roomRepository) Create(room *Room) error {
	id := uuid.NewString()
	q := `INSERT INTO rooms (id, book)
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
