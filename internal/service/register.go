package service

import (
	"context"

	user "github.com/Ekvo/go-postgres-grpc-apis/user/v1"
	"golang.org/x/crypto/bcrypt"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/service/deserializer"
)

func (s *service) SignUp(
	ctx context.Context,
	req *user.SignUpRequest) (*user.SignUpResponse, error) {
	deserialize := deserializer.NewUserDecode()
	if err := deserialize.Decode(req); err != nil {
		return nil, err
	}

	u := deserialize.Model()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.Password = string(hashedPassword)

	id, err := s.DBProvider.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	return &user.SignUpResponse{UserId: uint64(id)}, nil
}
