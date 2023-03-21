package internal

import (
	"io"
	"log"
	"net"
)

type SSHProxy struct {
	SrvAddr   string
	listener  net.Listener
	LocalAddr string
}

func NewSSHProxy(srvAddr, localAddr string) *SSHProxy {
	return &SSHProxy{
		SrvAddr:   srvAddr,
		LocalAddr: localAddr,
	}
}

func (s *SSHProxy) Start() {
	var err error
	s.listener, err = net.Listen("tcp", s.LocalAddr)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}
		go s.handleConnection(conn)
	}
}

func (s *SSHProxy) handleConnection(conn net.Conn) {
	defer conn.Close()

	localConn, err := net.Dial("tcp", s.SrvAddr)
	if err != nil {
		return
	}
	defer localConn.Close()
	go func() {
		if _, err := io.Copy(localConn, conn); err != nil {
			log.Println(err)
		}
	}()
	if _, err := io.Copy(conn, localConn); err != nil {
		log.Println(err)
	}
}
