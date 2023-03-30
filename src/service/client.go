package service

import "net"

type client struct {
	net.Conn
	Hostname string

	ExitChan chan bool
}

func NewClient(conn net.Conn) *client {
	return &client{
		Conn:     conn,
		ExitChan: make(chan bool),
	}
}
