package service

import (
	"io"
	"net"
	"sync"

	"github.com/lpxxn/plumber/src/log"
)

type Service struct {
	SrvAddr  string
	listener net.Listener
	subCons  sync.Map
}

func NewService(srvAddr string) *Service {
	return &Service{
		SrvAddr: srvAddr,
	}
}

func (s *Service) Start() {
	var err error
	s.listener, err = net.Listen("tcp", s.SrvAddr)
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

func (s *Service) handleConnection(conn net.Conn) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		log.Errorf("read magic error: %v", err)
		conn.Close()
		return
	}
	magicStr := string(buf)
	log.Infof("magicStr: %s", magicStr)
	s.subCons.Store(conn.RemoteAddr().String(), conn)
	// remove conn from subCons if conn is closed

}

// remove conn from subCons if conn is closed
func (s *Service) removeConn(conn net.Conn) {
	s.subCons.Delete(conn.RemoteAddr().String())
}

func (s *Service) Close() {
	s.subCons.Range(func(key, value interface{}) bool {
		value.(net.Conn).Close()
		return true
	})
}

// monitor conn
func (s *Service) monitorConn(conn net.Conn) {

}
