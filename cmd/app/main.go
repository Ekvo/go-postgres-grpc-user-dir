package main

import (
	"log"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/app"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/config"
)

func main() {
	cfg, err := config.NewConfig("./init/.env")
	if err != nil {
		log.Fatalf("main: config error - %v", err)
	}
	app.Run(cfg)
}
