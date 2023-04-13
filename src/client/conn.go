package client

import (
	"io"
	"net"
	"time"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/common"
)

type Conn struct {
	net.Conn

	// reading/writing interfaces
	r io.Reader
	w io.Writer

	Conf *config.CliConf
}

func (c *Conn) Connect() error {
	conn, err := net.Dial("tcp", c.Conf.SrvTCPAddr)
	if err != nil {
		return err
	}
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	}
	c.Conn = conn
	c.r = conn
	c.w = conn
	if _, err := c.Write([]byte(common.MagicString)); err != nil {
		c.Close()
		return err
	}
	return nil
}

func (c *Conn) Close() error {
	c.Conn.Close()
	return nil
}
