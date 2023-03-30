package service

import "net"

type client struct {
	net.Conn
	ID       string
	Hostname string

	ExitChan chan bool
}
