package client

import (
	"bufio"
	"net"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/common"
)

type Conn struct {
	*net.TCPConn

	// reading/writing interfaces
	r *bufio.Reader
	w *bufio.Writer

	Conf *config.CliConf
}

func NewConnect(conn *net.TCPConn) (*Conn, error) {
	c := &Conn{}
	c.TCPConn = conn
	c.r = bufio.NewReader(conn)
	c.w = bufio.NewWriter(conn)
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

func (c *Conn) Flush() error {
	return c.w.Flush()
}
