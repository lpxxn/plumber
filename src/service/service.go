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

// remove conn from subCons if conn is closed
func (s *Service) removeConn(conn net.Conn) {
	delete(s.subCons, conn.RemoteAddr().String())
}

// monitor conn
func (s *Service) monitorConn(conn net.Conn) {
	for {
		// read from conn
		// if error, remove conn from subCons
		data, err := conn.Read()
		if err != nil {
			s.removeConn(conn)
			break
		}
		// if read data, send data to all subCons
		for _, subConn := range s.subCons {
			_, err := subConn.Write(data)
			if err != nil {
				s.removeConn(subConn)
			}
		}
	}
}
