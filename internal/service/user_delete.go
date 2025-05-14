package service

import (
	"context"
	"log"

	user "github.com/Ekvo/go-grpc-apis/user/v1"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/service/deserializer"
)

// UserDelete - rules for delete User
// decode the user ID from the ctx
// remove user by ID from database
func (s *service) UserDelete(
	ctx context.Context,
	_ *user.UserDeleteRequest) (*user.UserDeleteResponse, error) {
	deserialize := deserializer.NewIDDecode()
	if err := deserialize.Decode(ctx); err != nil {
		log.Printf("service: UserDelete Decode error - {%v};", err)
		return nil, ErrServiceInternal
	}

	if err := s.DBProvider.RemoveUserByID(ctx, deserialize.UserID()); err != nil {
		log.Printf("service: UserDelete RemoveUserByID error - {%v};", err)
		return nil, ErrServiceNotFound
	}

	return &user.UserDeleteResponse{}, nil
}
