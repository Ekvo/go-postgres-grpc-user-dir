package deserializer

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	user "github.com/Ekvo/go-postgres-grpc-apis/user/v1"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/model"
	"github.com/Ekvo/go-postgres-grpc-user-dir/pkg/utils"
)

type UserDecode struct {
	Login     string
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time

	user model.User
}

func NewUserDecode() *UserDecode {
	return &UserDecode{}
}

func (ud *UserDecode) Model() *model.User {
	return &ud.user
}

func (ud *UserDecode) Decode(req *user.SignUpRequest) error {
	ud.parseReq(req)
	if err := ud.validReq(); err != nil {
		return err
	}
	ud.setUser()
	return nil
}

func (ud *UserDecode) setUser() {
	ud.user.Login = ud.Login
	ud.user.FirstName = ud.FirstName
	ud.user.LastName = ud.LastName
	ud.user.Email = ud.Email
	ud.user.Password = ud.Password
	ud.user.CreatedAt = ud.CreatedAt
}

func (ud *UserDecode) parseReq(req *user.SignUpRequest) {
	ud.Login = req.GetLogin()
	ud.FirstName = req.GetFirstName()
	ud.LastName = req.GetLastName()
	ud.Email = req.GetEmail()
	ud.Password = req.GetPassword()
	ud.CreatedAt = req.GetCreatedAt().AsTime()
}

var reEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (ud *UserDecode) validReq() error {
	msgErr := utils.Message{}
	if strings.TrimSpace(ud.Login) == "" {
		msgErr["login"] = ErrDeserializerEmpty
	}
	if strings.TrimSpace(ud.FirstName) == "" {
		msgErr["first-name"] = ErrDeserializerEmpty
	}
	if !reEmail.MatchString(ud.Email) {
		msgErr["email"] = ErrDeserializerInvalid
	}
	if strings.TrimSpace(ud.Password) == "" {
		msgErr["password"] = ErrDeserializerEmpty
	}
	if ud.CreatedAt.IsZero() || ud.CreatedAt.UTC().After(time.Now().UTC()) {
		msgErr["created-at"] = ErrDeserializerInvalid
	}
	if len(msgErr) > 0 {
		return fmt.Errorf("deserializer: invalid signup - %s", msgErr.String())
	}
	return nil
}
