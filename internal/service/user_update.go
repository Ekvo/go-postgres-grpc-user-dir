package service

import (
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"

	user "github.com/Ekvo/go-grpc-apis/user/v1"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/model"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/service/deserializer"
)

// UserUpdate - rules for update User data
// decode the new user data from request
// decode user ID from ctx
// call userUpdate
func (s *service) UserUpdate(
	ctx context.Context,
	req *user.UserUpdateRequest) (*user.UserUpdateResponse, error) {
	deserializeUserData := deserializer.NewUserUpdateDecode()
	if err := deserializeUserData.Decode(req); err != nil {
		return nil, err
	}

	deserializeUserID := deserializer.NewIDDecode()
	if err := deserializeUserID.Decode(ctx); err != nil {
		log.Printf("service: UserUpdate Decode error - {%v};", err)
		return nil, ErrServiceInternal
	}

	userNewData := deserializeUserData.Model()
	userNewData.ID = deserializeUserID.UserID()

	if err := s.userUpdate(ctx, userNewData); err != nil {
		return nil, err
	}
	return &user.UserUpdateResponse{}, nil
}

// userUpdate - prepares and writes user data to the database
// gets user data from the database by ID
// new password is not empty -> create hashedPassword
// creates a user for writing to the database (NewData) -> (internal/model/user.go)
// updates user data in the storage
func (s *service) userUpdate(ctx context.Context, userNewData *model.User) error {
	u, err := s.DBProvider.FindUserByID(ctx, userNewData.ID)
	if err != nil {
		log.Printf("service: userUpdate FindUserByID error - {%v};", err)
		return ErrServiceNotFound
	}

	if userNewData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userNewData.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("service: userUpdate GenerateFromPassword - error {%v}, password - {%s};", err, userNewData.Password)
			return ErrServiceInternal
		}
		userNewData.Password = string(hashedPassword)
	}

	if err := u.NewData(userNewData); err != nil {
		log.Printf("service: userUpdate NewData error - {%v};", err)
		return ErrServiceUpdateDataInvalid
	}

	if err := s.DBProvider.UpdateUser(ctx, u); err != nil {
		log.Printf("service: userUpdate UpdateUser error - {%v};", err)
		return ErrServiceInternal
	}

	return nil
}
