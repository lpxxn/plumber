package proxy

import (
	"github.com/lpxxn/plumber/src/log"
	"io"
	"net"
)

type SSHProxy struct {
	SrvAddr      string
	listener     net.Listener
	LocalSSHAddr string
}

func NewSSHProxy(srvAddr, localSSHAddr string) *SSHProxy {
	return &SSHProxy{
		SrvAddr:      srvAddr,
		LocalSSHAddr: localSSHAddr,
	}
}

func (s *SSHProxy) Start() {
	var err error
	s.listener, err = net.Listen("tcp", s.SrvAddr)
	if err != nil {
		log.Panicf("listen error: %v addr: %s", err, s.SrvAddr)
	}
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Panicf("Local Tcp accept error: %v", err)
			return
		}
		go s.handleConnection(conn)
	}
}

func (s *SSHProxy) handleConnection(remoteConn net.Conn) {
	defer remoteConn.Close()

	localConn, err := net.Dial("tcp", s.LocalSSHAddr)
	if err != nil {
		return
	}
	defer localConn.Close()
	if err := proxyConn(remoteConn, localConn); err != nil {
		// check err if connection is closed
		if err == io.EOF {
			log.Infof("connection closed")
		} else {
			log.Infof("proxy connection error: %v", err)
		}
	}
}

func proxyConn(localConn, remoteConn net.Conn) error {
	errCh := make(chan error, 1)
	// copy from local to remote
	go func() {
		_, err := io.Copy(remoteConn, localConn)
		if err != nil {
			log.Infof("copy from local to remote error: %v", err)
		}
		errCh <- err
	}()
	// copy from remote to local
	go func() {
		_, err := io.Copy(localConn, remoteConn)
		if err != nil {
			log.Infof("copy from remote to local error: %v", err)
		}
		errCh <- err
	}()
	return <-errCh
}
