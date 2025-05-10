package deserializer

import (
	"fmt"
	"strings"

	user "github.com/Ekvo/go-postgres-grpc-apis/user/v1"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/model"
	"github.com/Ekvo/go-postgres-grpc-user-dir/pkg/utils"
)

type LoginDecode struct {
	Email    string
	Password string

	login model.Login
}

func NewLoginDecode() *LoginDecode {
	return &LoginDecode{}
}

func (ul *LoginDecode) Model() *model.Login {
	return &ul.login
}

func (ul *LoginDecode) Decode(req *user.LoginRequest) error {
	ul.parseReq(req)
	if err := ul.validReq(); err != nil {
		return err
	}
	ul.setUser()
	return nil
}

func (ul *LoginDecode) setUser() {
	ul.login.Email = ul.Email
	ul.login.Password = ul.Password
}

func (ul *LoginDecode) parseReq(req *user.LoginRequest) {
	ul.Email = req.GetEmail()
	ul.Password = req.GetPassword()
}

func (ul *LoginDecode) validReq() error {
	msgErr := utils.Message{}
	if !reEmail.MatchString(ul.Email) {
		msgErr["email"] = ErrDeserializerInvalid
	}
	if strings.TrimSpace(ul.Password) == "" {
		msgErr["password"] = ErrDeserializerEmpty
	}
	if len(msgErr) > 0 {
		return fmt.Errorf("deserializer: invalid signup - %s", msgErr.String())
	}
	return nil
}
