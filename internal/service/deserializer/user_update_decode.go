package deserializer

import (
	"fmt"
	"strings"
	"time"

	user "github.com/Ekvo/go-grpc-apis/user/v1"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/model"
	"github.com/Ekvo/go-postgres-grpc-user-dir/pkg/utils"
)

type UserUpdateDecode struct {
	Login     string
	FirstName string
	LastName  string
	Email     string
	Password  string
	UpdatedAt time.Time

	user model.User
}

func NewUserUpdateDecode() *UserUpdateDecode {
	return &UserUpdateDecode{}
}

func (uud *UserUpdateDecode) Model() *model.User {
	return &uud.user
}

func (uud *UserUpdateDecode) Decode(req *user.UserUpdateRequest) error {
	uud.parseReq(req)
	if err := uud.validReq(); err != nil {
		return err
	}
	uud.setUser()
	return nil
}

func (uud *UserUpdateDecode) setUser() {
	uud.user.Login = uud.Login
	uud.user.FirstName = uud.FirstName
	uud.user.LastName = uud.LastName
	uud.user.Email = uud.Email
	uud.user.Password = uud.Password
	uud.user.UpdatedAt = &uud.UpdatedAt
}

func (uud *UserUpdateDecode) parseReq(req *user.UserUpdateRequest) {
	uud.Login = req.GetLogin()
	uud.FirstName = req.GetFirstName()
	uud.LastName = req.GetLastName()
	uud.Email = req.GetEmail()
	uud.Password = req.GetPassword()
	uud.UpdatedAt = req.GetUpdatedAt().AsTime()
}

func (uud *UserUpdateDecode) validReq() error {
	msgErr := utils.Message{}
	if uud.Email = strings.TrimSpace(uud.Email); uud.Email != "" && !reEmail.MatchString(uud.Email) {
		msgErr["email"] = ErrDeserializerInvalid
	}
	if uud.UpdatedAt.IsZero() {
		msgErr["updated-at"] = ErrDeserializerInvalid
	}
	if len(msgErr) > 0 {
		return fmt.Errorf("deserializer: invalid user update - %s", msgErr.String())
	}
	uud.Login = strings.TrimSpace(uud.Login)
	uud.FirstName = strings.TrimSpace(uud.FirstName)
	uud.LastName = strings.TrimSpace(uud.LastName)
	uud.Password = strings.TrimSpace(uud.Password)
	return nil
}
