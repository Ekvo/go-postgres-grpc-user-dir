// create net.Listener for server
package listen

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/config"
)

func NewListen(cfg *config.ServerConfig) (net.Listener, error) {
	srvPort := strconv.FormatUint(uint64(cfg.Port), 10)
	listen, err := net.Listen(cfg.Network, net.JoinHostPort("", srvPort))
	if err != nil {
		return nil, fmt.Errorf("listen: net.Listener error - {%w};", err)
	}

	log.Printf("listen: server listen and serv in port - {%s};", srvPort)

	return listen, nil
}
