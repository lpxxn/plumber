package client

import (
	"bufio"
	"net"
	"time"

	"github.com/lpxxn/plumber/config"
)

type Conn struct {
	net.Conn

	// reading/writing interfaces
	Reader *bufio.Reader
	Writer *bufio.Writer

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
	return nil
}

func (c *Conn) Close() error {
	c.Conn.Close()
	return nil
}
