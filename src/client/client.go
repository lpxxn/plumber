package client

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/common"
	"github.com/lpxxn/plumber/src/log"
	"github.com/lpxxn/plumber/src/protocol"
)

type Client struct {
	Hostname string

	exitChan chan bool
	Conn     *Conn

	sshProxy  *SSHProxy
	httpProxy *HttpProxy
	Conf      *config.CliConf
	exist     uint32
	close     uint32
}

func NewClient(conf *config.CliConf) *Client {
	return &Client{
		Conf:     conf,
		exitChan: make(chan bool),
	}
}

func (c *Client) GetExitChan() <-chan bool {
	return c.exitChan
}

func (c *Client) Exit() {
	if !atomic.CompareAndSwapUint32(&c.exist, 0, 1) {
		return
	}
	close(c.exitChan)
}

func (c *Client) Close() error {
	if !atomic.CompareAndSwapUint32(&c.close, 0, 1) {
		return nil
	}
	c.Conn.Close()
	return nil
}

func (c *Client) DryRun() error {
	conn, err := c.connectToSrv()
	if err != nil {
		return err
	}

	defer conn.Close()
	if c.sshProxy != nil {
		if err := c.sshProxy.DryRun(); err != nil {
			return err
		}
	}
	if c.httpProxy != nil {
		if err := c.httpProxy.DryRun(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Run() error {
	return c.run()
}

func (c *Client) run() error {
	reconnectTimes := c.Conf.ReConnectionTimes
	reConnFun := func() {
		if reconnectTimes > 0 {
			reconnectTimes--
		}
		time.Sleep(time.Second * 3)
		log.Info("reconnect to server...")
	}
	for reconnectTimes != 0 {
		select {
		case <-c.exitChan:
			goto exit
		default:
		}
		if err := c.ConnectToSrv(); err != nil {
			// if err is timeout
			//if ne, ok := err.(net.Error); ok && ne.Timeout() {
			//	continue
			//}
			//goto exit
			reConnFun()
			continue
		}
		defer c.Close()

		if err := c.HandleSSHProxy(); err != nil {
			log.Errorf("handle ssh proxy failed: %s", err.Error())
			goto exit
		}
		if err := c.HandleHttpProxy(); err != nil {
			log.Errorf("handle http proxy failed: %s", err.Error())
			goto exit
		}
		// session
		NewCliProtocol(c).IOLoop()
		reConnFun()
	}
exit:
	c.Exit()
	log.Info("cli exit")
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
	log.Infof("connect to server(%s)", c.Conf.SrvTCPAddr)
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
	return c.Conn.Flush()
}

func (c *Client) ping() error {
	pingCmd := protocol.PingCmd()
	if _, err := pingCmd.Write(c.Conn.w); err != nil {
		return err
	}
	return c.Conn.Flush()
}

func (c *Client) HandleHttpProxy() error {
	if c.Conf.HttpProxy == nil {
		return nil
	}
	httpProxy := NewHttpProxy(c.Conf.SrvIP, c.Conf.HttpProxy)

	if err := httpProxy.InitConn(); err != nil {
		log.Errorf("init http proxy failed: %s", err.Error())
		return err
	}
	log.Infof("start http proxy")
	c.httpProxy = httpProxy
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
	result, err := protocol.ReadFrameData(c.Conn.r)
	if err != nil {
		log.Errorf("read ssh proxy result failed: %s", err.Error())
		return err
	}
	log.Debugf("ssh proxy result: %s code: %d", result.Msg, result.Code)
	if result.Code != protocol.Success {
		log.Errorf("ssh proxy failed: %s", result.Msg)
		return err
	}
	log.Infof("ssh proxy success")
	c.sshProxy = NewSSHProxy(c.Conf.SrvIP, c.Conf.SSH)

	if err := c.sshProxy.Handle(); err != nil {
		c.sshProxy.Close()
		return err
	}
	go func() {
		<-c.sshProxy.Exit
		c.sshProxy.Close()
		time.Sleep(time.Second * 1)
		c.HandleSSHProxy()
	}()
	return nil
}

func (c *Client) IsConnValid() bool {
	if c.Conn == nil {
		return false
	}

	// check c.Conn is disconnect
	if err := c.ping(); err != nil {
		return false
	}
	return true
}
