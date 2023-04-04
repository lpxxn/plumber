package config

import "errors"

type SrvConf struct {
	// Server TCP Port
	TCPAddr string `yaml:"tcpAddr"`
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
