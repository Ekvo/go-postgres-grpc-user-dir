package service

import (
	"context"
	"net"
	"testing"
	"time"

	user "github.com/Ekvo/go-grpc-apis/user/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/db/mock"
)

func newTest_UserRegisterRequest(t time.Time) *user.UserRegisterRequest {
	return &user.UserRegisterRequest{
		Login:     "ekvo",
		FirstName: "Alexander",
		Email:     "test@gmail.com",
		Password:  "qwert12345",
		CreatedAt: timestamppb.New(t),
	}
}

func Test_UserService(t *testing.T) {
	asserts := assert.New(t)
	requires := require.New(t)

	listener := bufconn.Listen(1024)
	srv := grpc.NewServer(grpc.UnaryInterceptor(Authorization))
	usecase := NewService(NewDepends(mock.NewMockProvider()))
	user.RegisterUserServiceServer(srv, usecase)

	go func() {
		if err := srv.Serve(listener); err != nil {
			requires.FailNowf("srv.Serve", "service_test: Serve error - {%v};", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	requires.NoError(err, "service_test: DialContext error - {%v};", err)

	client := user.NewUserServiceClient(conn)

	now := time.Now().UTC()
	res, err := client.UserRegister(context.Background(), newTest_UserRegisterRequest(now))
	asserts.NoError(err, "UserRegister error")

	asserts.Equal(uint64(1), res.UserId)
}
