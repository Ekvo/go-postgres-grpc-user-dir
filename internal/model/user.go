package model

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrModelUserDifferentID = errors.New("different id's")

	ErrModelUserDateEarly = errors.New("earlier date of recording")
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

func (u *User) NewData(user *User) error {
	if u.ID != user.ID {
		return ErrModelUserDifferentID
	}
	// user.UpdatedAt - check in (deserializer/user_update_decode.go)
	if !u.CreatedAt.UTC().Before(user.UpdatedAt.UTC()) ||
		(u.UpdatedAt != nil && !u.UpdatedAt.UTC().Before(user.UpdatedAt.UTC())) {
		return ErrModelUserDateEarly
	}
	if user.Login != "" {
		u.Login = user.Login
	}
	if user.Password != "" {
		u.Password = user.Password
	}
	if user.FirstName != "" {
		u.FirstName = user.FirstName
	}
	if user.LastName != "" {
		u.LastName = user.LastName
	}
	if user.Email != "" {
		u.Email = user.Email
	}
	u.UpdatedAt = user.UpdatedAt
	return nil
}
