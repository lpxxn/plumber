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
	return nil
}

type SrvHttpProxyConf struct {
	Domain           string            `yaml:"domain"`
	Port             int               `yaml:"port"`
	DefaultForwardTo string            `yaml:"defaultForwardTo"`
	Forwards         []*SrvForwardConf `yaml:"forwards"`
}

type SrvForwardConf struct {
	Path string `yaml:"path"`
	// Forward to client or other server
	ForwardTo string `yaml:"forwardTo"`
}
