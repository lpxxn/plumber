package config

import (
	"errors"

	"github.com/lpxxn/plumber/src/common"
)

type CliConf struct {
	SrvTCPAddr string         `yaml:"srvTcpAddr"`
	SSH        *SSHConf       `yaml:"ssh"`
	HttpProxy  *HttpProxyConf `yaml:"httpProxy"`
}

type SSHConf struct {
	SrvPort      int      `yaml:"srvPort"`
	LocalSSHAddr string   `yaml:"localSSHAddr"`
	WhiteList    []string `yaml:"whiteList"`
}

type HttpProxyConf struct {
	SrvPort int `yaml:"srvPort"`
}

func (c *CliConf) Validate() error {
	if c.SrvTCPAddr == "" {
		return errors.New("srvTcpAddr is empty")
	}
	if _, err := common.TcpAddr(c.SrvTCPAddr); err != nil {
		return err
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
	_, err := common.TcpAddr(s.LocalSSHAddr)
	return err
}

func NewCliConf() *CliConf {
	return &CliConf{}
}
