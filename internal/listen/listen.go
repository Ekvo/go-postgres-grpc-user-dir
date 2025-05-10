package listen

import (
	"fmt"
	"log"
	"net"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/config"
)

type Listen struct {
	net.Listener
}

func NewListen(cfg *config.Config) (*Listen, error) {
	listen, err := net.Listen(cfg.SRVNetwork, net.JoinHostPort("", cfg.SRVPort))
	if err != nil {
		return nil, fmt.Errorf("listen: net.Listen error - %w", err)
	}
	log.Printf("listen: server listen and serv in port - %s", cfg.SRVPort)
	return &Listen{Listener: listen}, nil
}
