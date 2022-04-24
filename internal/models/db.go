package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNoResults = errors.New("no results found")
)

type DB struct {
	DB *sql.DB
}

func NewDB(db *sql.DB) *DB {
	return &DB{
		DB: db,
	}
}

func (db DB) newContext(seconds uint16) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
}

func (db DB) Query(query string, args ...any) (*sql.Rows, error) {
	c, cancel := db.newContext(5)
	defer cancel()
	return db.DB.QueryContext(c, query, args)
}

func (db DB) QueryRow(query string, args ...any) *sql.Row {
	c, cancel := db.newContext(5)
	defer cancel()
	return db.DB.QueryRowContext(c, query, args)
}

func (db DB) Execute(query string, args ...any) (sql.Result, error) {
	c, cancel := db.newContext(5)
	defer cancel()
	return db.DB.ExecContext(c, query, args)
}
