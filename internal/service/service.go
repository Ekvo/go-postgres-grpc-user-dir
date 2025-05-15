package service

import (
	"errors"

	user "github.com/Ekvo/go-grpc-apis/user/v1"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/db"
)

// errors for response
var (
	ErrServiceInternal = errors.New("internal error")

	ErrServiceNotFound = errors.New("not found")

	ErrServiceAlreadyExists = errors.New("already exists")

	ErrServiceAuthorizationInvalid = errors.New("invalid authorization")

	ErrServicePasswordInvalid = errors.New("invalid password")

	ErrServiceUpdateDataInvalid = errors.New("invalid update data")
)

type Service interface {
	user.UserServiceServer
}

// Depends- if necessary add another base
type Depends struct {
	DBProvider db.Provider
}

func NewDepends(dbProvider db.Provider) Depends {
	return Depends{DBProvider: dbProvider}
}

type service struct {
	Depends
}

func NewService(dep Depends) *service {
	return &service{Depends: dep}
}
