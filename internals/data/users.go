package data

import (
	"database/sql"
	"time"

	"universityforum.miguelavila.net/internals/validator"
)

type Newuser struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"hashed_password"`
	CreatedAt time.Time `json:"-"`
	Active    bool      `json:"active"`
}

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

func ValidateUser(v *validator.Validator, user *Newuser) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 200, "name", "must no more 200 characters")

	v.Check(user.Email != "", "email", "must be provided")
	v.Check(validator.Matches(user.Email, validator.EmailRX), "email", "must be a valid email")

	v.Check(user.Password != "", "hashedPassword", "must be provided")
	v.Check(len(user.Name) <= 200, "name", "must no more 200 characters")
	v.Check(len(user.Name) <= 200, "name", "must no more 200 characters")

}
