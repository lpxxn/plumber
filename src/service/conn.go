package service

import (
	"bufio"
	"net"
)

type Conn struct {
	net.Conn
	// reading/writing interfaces
	Reader *bufio.Reader
	Writer *bufio.Writer
}
