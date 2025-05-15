package service

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"
	"time"

	user "github.com/Ekvo/go-grpc-apis/user/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/config"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/db/mock"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/lib/jwtsign"
)

type dataServer struct {
	lis    *bufconn.Listener
	srv    *grpc.Server
	client user.UserServiceClient
}

// newDataServer - implemet and start server
func newDataServer() (*dataServer, error) {
	_ = jwtsign.NewSecretKey(&config.Config{JWTSecretKey: "secret"})

	listener := bufconn.Listen(1024 * 0124)
	srv := grpc.NewServer(grpc.UnaryInterceptor(Authorization))
	usecase := NewService(NewDepends(mock.NewMockProvider()))
	user.RegisterUserServiceServer(srv, usecase)

	go func() {
		if err := srv.Serve(listener); err != nil {
			_ = listener.Close()
			log.Fatalf("go service_test: newDataServer Serve error - {%v};", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &dataServer{
		lis:    listener,
		srv:    srv,
		client: user.NewUserServiceClient(conn),
	}, nil
}

// createDataFroAutirizationWithContext - create and load start date to store (mock)
// get Bearer token from clien, set token to context
func (ds *dataServer) createDataFroAutirizationWithContext(date time.Time) (context.Context, error) {
	ctx := context.Background()

	if _, err := ds.client.UserRegister(ctx, newUserRegisterRequest(date)); err != nil {
		return nil, err
	}

	token, err := ds.client.UserLogin(ctx, newUserLoginRequest())
	if err != nil {
		return nil, err
	}

	md := metadata.Pairs("authorization", "bearer "+token.Token)

	return metadata.NewOutgoingContext(ctx, md), nil
}

func (ds *dataServer) close(t *testing.T) {
	t.Cleanup(func() {
		ds.srv.Stop()
		_ = ds.lis.Close()
	})
}

func newUserRegisterRequest(date time.Time) *user.UserRegisterRequest {
	return &user.UserRegisterRequest{
		Login:     `avp`,
		FirstName: `NameTest`,
		Email:     `test@example.com`,
		Password:  `testpassword`,
		CreatedAt: timestamppb.New(date),
	}
}

func newUserLoginRequest() *user.UserLoginRequest {
	return &user.UserLoginRequest{
		Email:    `test@example.com`,
		Password: `testpassword`,
	}
}

func Test_UserRegister_Service(t *testing.T) {
	log.Printf("service_test: Test_UserRegister_Service - START")

	asserts := assert.New(t)
	requires := require.New(t)

	var testData = []struct {
		title       string
		data        *user.UserRegisterRequest
		expectedRes uint64
		expectedErr error
		msg         string
	}{
		{
			title: `valid register`,
			data: &user.UserRegisterRequest{
				Login:     `fury`,
				FirstName: `Mary`,
				Email:     `fury1995@example.com`,
				Password:  `Maryshy`,
				CreatedAt: timestamppb.Now(),
			},
			expectedRes: 1,
			expectedErr: nil,
			msg:         `this test should be valid, return a new id for the new user and the error should be zero`,
		},
		{
			title: `wrong register,login alrady exists`,
			data: &user.UserRegisterRequest{
				Login:     `fury`,
				FirstName: `Mary`,
				Email:     `fury2000@example.com`,
				Password:  `Maryshy`,
				CreatedAt: timestamppb.Now(),
			},
			expectedRes: 0,
			expectedErr: ErrServiceAlreadyExists,
			msg:         `wrong result, login already exists return error`,
		},
		{
			title: `wrong register, email alrady exists`,
			data: &user.UserRegisterRequest{
				Login:     `MaryFox`,
				FirstName: `Mary`,
				Email:     `fury1995@example.com`,
				Password:  `goldenfox`,
				CreatedAt: timestamppb.Now(),
			},
			expectedRes: 0,
			expectedErr: ErrServiceAlreadyExists,
			msg:         `wrong result, email already exists return error`,
		},
		{
			title: `invalid register, no login`,
			data: &user.UserRegisterRequest{
				Login:     `serg`,
				FirstName: `Serge`,
				Email:     `sergeUexample.com`,
				Password:  `Ultra`,
				CreatedAt: timestamppb.Now(),
			},
			expectedRes: 0,
			expectedErr: errors.New(`deserializer: invalid signup - {email:invalid}`),
			msg:         `invalid email, error is exist`,
		},
		{
			title:       `invalid register, no login`,
			data:        &user.UserRegisterRequest{},
			expectedRes: 0,
			expectedErr: errors.New(`deserializer: invalid signup - {email:invalid},{first-name:empty},{login:empty},{password:empty}`),
			msg:         `invalid user data, error is exist`,
		},
	}

	dataService, err := newDataServer()
	if err != nil {
		log.Fatalf("service_test: Test_UserRegister_Service - {%v};", err)
	}
	defer dataService.close(t)

	for i, test := range testData {
		log.Printf("\t%d - %s", i+1, test.title)

		res, err := dataService.client.UserRegister(context.Background(), test.data)

		if err == nil {
			asserts.NoError(test.expectedErr, test.msg)
			asserts.Equal(test.expectedRes, res.UserId, test.msg)
		} else {
			st, ok := status.FromError(err)
			if !ok {
				requires.FailNow("this is not a GRPC error it is ALIEN")
			}
			requires.NotNil(test.expectedErr, test.msg)
			asserts.Equal(test.expectedErr.Error(), st.Message(), test.msg)
			asserts.Zero(test.expectedRes, test.msg)
		}
	}
	log.Printf("service_test: Test_UserRegister_Service - END")
}

func Test_UserLogin_Service(t *testing.T) {
	log.Printf("service_test: Test_UserLogin_Service - START")

	asserts := assert.New(t)
	requires := require.New(t)

	var testData = []struct {
		title       string
		data        *user.UserLoginRequest
		expectedRes string
		expectedErr error
		msg         string
	}{
		{
			title:       `valid login`,
			data:        newUserLoginRequest(),
			expectedRes: `.+`,
			expectedErr: nil,
			msg:         `token must be exist, error is nil`,
		},
		{
			title: `wrong login, email not faound`,
			data: &user.UserLoginRequest{
				Email:    `soem@example.com`,
				Password: `testpassword`,
			},
			expectedRes: ``,
			expectedErr: ErrServiceNotFound,
			msg:         `token must be exist, error is nil`,
		},
		{
			title: `valid login`,
			data: &user.UserLoginRequest{
				Email:    `test@example.com`,
				Password: `somepassword`,
			},
			expectedRes: ``,
			expectedErr: ErrServicePasswordInvalid,
			msg:         `token must be exist, error is nil`,
		},
		{
			title:       `valid login`,
			data:        &user.UserLoginRequest{},
			expectedRes: ``,
			expectedErr: errors.New(`deserializer: invalid login - {email:invalid},{password:empty}`),
			msg:         `token must be exist, error is nil`,
		},
	}

	dataService, err := newDataServer()
	if err != nil {
		log.Printf("service_test: Test_UserLogin_Service newDataServer error - {%v};", err)
		return
	}
	defer dataService.close(t)

	_, err = dataService.client.UserRegister(context.Background(), newUserRegisterRequest(time.Now().UTC()))
	if err != nil {
		log.Printf("service_test: Test_UserLogin_Service UserRegister error - {%v};", err)
		return
	}

	for i, test := range testData {
		log.Printf("\t%d - %s", i+1, test.title)

		res, err := dataService.client.UserLogin(context.Background(), test.data)

		if err == nil {
			asserts.NoError(test.expectedErr, test.msg)
			asserts.Regexp(test.expectedRes, res.Token, test.msg)
		} else {
			st, ok := status.FromError(err)
			if !ok {
				requires.FailNow("this is not a GRPC error it is ALIEN")
			}
			requires.NotNil(test.expectedErr, test.msg)
			asserts.Equal(test.expectedErr.Error(), st.Message(), test.msg)
			asserts.Zero(test.expectedRes, test.msg)
		}
	}
	log.Printf("service_test: Test_UserLogin_Service - END")
}

func Test_UserData_Service(t *testing.T) {
	log.Printf("service_test: Test_UserData_Service - START")

	asserts := assert.New(t)
	requires := require.New(t)

	dataService, err := newDataServer()
	if err != nil {
		log.Printf("service_test: Test_UserData_Service newDataServer error - {%v};", err)
		return
	}
	defer dataService.close(t)

	now := time.Now().UTC()

	ctx, err := dataService.createDataFroAutirizationWithContext(now)
	if err != nil {
		log.Printf("service_test: Test_UserData_Service createDataFroAutirizationWithContext error - {%v};", err)
		return
	}
	log.Printf("service_test: Test_UserData_Service - valid test")

	res, err := dataService.client.UserData(ctx, &user.UserDataRequest{})
	requires.NoError(err, "get UserData incorrectly")

	userRes := res.User
	expectedUser := user.UserDataResponse{
		User: &user.User{
			Id:        userRes.Id,
			Login:     `avp`,
			FirstName: `NameTest`,
			Email:     `test@example.com`,
			CreatedAt: timestamppb.New(now),
		},
	}.User
	asserts.Equal(expectedUser, userRes, "user data not equal")

	log.Printf("service_test: Test_UserData_Service - wrong test")

	res, err = dataService.client.UserData(context.Background(), &user.UserDataRequest{})
	requires.NotNil(err, "shoud be not nil")
	st, ok := status.FromError(err)
	if !ok {
		requires.FailNow("this is not a GRPC error it is ALIEN")
	}
	asserts.Equal(ErrServiceAuthorizationInvalid.Error(), st.Message(), "differen errors")
	asserts.Nil(res, "should be nil")

	log.Printf("service_test: Test_UserData_Service - END")
}

func Test_UserUpdate_Service(t *testing.T) {
	log.Printf("service_test: Test_UserUpdate_Service - START")

	asserts := assert.New(t)
	requires := require.New(t)

	createTime := time.Now().UTC()
	updateTime := createTime.Add(time.Hour)

	var testData = []struct {
		title       string
		data        *user.UserUpdateRequest
		expectedErr error
		msg         string
	}{
		{
			title: `valid update user data`,
			data: &user.UserUpdateRequest{
				Login:     `fury`,
				FirstName: `Mary`,
				Email:     `fury1995@example.com`,
				Password:  `Marystrong`,
				UpdatedAt: timestamppb.New(updateTime),
			},
			expectedErr: nil,
			msg:         `this test should be valid, return empty result (not nil), error - nil`,
		},
		{
			title: `invalid update (important)`,
			data: &user.UserUpdateRequest{
				Login:     `fury`,
				FirstName: `Mary`,
				Email:     `fury1995@example.com`,
				Password:  `Marystrong`,
				UpdatedAt: timestamppb.New(updateTime),
			},
			expectedErr: ErrServiceUpdateDataInvalid,
			msg:         `valid, time for update not after than previous update, error`,
		},
		{
			title: `valid update user data, empty password `,
			data: &user.UserUpdateRequest{
				Login:     `fury`,
				FirstName: `Mary`,
				Email:     `fury1995@example.com`,
				UpdatedAt: timestamppb.New(updateTime.Add(time.Hour)),
			},
			expectedErr: nil,
			msg:         `this test should be valid, return empty result (not nil), error - nil`,
		},
		{
			title:       `wrong update user data`,
			data:        &user.UserUpdateRequest{},
			expectedErr: errors.New(`deserializer: invalid user update - {email:invalid},{first-name:empty},{login:empty}`),
			msg:         `empty request, error is  exist`,
		},
	}

	dataService, err := newDataServer()
	if err != nil {
		log.Printf("service_test: Test_UserData_Service newDataServer error - {%v};", err)
		return
	}
	defer dataService.close(t)

	ctx, err := dataService.createDataFroAutirizationWithContext(createTime)
	if err != nil {
		log.Printf("service_test: Test_UserUpdate_Service createDataFroAutirizationWithContext error - {%v};", err)
		return
	}

	for i, test := range testData {
		log.Printf("\t%d - %s", i+1, test.title)

		res, err := dataService.client.UserUpdate(ctx, test.data)

		if err == nil {
			requires.NotNil(res, "result should be not nil")
		} else {
			st, ok := status.FromError(err)
			if !ok {
				requires.FailNow("this is not a GRPC error it is ALIEN")
			}
			requires.NotNil(test.expectedErr, test.msg)
			asserts.Equal(test.expectedErr.Error(), st.Message(), test.msg)
		}
	}
	log.Printf("service_test: Test_UserUpdate_Service - END")
}

func Test_UserDelete_Service(t *testing.T) {
	log.Printf("service_test: Test_UserDelete_Service - START")

	asserts := assert.New(t)
	requires := require.New(t)

	dataService, err := newDataServer()
	if err != nil {
		log.Printf("service_test: Test_UserDelete_Service newDataServer error - {%v};", err)
		return
	}
	defer dataService.close(t)

	now := time.Now().UTC()

	ctx, err := dataService.createDataFroAutirizationWithContext(now)
	if err != nil {
		log.Printf("service_test: Test_UserDelete_Service createDataFroAutirizationWithContext error - {%v};", err)
		return
	}
	log.Printf("service_test: Test_UserDelete_Service - valid test")

	res, err := dataService.client.UserDelete(ctx, &user.UserDeleteRequest{})
	asserts.NoError(err, "correct delete from store")
	asserts.NotNil(res, "shouldn't be nil")

	log.Printf("service_test: Test_UserDelete_Service - wrong test")

	res, err = dataService.client.UserDelete(ctx, &user.UserDeleteRequest{})
	requires.Error(err, "error shoud be not nil")
	st, ok := status.FromError(err)
	if !ok {
		requires.FailNow("this is not a GRPC error it is ALIEN")
	}
	asserts.Equal(ErrServiceNotFound.Error(), st.Message(), "differen errors")
	asserts.Nil(res, "should be nil")

	log.Printf("service_test: Test_UserDelete_Service - END")
}
