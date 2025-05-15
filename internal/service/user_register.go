package service

import (
	"context"
	"log"

	user "github.com/Ekvo/go-grpc-apis/user/v1"

	"golang.org/x/crypto/bcrypt"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/service/deserializer"
)

// UserRegister - rules for creating a new user in User Srvice
// decode the user from the request
// create a hashed password for the user
// write the user to the database
// return the new user ID
func (s *service) UserRegister(
	ctx context.Context,
	req *user.UserRegisterRequest) (*user.UserRegisterResponse, error) {
	deserialize := deserializer.NewUserDecode()
	if err := deserialize.Decode(req); err != nil {
		return nil, err
	}

	u := deserialize.Model()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("service: UserRegister GenerateFromPassword - error {%v};", err)
		return nil, ErrServiceInternal
	}
	u.Password = string(hashedPassword)

	id, err := s.DBProvider.CreateUser(ctx, u)
	if err != nil {
		log.Printf("service: UserRegister CreateUser - error {%v};", err)
		return nil, ErrServiceAlreadyExists
	}

	return &user.UserRegisterResponse{UserId: uint64(id)}, nil
}
