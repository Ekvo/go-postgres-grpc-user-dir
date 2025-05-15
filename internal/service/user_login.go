package service

import (
	"context"
	"log"

	user "github.com/Ekvo/go-grpc-apis/user/v1"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/service/deserializer"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/service/serializer"
)

// UserLogin - rules for entering the User Service
// decode user from request
// find user by email in database, then check password
// create bearer token for response
func (s *service) UserLogin(
	ctx context.Context,
	req *user.UserLoginRequest) (*user.UserLoginResponse, error) {
	deserialize := deserializer.NewLoginDecode()
	if err := deserialize.Decode(req); err != nil {
		return nil, err
	}

	login := deserialize.Model()
	u, err := s.DBProvider.FindUserByEmail(ctx, login.Email)
	if err != nil {
		log.Printf("service: UserLogin FindUserByEmail error - {%v};", err)
		return nil, ErrServiceNotFound
	}
	if err := u.ValidPassword(login.Password); err != nil {
		log.Printf("service: UserLogin ValidPassword error - {%v};", err)
		return nil, ErrServicePasswordInvalid
	}

	serialize := serializer.LoginEncode{ID: u.ID}
	userLoginResponse, err := serialize.Response()
	if err != nil {
		log.Printf("service: UserLogin LoginEncode error - {%v};", err)
		return nil, ErrServiceInternal
	}

	return userLoginResponse, nil
}
