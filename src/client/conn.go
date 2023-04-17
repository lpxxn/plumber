package client

import (
	"io"
	"net"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/common"
)

type Conn struct {
	*net.TCPConn

	// reading/writing interfaces
	r io.Reader
	w io.Writer

	Conf *config.CliConf
}

func NewConnect(conn *net.TCPConn) (*Conn, error) {
	c := &Conn{}
	c.TCPConn = conn
	c.r = conn
	c.w = conn
	if _, err := c.Write([]byte(common.MagicString)); err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}

func (c *Conn) Close() error {
	c.TCPConn.Close()
	return nil
}

func (c *Conn) SendCommand() {

}

type flusher interface {
	Flush() error
}

func (c *Conn) Flush() error {
	if w, ok := c.w.(flusher); ok {
		return w.Flush()
	}
	return nil
}
