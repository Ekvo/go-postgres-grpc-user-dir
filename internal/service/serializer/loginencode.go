package serializer

import (
	"strconv"

	user "github.com/Ekvo/go-postgres-grpc-apis/user/v1"

	"github.com/Ekvo/go-postgres-grpc-user-dir/pkg/utils"
)

type LoginEncode struct {
	ID uint
}

func (le *LoginEncode) Response(secretKey string) (*user.LoginResponse, error) {
	token, err := utils.GenerateJWT(secretKey, strconv.FormatUint(uint64(le.ID), 10))
	return &user.LoginResponse{Token: token}, err
}
