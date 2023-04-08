package service

import (
	"fmt"
	"io"
	"time"

	"github.com/lpxxn/plumber/src/log"
	"github.com/lpxxn/plumber/src/protocol"
)

type ServProtocol struct {
	Plumber *PlumberSrv
}

func NewServProtocol(srv *PlumberSrv) *ServProtocol {
	return &ServProtocol{Plumber: srv}
}

func (s *ServProtocol) IOLoop(c protocol.Client) error {
	client := c.(*client)
	var err error
	var line []byte
	var zeroTime time.Time

	// read data from client
	for {
		client.Conn.SetReadDeadline(zeroTime)
		// ReadSlice does not allocate new space for the data each request
		// ie. the returned slice is only valid until the next call to it
		line, err = client.Reader.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = fmt.Errorf("failed to read command - %s", err)
			}
			break
		}
		// trim \n
		line = line[:len(line)-1]
		// optional: trim \r
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}

	}
	log.Infof("client(%s) host %s exit", client.Conn.RemoteAddr(), client.Hostname)
	close(client.ExitChan)
	return nil
}
