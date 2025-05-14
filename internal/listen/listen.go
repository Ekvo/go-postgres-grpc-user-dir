package listen

import (
	"fmt"
	"log"
	"net"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/config"
)

func NewListen(cfg *config.Config) (net.Listener, error) {
	listen, err := net.Listen(cfg.SRVNetwork, net.JoinHostPort("", cfg.SRVPort))
	if err != nil {
		return nil, fmt.Errorf("listen: net.Listener error - %w", err)
	}
	log.Printf("listen: server listen and serv in port - %s", cfg.SRVPort)
	return listen, nil
}
