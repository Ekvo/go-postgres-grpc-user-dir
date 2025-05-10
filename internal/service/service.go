package service

import (
	user "github.com/Ekvo/go-postgres-grpc-apis/user/v1"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/config"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/db"
)

type Service interface {
	user.UserServiceServer
}

type Options struct {
	JWTSecret string
}

func NewOptions(cfg *config.Config) Options {
	return Options{JWTSecret: cfg.JWTSecretKey}
}

type Depends struct {
	DBProvider db.Provider
}

func NewDepends(dbProvider db.Provider) Depends {
	return Depends{DBProvider: dbProvider}
}

type service struct {
	Options
	Depends
}

func NewService(opt Options, dep Depends) *service {
	return &service{Options: opt, Depends: dep}
}
