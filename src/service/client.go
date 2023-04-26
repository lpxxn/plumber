package service

import (
	"bufio"
	"fmt"
	"net"
	"sync"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/log"
	"github.com/lpxxn/plumber/src/protocol"
)

const defaultBufferSize = 16 * 1024

type client struct {
	net.Conn
	Hostname string

	ExitChan chan bool
	// reading/writing interfaces
	Reader    *bufio.Reader
	Writer    *bufio.Writer
	Identity  *protocol.Identify
	writeLock sync.RWMutex
	sshProxy  *sshProxy
}

type sshProxy struct {
	// if the client is ssh proxy client, this field will be set
	LocalListener net.Listener
	SSHConfig     *config.SSHConf
}

func (s *sshProxy) Close() error {
	if s.LocalListener != nil {
		s.LocalListener.Close()
		s.LocalListener = nil
	}
	return nil
}

func (s *sshProxy) NewTCPServer() error {
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", s.SSHConfig.SrvPort))
	if err != nil {
		log.Errorf("sshProxy listen on %d failed: %v", s.SSHConfig.SrvPort, err)
		return err
	}
	return TCPServer(listner, s)
}

func (s *sshProxy) Handle(con net.Conn) {
	log.Infof("sshProxy: new connection from %s", con.RemoteAddr())
}

func NewClient(conn net.Conn) *client {
	return &client{
		Conn:     conn,
		Reader:   bufio.NewReaderSize(conn, defaultBufferSize),
		Writer:   bufio.NewWriterSize(conn, defaultBufferSize),
		ExitChan: make(chan bool),
		sshProxy: &sshProxy{},
	}
}

func (c *client) Close() error {
	if c.sshProxy != nil {
		c.sshProxy.Close()
	}
	return c.Conn.Close()
}

func (c *client) SendCommand(b []byte) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	_, err := c.Writer.Write(b)
	if err != nil {
		return err
	}
	return c.Writer.Flush()
}
