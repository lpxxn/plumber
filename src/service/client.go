package service

import (
	"bufio"
	"net"
	"sync"
	"sync/atomic"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/common"
	"github.com/lpxxn/plumber/src/log"
	"github.com/lpxxn/plumber/src/protocol"
	"github.com/lpxxn/plumber/src/proxy"
	"golang.org/x/sync/errgroup"
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

func (c *client) StartSSHProxy(config *config.SSHConf) error {
	c.sshProxy = &proxy.SSHProxy{
		SSHConfig: config,
	}

	if err := c.sshProxy.NewTCPServer(); err != nil {
		return err
	}
	if _, err := protocol.SendFrameData(c.Conn, &protocol.FrameResp{
		Code: protocol.Success,
		Msg:  "success",
	}); err != nil {
		log.Errorf("write ready command failed: %v", err)
		return err
	}
	if err := c.sshProxy.WaitForTunnelConn(); err != nil {
		return err
	}
	log.Infof("client %s start ssh proxy success", c.RemoteAddr())
	return c.sshProxy.Start(func(conn net.Conn) {
		log.Infof("new tcp connection from ssh proxy: %s", conn.RemoteAddr())
		eg := errgroup.Group{}
		eg.Go(func() error {
			return common.CopyDate(conn, c.sshProxy.RemoteTunnelConn)
		})
		eg.Go(func() error {
			return common.CopyDate(c.sshProxy.RemoteTunnelConn, conn)
		})
		err := eg.Wait()
		if err != nil {
			log.Errorf("copy data failed: %v", err)
		}
	})
}
