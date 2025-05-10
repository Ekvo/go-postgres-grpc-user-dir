package service

import (
	"context"
	"strconv"

	user "github.com/Ekvo/go-postgres-grpc-apis/user/v1"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/service/deserializer"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/service/serializer"
	"github.com/Ekvo/go-postgres-grpc-user-dir/pkg/utils"
)

func (s *service) UserData(
	ctx context.Context, req *user.UserDataRequest) (*user.UserDataResponse, error) {
	deserialize := deserializer.NewTokenDecode()
	if err := deserialize.Decode(ctx); err != nil {
		return nil, err
	}

	idStr, err := utils.VerifyJWT(s.JWTSecret, deserialize.Token())
	if err != nil {
		return nil, err
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	u, err := s.DBProvider.FindUserByID(ctx, uint(id))
	if err != nil {
		return nil, err
	}

	serialize := serializer.UserEncode{User: *u}
	return serialize.Response(), nil
}
