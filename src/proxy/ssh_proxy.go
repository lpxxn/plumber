package proxy

import (
	"fmt"
	"net"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/common"
	"github.com/lpxxn/plumber/src/log"
)

type SSHProxy struct {
	// if the client is ssh proxy client, this field will be set
	LocalListener    net.Listener
	SSHConfig        *config.SSHConf
	RemoteTunnelConn net.Conn
}

func (s *SSHProxy) Close() error {
	if s.LocalListener != nil {
		s.LocalListener.Close()
		s.LocalListener = nil
	}
	log.Infof("SSHProxy: closed")
	return nil
}

func (s *SSHProxy) NewTCPServer() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.SSHConfig.SrvPort))
	if err != nil {
		log.Errorf("SSHProxy listen on %d failed: %v", s.SSHConfig.SrvPort, err)
		return err
	}
	s.LocalListener = listener
	log.Infof("SSHProxy: listening on %s", s.LocalListener.Addr())
	return nil
}

func (s *SSHProxy) WaitForTunnelConn() error {
	log.Infof("SSHProxy: waiting for tunnel connection")
	conn, err := s.LocalListener.Accept()
	if err != nil {
		log.Errorf("SSHProxy accept failed: %v", err)
		return err
	}
	if err := common.VerifySSHMagicStrConnection(conn); err != nil {
		conn.Close()
		return err
	}
	s.RemoteTunnelConn = conn
	log.Infof("SSHProxy: got tunnel connection")
	return nil
}

func (s *SSHProxy) Start(handler TCPHandler) error {
	go func() {
		if err := TCPServer(s.LocalListener, handler); err != nil {
			log.Errorf("SSHProxy TCPServer failed: %v", err)
		}
	}()
	return nil
}
