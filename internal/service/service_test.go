package service

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"

	user "github.com/Ekvo/go-grpc-apis/user/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/db/mock"
)

type dataServer struct {
	lis    *bufconn.Listener
	srv    *grpc.Server
	client user.UserServiceClient
}

func (ds *dataServer) Close(t *testing.T) {
	t.Cleanup(func() {
		ds.srv.Stop()
		_ = ds.lis.Close()
	})
}

func newDataServer(t *testing.T) *dataServer {
	listener := bufconn.Listen(1024 * 0124)
	srv := grpc.NewServer(grpc.UnaryInterceptor(Authorization))
	usecase := NewService(NewDepends(mock.NewMockProvider()))
	user.RegisterUserServiceServer(srv, usecase)

	go func() {
		if err := srv.Serve(listener); err != nil {
			require.FailNowf(t, "srv.Serve", "service_test: Serve error - {%v};", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	require.NoErrorf(t, err, "service_test: DialContext error - {%v};", err)

	return &dataServer{
		lis:    listener,
		srv:    srv,
		client: user.NewUserServiceClient(conn),
	}
}

func Test_UserRegister_Service(t *testing.T) {
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

	dataService := newDataServer(t)
	defer dataService.Close(t)

	asserts := assert.New(t)
	requires := require.New(t)

	for i, test := range testData {
		log.Printf("\t%d - %s", i+1, test.title)

		res, err := dataService.client.UserRegister(context.Background(), test.data)

		if err == nil {
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
}
