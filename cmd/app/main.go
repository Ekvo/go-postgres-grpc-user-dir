package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/app"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/config"
)

func main() {
	cfg, err := config.NewConfig("./init/.env")
	if err != nil {
		log.Fatalf("main: config error - %v", err)
	}

	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatalf("main: application error - %v", err)
	}

	application.Run()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	application.Stop()
}
