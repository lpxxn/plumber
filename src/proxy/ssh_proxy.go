package proxy

import (
	"fmt"
	"net"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/log"
)

type SSHProxy struct {
	// if the client is ssh proxy client, this field will be set
	LocalListener net.Listener
	SSHConfig     *config.SSHConf
}

func (s *SSHProxy) Close() error {
	if s.LocalListener != nil {
		s.LocalListener.Close()
		s.LocalListener = nil
	}
	return nil
}

func (s *SSHProxy) NewTCPServer(handler TCPHandler) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.SSHConfig.SrvPort))
	if err != nil {
		log.Errorf("SSHProxy listen on %d failed: %v", s.SSHConfig.SrvPort, err)
		return err
	}
	return TCPServer(listener, handler)
}
