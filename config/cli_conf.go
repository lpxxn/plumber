package config

import (
	"errors"

	"github.com/lpxxn/plumber/src/common"
)

type CliConf struct {
	SrvTCPAddr        string               `yaml:"srvTcpAddr"`
	SrvIP             string               `yaml:"-"`
	SSH               *SSHConf             `yaml:"ssh"`
	HttpProxy         *ClientHttpProxyConf `yaml:"http"`
	ReConnectionTimes int32                `yaml:"reConnectionTimes"`
}

type SSHConf struct {
	SrvPort      int      `yaml:"srvPort"`
	LocalSSHAddr string   `yaml:"localSSHAddr"`
	WhiteList    []string `yaml:"whiteList"`
	ReConnTimes  int      `yaml:"reConnTimes"`
}

func (c *CliConf) Validate() error {
	if c.SrvTCPAddr == "" {
		return errors.New("srvTcpAddr is empty")
	}
	srvIP, err := common.TcpAddr(c.SrvTCPAddr)
	if err != nil {
		return err
	}
	c.SrvIP = srvIP.IP.String()
	if c.ReConnectionTimes == 0 {
		c.ReConnectionTimes = -1
	}
	if c.SSH != nil {
		if err := c.SSH.Validate(); err != nil {
			return err
		}
	}
	if c.HttpProxy != nil {
		if err := c.HttpProxy.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (s *SSHConf) Validate() error {
	if s.SrvPort == 0 {
		return errors.New("srvPort is empty")
	}
	if s.LocalSSHAddr == "" {
		return errors.New("localSSHAddr is empty")
	}
	if s.ReConnTimes <= 0 {
		s.ReConnTimes = -1
	}
	_, err := common.TcpAddr(s.LocalSSHAddr)
	return err
}

func NewCliConf() *CliConf {
	return &CliConf{}
}

type ClientHttpProxyConf struct {
	RemotePort   int32  `yaml:"remotePort"`
	UID          string `yaml:"uid"`
	LocalSrvAddr string `yaml:"localSrvAddr"`
}

func (c *ClientHttpProxyConf) Validate() error {
	if c.RemotePort <= 0 {
		return errors.New("remotePort is empty")
	}
	if c.UID == "" {
		return errors.New("uid is empty")
	}
	if c.LocalSrvAddr == "" {
		return errors.New("localSrvAddr is empty")
	}
	_, err := common.TcpAddr(c.LocalSrvAddr)
	return err
}
