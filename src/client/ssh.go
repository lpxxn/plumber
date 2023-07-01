package client

import (
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/yamux"
	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/common"
	"github.com/lpxxn/plumber/src/log"
	"golang.org/x/sync/errgroup"
)

type Preparer func() error

type SSHProxy struct {
	sshConf      *config.SSHConf `yaml:"ssh"`
	SSHProxyConn net.Conn
	ServIP       string
	Exit         chan struct{}
}

func NewSSHProxy(servIP string, sshConf *config.SSHConf) *SSHProxy {
	return &SSHProxy{
		sshConf: sshConf,
		ServIP:  servIP,
		Exit:    make(chan struct{}),
	}
}

func (s *SSHProxy) Close() error {
	if s.SSHProxyConn != nil {
		s.SSHProxyConn.Close()
		s.SSHProxyConn = nil
	}
	return nil
}

func (s *SSHProxy) DryRun() error {
	return s.testConnection()
}

func (s *SSHProxy) testConnection() error {
	err := s.ConnSSHProxy(s.ServIP)
	if err != nil {
		log.Errorf("connect to ssh proxy server failed: %s", err.Error())
		return err
	}
	defer s.SSHProxyConn.Close()
	localSSHFwdConn, err := s.ConnForwardSSHSrv()
	if err != nil {
		log.Errorf("connect to local ssh server failed: %s", err.Error())
		return err
	}
	defer localSSHFwdConn.Close()
	return nil
}

func (s *SSHProxy) reduceReConnTimes() {
	if s.sshConf.ReConnTimes > 0 {
		s.sshConf.ReConnTimes--
	}
}

func (s *SSHProxy) fallback() {
	log.Infof("fallback to direct connection")
	s.reduceReConnTimes()
	time.Sleep(time.Second * 2)
}

func (s *SSHProxy) Handle() error {
	err := s.ConnSSHProxy(s.ServIP)
	if err != nil {
		log.Errorf("connect to ssh proxy server failed: %s", err.Error())
		return nil
	}
	if testLocalSSHFwdConn, err := s.ConnForwardSSHSrv(); err != nil {
		log.Errorf("connect to local ssh server [%s] failed: %s", s.sshConf.LocalSSHAddr, err.Error())
		return err
	} else {
		testLocalSSHFwdConn.Close()
	}
	session, err := yamux.Server(s.SSHProxyConn, common.NewYamuxConfig())
	if err != nil {
		log.Errorf("create yamux session failed: %s", err.Error())
		return err
	}
	go func() {
		defer session.Close()
		for {
			stream, err := session.AcceptStream()
			if err != nil {
				log.Errorf("accept stream failed: %s", err.Error())
				close(s.Exit)
				return
			}
			log.Infof("stream %d accepted", stream.StreamID())
			go func() {
				defer stream.Close()
				localSSHFwdConn, err := s.ConnForwardSSHSrv()
				if err != nil {
					log.Errorf("connect to local ssh server [%s] failed: %s", s.sshConf.LocalSSHAddr, err.Error())
					panic(err)
				}
				defer localSSHFwdConn.Close()
				eg := errgroup.Group{}
				eg.Go(func() error {
					return common.CopyDate(stream, localSSHFwdConn)
				})
				eg.Go(func() error {
					return common.CopyDate(localSSHFwdConn, stream)
				})
				log.Infof("wait for ssh proxy exit")
				if err = eg.Wait(); err != nil {
					log.Errorf("copy data failed: %s", err.Error())
				}
				log.Infof("stream %d closed", stream.StreamID())
			}()
		}
	}()

	return nil
}

func (s *SSHProxy) ConnSSHProxy(srvIP string) error {
	sshProxyConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", srvIP, s.sshConf.SrvPort))
	if err != nil {
		log.Errorf("connect to ssh proxy server failed: %s", err.Error())
		return err
	}
	if _, err := sshProxyConn.Write([]byte(common.SSHMagicString)); err != nil {
		log.Errorf("write magic string to ssh proxy server failed: %s", err.Error())
		s.Close()
		return err
	}
	s.SSHProxyConn = sshProxyConn
	log.Infof("connect to ssh proxy server [%s] success", sshProxyConn.RemoteAddr())
	return nil
}

func (s *SSHProxy) ConnForwardSSHSrv() (net.Conn, error) {
	sshProxyConn, err := net.Dial("tcp", s.sshConf.LocalSSHAddr)
	if err != nil {
		log.Errorf("connect to local ssh [%s] failed: %s", s.sshConf.LocalSSHAddr, err.Error())
		return nil, err
	}
	log.Infof("connect to local ssh [%s] success", sshProxyConn.LocalAddr())
	return sshProxyConn, nil
}
