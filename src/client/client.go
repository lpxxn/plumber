package client

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/common"
	"github.com/lpxxn/plumber/src/log"
	"github.com/lpxxn/plumber/src/protocol"
	"golang.org/x/sync/errgroup"
)

type Client struct {
	Hostname string

	ExitChan chan bool
	Conn     *Conn

	Conf *config.CliConf
}

func NewClient(conf *config.CliConf) *Client {
	return &Client{
		Conf:     conf,
		ExitChan: make(chan bool),
	}
}

func (c *Client) Close() error {
	return nil
}

func (c *Client) DryRun() error {
	conn, err := c.connectToSrv()
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)
	defer conn.Close()

	return nil
}

func (c *Client) ConnectToSrv() error {
	conn, err := c.connectToSrv()
	if err != nil {
		return err
	}
	if err := c.initTcpConn(conn.(*net.TCPConn)); err != nil {
		return err
	}

	if err := c.sendIdentify(); err != nil {
		log.Errorf("send identify to server failed: %s", err.Error())
		return err
	}
	c.Conn.r = bufio.NewReader(c.Conn.r)
	c.Conn.w = bufio.NewWriter(c.Conn.w)
	return nil
}

func (c *Client) initTcpConn(netCon *net.TCPConn) error {
	con, err := NewConnect(netCon)
	if err != nil {
		return err
	}
	c.Conn = con
	return nil
}

func (c *Client) connectToSrv() (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 15 * time.Second,
	}
	conn, err := dialer.Dial("tcp", c.Conf.SrvTCPAddr)
	if err != nil {
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			log.Errorf("connect to server(%s) timeout: %s", c.Conf.SrvTCPAddr, err.Error())
		}
		return nil, err
	}
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	}
	log.Infof("connect to server(%s) success", c.Conf.SrvTCPAddr)
	return conn, nil
}

func (c *Client) sendIdentify() error {
	// host name
	localIP, _ := common.LocalPrivateIPV4()
	identity := &protocol.Identify{
		Hostname: common.GetHostname(),
		LocalIP:  localIP.String(),
	}
	cmd, err := protocol.IdentifyCmd(identity)
	if err != nil {
		return err
	}
	if _, err := cmd.Write(c.Conn.w); err != nil {
		return err
	}
	return nil
}

// begin to handle ssh proxy
func (c *Client) HandleSSHProxy() error {
	if c.Conf.SSH == nil {
		return nil
	}
	log.Infof("start ssh proxy")
	cmd, err := protocol.SSHProxyCmd(c.Conf.SSH)
	if err != nil {
		log.Errorf("create ssh proxy cmd failed: %s", err.Error())
		return err
	}
	if _, err := cmd.Write(c.Conn.w); err != nil {
		log.Errorf("send ssh proxy cmd failed: %s", err.Error())
		return err
	}
	if err := c.Conn.Flush(); err != nil {
		return err
	}
	result, err := protocol.ReadCommonResultData(c.Conn.r)
	if err != nil {
		log.Errorf("read ssh proxy result failed: %s", err.Error())
		return err
	}
	if result.Code != protocol.Success {
		log.Errorf("ssh proxy failed: %s", result.Msg)
		return err
	}
	log.Infof("ssh proxy success")

	sshProxyConn, err := c.ConnSSHProxy()
	if err != nil {
		return err
	}

	localSSHFwdConn, err := c.ConnForwardSSHSrv()
	if err != nil {
		return err
	}
	eg := errgroup.Group{}
	eg.Go(func() error {
		return common.CopyDate(sshProxyConn, localSSHFwdConn)
	})
	eg.Go(func() error {
		return common.CopyDate(localSSHFwdConn, sshProxyConn)
	})
	err = eg.Wait()
	return err
}

func (c *Client) ConnSSHProxy() (net.Conn, error) {
	sshProxyConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.Conf.SrvIP, c.Conf.SSH.SrvPort))
	if err != nil {
		log.Errorf("connect to ssh proxy server failed: %s", err.Error())
		return nil, err
	}
	return sshProxyConn, nil
}

func (c *Client) ConnForwardSSHSrv() (net.Conn, error) {
	sshProxyConn, err := net.Dial("tcp", c.Conf.SSH.LocalSSHAddr)
	if err != nil {
		log.Errorf("connect to local ssh [%s] failed: %s", c.Conf.SSH.LocalSSHAddr, err.Error())
		return nil, err
	}
	return sshProxyConn, nil
}
