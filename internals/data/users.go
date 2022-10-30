package data

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword []byte    `json:"hashed_password"`
	CreatedAt      time.Time `json:"-"`
	Active         bool      `json:"active"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

func (m *UserModel) Authenticate(name, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Get(id int) (*User, error) {
	return nil, nil

}
