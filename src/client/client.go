package client

import (
	"github.com/lpxxn/plumber/config"
)

type Client struct {
	Hostname string

	ExitChan chan bool
	Conn     *Conn

	Conf *config.CliConf
}

func (c *Client) Close() error {
	return nil
}
