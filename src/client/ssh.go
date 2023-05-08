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
	sshConf    *config.SSHConf `yaml:"ssh"`
	ServIP     string
	Exit       chan struct{}
	Preparer   Preparer
	AfterExist func()
}

func NewSSHProxy(servIP string, sshConf *config.SSHConf, prepare Preparer, exist func()) *SSHProxy {
	return &SSHProxy{
		sshConf:    sshConf,
		ServIP:     servIP,
		Exit:       make(chan struct{}),
		Preparer:   prepare,
		AfterExist: exist,
	}
}

func (s *SSHProxy) Close() error {
	return nil
}

func (s *SSHProxy) DryRun() error {
	return s.testConnection()
}

func (s *SSHProxy) testConnection() error {
	sshProxyConn, err := s.ConnSSHProxy(s.ServIP)
	if err != nil {
		log.Errorf("connect to ssh proxy server failed: %s", err.Error())
		return err
	}
	defer sshProxyConn.Close()
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
	for s.sshConf.ReConnTimes != 0 {
		select {
		case <-s.Exit:
			log.Infof("ssh handler exit")
			goto exit
		default:
		}
		if err := s.Preparer(); err != nil {
			log.Errorf("prepare failed: %s", err.Error())
			goto exit
		}
		sshProxyConn, err := s.ConnSSHProxy(s.ServIP)
		if err != nil {
			log.Errorf("connect to ssh proxy server failed: %s", err.Error())
			s.fallback()
			continue
		}
		defer sshProxyConn.Close()

		session, err := yamux.Server(sshProxyConn, common.NewYamuxConfig())
		if err != nil {
			log.Errorf("create yamux session failed: %s", err.Error())
			return err
		}
		for {
			stream, err := session.AcceptStream()
			if err != nil {
				log.Errorf("accept stream failed: %s", err.Error())
				break
			}
			log.Infof("stream %d accepted", stream.StreamID())
			go func() {
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
				err = eg.Wait()
			}()
		}
	}
exit:
	log.Infof("ssh handler exit")
	s.AfterExist()
	return nil
}

func (s *SSHProxy) ConnSSHProxy(srvIP string) (net.Conn, error) {
	sshProxyConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", srvIP, s.sshConf.SrvPort))
	if err != nil {
		log.Errorf("connect to ssh proxy server failed: %s", err.Error())
		return nil, err
	}
	if _, err := sshProxyConn.Write([]byte(common.SSHMagicString)); err != nil {
		log.Errorf("write magic string to ssh proxy server failed: %s", err.Error())
		s.Close()
		return nil, err
	}
	log.Infof("connect to ssh proxy server [%s] success", sshProxyConn.RemoteAddr())
	return sshProxyConn, nil
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
