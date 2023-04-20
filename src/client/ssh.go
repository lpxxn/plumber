package client

import (
	"github.com/lpxxn/plumber/src/protocol"
)

type SSHHandler struct {
	client *Client
}

func NewHandleSSH(client *Client) *SSHHandler {
	return &SSHHandler{
		client: client,
	}
}

func (s *SSHHandler) Close() error {
	return nil
}

func (s *SSHHandler) SendSSHCommand() error {
	if s.client.Conf.SSH == nil {
		return nil
	}
	cmd, err := protocol.SSHProxyCmd(s.client.Conf.SSH)
	if err != nil {
		return err
	}
	if _, err := cmd.Write(s.client.Conn.w); err != nil {
		return err
	}
	return s.client.Conn.Flush()
}
