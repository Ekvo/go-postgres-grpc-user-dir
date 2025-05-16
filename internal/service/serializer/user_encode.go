// create user data for Response
package serializer

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	user "github.com/Ekvo/go-grpc-apis/user/v1"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/model"
)

type UserEncode struct {
	model.User
}

func (ue *UserEncode) Response() *user.UserDataResponse {
	userResponse := &user.User{
		Id:        uint64(ue.ID),
		Login:     ue.Login,
		FirstName: ue.FirstName,
		LastName:  ue.LastName,
		Email:     ue.Email,
		CreatedAt: timestamppb.New(ue.CreatedAt),
	}
	if updateTime := ue.UpdatedAt; updateTime != nil && !updateTime.IsZero() {
		userResponse.UpdatedAt = timestamppb.New(*updateTime)
	}
	return &user.UserDataResponse{User: userResponse}
}
