package service

import (
	"context"

	user "github.com/Ekvo/go-postgres-grpc-apis/user/v1"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/service/deserializer"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/service/serializer"
)

func (s *service) Login(
	ctx context.Context,
	req *user.LoginRequest) (*user.LoginResponse, error) {
	deserialize := deserializer.NewLoginDecode()
	if err := deserialize.Decode(req); err != nil {
		return nil, err
	}

	login := deserialize.Model()
	u, err := s.DBProvider.FindUserByEmail(ctx, login.Email)
	if err != nil {
		return nil, err
	}
	if err := u.ValidPassword(login.Password); err != nil {
		return nil, err
	}

	serialize := serializer.LoginEncode{ID: u.ID}
	return serialize.Response(s.JWTSecret)
}
