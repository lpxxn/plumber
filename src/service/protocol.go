package service

import (
	"fmt"
	"io"
	"time"

	"github.com/lpxxn/plumber/src/common"
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
	var header []byte
	var zeroTime time.Time

	// read data from client
	for {
		client.Conn.SetReadDeadline(zeroTime)
		// ReadSlice does not allocate new space for the data each request
		// ie. the returned slice is only valid until the next call to it
		header, err = client.Reader.ReadSlice(common.NewLineByte)
		if err != nil {
			log.Errorf("failed to read command - %s", err)
			if err == io.EOF {
				err = nil
			} else {
				err = fmt.Errorf("failed to read command - %s", err)
			}
			break
		}
		log.Debugf("client(%s) host %s recv: %s", client.Conn.RemoteAddr(), client.Hostname, header)
		// trim \n
		header = header[:len(header)-1]
		cmdType := protocol.CommandType(header[0])
	}
	log.Infof("client(%s) host %s exit", client.Conn.RemoteAddr(), client.Hostname)
	close(client.ExitChan)
	return nil
}

func (s *ServProtocol) ExecCommand(c *client, cmd string) error {

	return nil
}
