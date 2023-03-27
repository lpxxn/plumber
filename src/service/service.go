package service

import "net"

type Service struct {
	SrvAddr  string
	listener net.Listener
	subCons  map[string]net.Conn
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
	defer conn.Close()
	s.subCons[conn.RemoteAddr().String()] = conn
	// remove conn from subCons if conn is closed
	
}
