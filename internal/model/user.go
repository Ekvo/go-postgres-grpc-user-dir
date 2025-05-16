package model

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrModelUserDifferentID - used when we try to update user data
	ErrModelUserDifferentID = errors.New("different id's")

	// ErrModelUserDateEarly - mark that the data to update is earlier than the creation date or the previous update
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

// ValidPassword - compare passwords with help 'bcrypt'
func (u *User) ValidPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err
}

// ValidUpdate - !!! u -> old data, user -> new data !!!
// compare u.ID and user.ID -> not equal -> error
// then if u was created not before user update was submitted -> error
// u last update not before user update -> error
func (u *User) ValidUpdate(user *User) error {
	if u.ID != user.ID {
		return ErrModelUserDifferentID
	}
	// user.UpdatedAt - check in (deserializer/user_update_decode.go)
	if !u.CreatedAt.UTC().Before(user.UpdatedAt.UTC()) ||
		(u.UpdatedAt != nil && !u.UpdatedAt.UTC().Before(user.UpdatedAt.UTC())) {
		return ErrModelUserDateEarly
	}
	return nil
}
