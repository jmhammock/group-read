package user

import (
	"context"
	"database/sql"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

func (r UserRepository) Create(u *User) error {
	q := `INSERT INTO users (email, password, display_name)
	VALUES ($1, $2, $3, $4);`

	c, cancel := r.context()
	defer cancel()

	_, err := r.db.ExecContext(c, q, u.Email, u.Password, u.DisplayName)

	return err
}

func (r UserRepository) GetByEmail(email string) (*User, error) {
	q := `SELECT id, email, password, display_name FROM users WHERE email = $1;`

	c, cancel := r.context()
	defer cancel()

	row := r.db.QueryRowContext(c, q, email)

	var user *User
	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.DisplayName,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r UserRepository) GetById(id int64) (*User, error) {
	q := `SELECT id, email, password, display_name FROM users WHERE id = $1;`

	c, cancel := r.context()
	defer cancel()

	row := r.db.QueryRowContext(c, q, id)

	var user *User
	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.DisplayName,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
