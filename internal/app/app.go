package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	user "github.com/Ekvo/go-postgres-grpc-apis/user/v1"
	"google.golang.org/grpc"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/config"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/db"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/listen"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/service"
)

type App struct {
	service.Service
}

func NewApplication(usecase service.Service) *App {
	return &App{Service: usecase}
}

func Run(cfg *config.Config) {
	ctx := context.Background()

	dbProvider, err := db.OpenPool(ctx, cfg)
	if err != nil {
		log.Fatalf("app: db open error - %v", err)
	}
	defer dbProvider.ClosePool()

	app := NewApplication(
		service.NewService(
			service.NewOptions(cfg),
			service.NewDepends(dbProvider),
		))

	srv := grpc.NewServer()

	user.RegisterUserServiceServer(srv, app)

	go func() {
		listener, err := listen.NewListen(cfg)
		if err != nil {
			log.Fatalf("go app: net.Listen error - %v", err)
		}
		if err := srv.Serve(listener.Listener); err != nil {
			log.Fatalf("go app: server error - %v", err)
		}
		log.Print("go app: stopped serving\n")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	srv.GracefulStop()

	select {
	case <-ctx.Done():
		log.Print("app: shutdown took too long, forcing stop")
		srv.Stop()
	case <-time.After(10 * time.Second):
		log.Print("app: server stopped gracefully")
	}
}
