package service

import (
	"net"
	"sync"
	"sync/atomic"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/common"
	"github.com/lpxxn/plumber/src/log"
)

type PlumberSrv struct {
	SrvAddr  string
	listener net.Listener
	subCons  sync.Map
	// tcp listener map[localPort]listener
	TcpListenerMap map[int]net.Listener

	isExiting int32
	Conf      *config.SrvConf
}

func NewService(conf *config.SrvConf) *PlumberSrv {
	return &PlumberSrv{
		SrvAddr: conf.TCPAddr,
		Conf:    conf,
	}
}

func (s *PlumberSrv) Run() {
	if err := s.HandleClientCommands(); err != nil {
		panic(err)
	}
	if err := s.HandleHttpForward(); err != nil {
		panic(err)
	}
}

func (s *PlumberSrv) HandleClientCommands() error {
	var err error
	s.listener, err = net.Listen("tcp", s.SrvAddr)
	if err != nil {
		return err
	}
	log.Infof("listen on %s", s.SrvAddr)
	log.Infof("start to accept connections")
	log.Debug(common.LocalPrivateIPV4())
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				panic(err)
			}
			go s.handleConnection(conn)
		}
	}()
	return nil
}

func (s *PlumberSrv) HandleHttpForward() error {
	if s.Conf.HttpProxy == nil {
		return nil
	}

	return nil
}

func (s *PlumberSrv) Exit() {
	if !atomic.CompareAndSwapInt32(&s.isExiting, 0, 1) {
		return
	}
	log.Info("stop to accept connections...")
	s.listener.Close()
	s.Close()
	log.Info("exit")
}

func (s *PlumberSrv) handleConnection(conn net.Conn) {
	if err := common.VerifyMagicStrConnection(conn); err != nil {
		conn.Close()
		return
	}

	protocol := NewServProtocol(s)

	client := NewClient(conn)
	s.subCons.Store(conn.RemoteAddr(), client)
	// remove conn from subCons if conn is closed
	go protocol.IOLoop(client)
	<-client.exitChan
	s.subCons.Delete(conn.RemoteAddr())
	client.Close()
}

// remove conn from subCons if conn is closed
func (s *PlumberSrv) removeConn(conn net.Conn) {
	s.subCons.Delete(conn.RemoteAddr().String())
}

func (s *PlumberSrv) Close() {
	s.subCons.Range(func(key, value interface{}) bool {
		value.(*client).Close()
		return true
	})
}

// Verify whether the IP is in the whitelist
func (s *PlumberSrv) VerifyIP(conn net.Conn) bool {
	return true
}
