package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID uint

	Login string

	Password string

	FirstName string
	LastName  string

	Email string

	CreatedAt time.Time
	UpdatedAt *time.Time
}

func (u *User) ValidPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err
}
