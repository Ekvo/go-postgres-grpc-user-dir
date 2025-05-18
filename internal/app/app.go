// initialization of application start and stop
package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	user "github.com/Ekvo/go-grpc-apis/user/v1"
	"google.golang.org/grpc"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/config"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/db"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/lib/jwtsign"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/listen"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/service"
)

// Application - contains a server, server listener,business service, an interface for working with the store
type Application struct {
	userRepository db.Provider
	userService    service.Service
	srv            *grpc.Server
	listener       net.Listener
}

// NewApplication
// create: secretKey for jwt, net.Listener for server, open pgx.pool, service.NewService, grpc.NewServer
// save all main variables inside &Application{}
func NewApplication(cfg *config.Config) (*Application, error) {
	if err := jwtsign.NewSecretKey(cfg); err != nil {
		return nil, fmt.Errorf("app: jwt error - %w", err)
	}

	listener, err := listen.NewListen(cfg)
	if err != nil {
		return nil, fmt.Errorf("app: Listener error - %w", err)
	}

	dbProvider, err := db.OpenPool(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("app: db error - %w", err)
	}

	app := &Application{}
	app.userRepository = dbProvider
	app.userService = service.NewService(service.NewDepends(dbProvider))
	app.srv = grpc.NewServer(grpc.UnaryInterceptor(service.Authorization))
	app.listener = listener

	return app, nil
}

// Run - registers server then start server inside go func()
func (a *Application) Run() {
	user.RegisterUserServiceServer(a.srv, a.userService)

	go func() {
		log.Print("go app: start server\n")
		if err := a.srv.Serve(a.listener); err != nil {
			a.userRepository.ClosePool()
			log.Fatalf("go app: server error - %v", err)
		}
		log.Print("go app: stopped serving\n")
	}()
}

// Stop - close pgx.pool, call GracefulStop() with select {<- ctx, time.After}
func (a *Application) Stop() {
	gracefully := true
	timer := time.AfterFunc(10*time.Second, func() {
		gracefully = false
		log.Print("app: server stopped - forcing stop\n")
		a.srv.Stop()
	})
	defer func() {
		timer.Stop()
		_ = a.listener.Close()
		a.userRepository.ClosePool()
	}()

	a.srv.GracefulStop()
	if gracefully {
		log.Print("app: server stopped - gracefully\n")
	}

}
