// Filename: internal/data/models.go

package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrEditConflict       = errors.New("edit conflict")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateEmail     = errors.New("duplicate email address")
)

// A wrapper for out data models
type Models struct {
	Forum ForumModel
	User  UserModel
}

// NewModels() allows us to create new models
func NewModels(db *sql.DB) *Models {
	return &Models{
		Forum: ForumModel{DB: db},
		User:  UserModel{DB: db},
	}
}
