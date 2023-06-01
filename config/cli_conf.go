package config

import (
	"errors"

	"github.com/lpxxn/plumber/src/common"
)

type CliConf struct {
	Name              string               `yaml:"name,required"`
	SrvTCPAddr        string               `yaml:"srvTcpAddr"`
	SrvIP             string               `yaml:"-"`
	SSH               *SSHConf             `yaml:"ssh"`
	HttpProxy         *ClientHttpProxyConf `yaml:"httpProxy"`
	ReConnectionTimes int32                `yaml:"reConnectionTimes"`
}

type SSHConf struct {
	SrvPort      int      `yaml:"srvPort"`
	LocalSSHAddr string   `yaml:"localSSHAddr"`
	WhiteList    []string `yaml:"whiteList"`
	ReConnTimes  int      `yaml:"reConnTimes"`
}

type ClientHttpProxyConf struct {
	LocalHttpSrvPort int `yaml:"LocalHttpSrvPort"`
}

func (c *CliConf) Validate() error {
	if c.Name == "" {
		return errors.New("name is empty")
	}
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
		return c.SSH.Validate()
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
