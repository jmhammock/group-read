package book

import (
	"context"
	"database/sql"
	"time"
)

type BookRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

func (r BookRepository) context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

func (r BookRepository) Create(b *Book) error {
	q := `INSERT INTO books (title, location) VALUES ($1, $2);`

	c, cancel := r.context()
	defer cancel()

	_, err := r.db.ExecContext(c, q, b.Title, b.Location)
	if err != nil {
		return err
	}

	return nil
}

func (r BookRepository) GetById(id int64) (*Book, error) {
	q := `SELECT id, title, location FROM books WHERE id = $1;`

	c, cancel := r.context()
	defer cancel()

	row := r.db.QueryRowContext(c, q, id)

	var book *Book
	err := row.Scan(
		&book.Id,
		&book.Title,
		&book.Location,
	)
	if err != nil {
		return nil, err
	}

	return book, nil
}
