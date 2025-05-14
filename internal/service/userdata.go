package service

import (
	"context"
	"log"

	user "github.com/Ekvo/go-grpc-apis/user/v1"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/service/deserializer"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/service/serializer"
)

// UserData - get user data from database
// get userID from ctx
// find user by ID from database
// create and return response
func (s *service) UserData(
	ctx context.Context, req *user.UserDataRequest) (*user.UserDataResponse, error) {
	deserialize := deserializer.NewIDDecode()
	if err := deserialize.Decode(ctx); err != nil {
		log.Printf("service: UserData Decode error - %v", err)
		return nil, ErrServiceInternal
	}

	u, err := s.DBProvider.FindUserByID(ctx, deserialize.UserID())
	if err != nil {
		log.Printf("service: UserData FindUserByID error - %v", err)
		return nil, ErrServiceNotFound
	}

	serialize := serializer.UserEncode{User: *u}

	return serialize.Response(), nil
}
