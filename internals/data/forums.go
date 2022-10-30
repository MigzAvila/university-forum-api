// Filename : internal/data/forums.go

package data

import (
	"database/sql"
	"time"
)

type Forum struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Version     int32     `json:"version"`
}

// define a ForumModel object that wraps a sql.DB connection pool
type ForumModel struct {
	DB *sql.DB
}

type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword []byte    `json:"hashed_password"`
	CreatedAt      time.Time `json:"-"`
	Active         bool      `json:"active"`
}

// func
