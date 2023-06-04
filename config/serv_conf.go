package config

import "errors"

type SrvConf struct {
	// Server TCP Port
	TCPAddr string `yaml:"tcpAddr"`
	// white list
	WhiteList []string          `yaml:"whiteList"`
	HttpProxy *SrvHttpProxyConf `yaml:"httpProxy"`
}

func NewSrvConf() *SrvConf {
	return &SrvConf{}
}

func (s *SrvConf) Validate() error {
	if s.TCPAddr == "" {
		return errors.New("tcpAddr is empty")
	}
	if s.HttpProxy != nil {
		if err := s.HttpProxy.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type SrvHttpProxyConf struct {
	Domain           string            `yaml:"domain"`
	Port             int               `yaml:"port"`
	DefaultForwardTo string            `yaml:"defaultForwardTo"`
	Forwards         []*SrvForwardConf `yaml:"forwards"`
}

func (c SrvHttpProxyConf) Validate() error {
	if c.Domain == "" {
		return errors.New("domain is empty")
	}
	if c.Port == 0 {
		return errors.New("port is empty")
	}
	if c.DefaultForwardTo == "" {
		return errors.New("defaultForwardTo is empty")
	}
	for _, forward := range c.Forwards {
		if err := forward.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type SrvForwardConf struct {
	Path string `yaml:"path"`
	// Forward to client or other server
	ForwardTo string `yaml:"forwardTo"`
}

func (c SrvForwardConf) Validate() error {
	return nil
}

type SrvForwardConfList []*SrvForwardConf

func (s SrvForwardConfList) Validate() error {
	for _, forward := range s {
		if err := forward.Validate(); err != nil {
			return err
		}
	}
	return nil
}
