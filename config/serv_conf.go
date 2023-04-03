package config

type SrvConf struct {
	// Server TCP Port
	TCPAddr string `yaml:"tcpAddr"`
}

func NewSrvConf() *SrvConf {
	return &SrvConf{}
}
