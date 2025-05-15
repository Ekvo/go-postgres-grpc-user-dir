// rulse for user registration
package deserializer

import (
	"fmt"
	"strings"
	"time"

	user "github.com/Ekvo/go-grpc-apis/user/v1"

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

func (ud *UserDecode) Decode(req *user.UserRegisterRequest) error {
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

func (ud *UserDecode) parseReq(req *user.UserRegisterRequest) {
	ud.Login = req.GetLogin()
	ud.FirstName = req.GetFirstName()
	ud.LastName = req.GetLastName()
	ud.Email = req.GetEmail()
	ud.Password = req.GetPassword()
	ud.CreatedAt = req.GetCreatedAt().AsTime()
}

// validReq - check critical fields for new user registration
func (ud *UserDecode) validReq() error {
	msgErr := utils.Message{}
	if ud.Login = strings.TrimSpace(ud.Login); ud.Login == "" {
		msgErr["login"] = ErrDeserializerEmpty
	}
	if ud.FirstName = strings.TrimSpace(ud.FirstName); ud.FirstName == "" {
		msgErr["first-name"] = ErrDeserializerEmpty
	}
	if ud.Email = strings.TrimSpace(ud.Email); !reEmail.MatchString(ud.Email) {
		msgErr["email"] = ErrDeserializerInvalid
	}
	if ud.Password = strings.TrimSpace(ud.Password); ud.Password == "" {
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
