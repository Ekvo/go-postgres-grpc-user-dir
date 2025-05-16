// create jwt.Toke for Response
package serializer

import (
	"strconv"

	user "github.com/Ekvo/go-grpc-apis/user/v1"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/lib/jwtsign"
)

type LoginEncode struct {
	ID uint
}

func (le *LoginEncode) Response() (*user.UserLoginResponse, error) {
	content := jwtsign.Content{}
	content["user_id"] = strconv.FormatUint(uint64(le.ID), 10)
	token, err := jwtsign.TokenGenerator(content)
	return &user.UserLoginResponse{Token: token}, err
}
