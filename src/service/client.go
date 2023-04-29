package service

import (
	"bufio"
	"net"
	"sync"
	"sync/atomic"

	"github.com/lpxxn/plumber/src/log"
	"github.com/lpxxn/plumber/src/protocol"
	"github.com/lpxxn/plumber/src/proxy"
)

const defaultBufferSize = 16 * 1024

type client struct {
	net.Conn
	Hostname string

	// reading/writing interfaces
	Reader    *bufio.Reader
	Writer    *bufio.Writer
	Identity  *protocol.Identify
	writeLock sync.RWMutex
	sshProxy  *proxy.SSHProxy

	exitChan chan bool
	isClosed int32
}

func NewClient(conn net.Conn) *client {
	return &client{
		Conn:     conn,
		Reader:   bufio.NewReaderSize(conn, defaultBufferSize),
		Writer:   bufio.NewWriterSize(conn, defaultBufferSize),
		exitChan: make(chan bool),
		sshProxy: &proxy.SSHProxy{},
	}
}

func (c *client) Close() error {
	if !atomic.CompareAndSwapInt32(&c.isClosed, 0, 1) {
		log.Errorf("client %s is already closed", c.RemoteAddr())
		return nil
	}
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
