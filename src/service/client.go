package service

import (
	"bufio"
	"net"

	"github.com/lpxxn/plumber/src/protocol"
)

const defaultBufferSize = 16 * 1024

type client struct {
	net.Conn
	Hostname string

	ExitChan chan bool
	// reading/writing interfaces
	Reader   *bufio.Reader
	Writer   *bufio.Writer
	Identity *protocol.Identify
}

func NewClient(conn net.Conn) *client {
	return &client{
		Conn:     conn,
		Reader:   bufio.NewReaderSize(conn, defaultBufferSize),
		Writer:   bufio.NewWriterSize(conn, defaultBufferSize),
		ExitChan: make(chan bool),
	}
}

func (c *client) Close() error {
	return c.Conn.Close()
}
