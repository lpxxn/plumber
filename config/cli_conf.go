package config

type CliConf struct {
	SrvTCPAddr string   `yaml:"srvTcpAddr"`
	SSH        *SSHConf `yaml:"ssh"`
}

type SSHConf struct {
	SrvPort      int    `yaml:"srvPort"`
	LocalSSHAddr string `yaml:"localSSHAddr"`
}
