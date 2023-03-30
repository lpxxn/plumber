package service

import (
	"io"
	"net"
	"sync"

	"github.com/lpxxn/plumber/src/common"
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
	if magicStr != common.MagicString {
		log.Errorf("magic string not match: %s", magicStr)
		conn.Close()
		return
	}
	client := NewClient(conn)
	s.subCons.Store(conn.RemoteAddr(), client)
	// remove conn from subCons if conn is closed

	s.subCons.Delete(conn.RemoteAddr())
	client.Close()
}

// remove conn from subCons if conn is closed
func (s *Service) removeConn(conn net.Conn) {
	s.subCons.Delete(conn.RemoteAddr().String())
}

func (s *Service) Close() {
	s.subCons.Range(func(key, value interface{}) bool {
		value.(*client).Close()
		return true
	})
}

func IOLop(client *client) error {

	// read data from client
	for {

	}
	log.Infof("client(%s) host %s exit", client.Conn.RemoteAddr(), client.Hostname)
	close(client.ExitChan)
	return nil
}
