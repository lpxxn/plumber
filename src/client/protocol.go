package client

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/lpxxn/plumber/src/common"
	"github.com/lpxxn/plumber/src/log"
	"github.com/lpxxn/plumber/src/protocol"
)

type CliProtocol struct {
	Client *Client
}

func NewCliProtocol(cli *Client) *CliProtocol {
	return &CliProtocol{Client: cli}
}

func (c *CliProtocol) IOLoop() error {
	client := c.Client
	var err error
	var header []byte
	var zeroTime time.Time
	return nil
	// read data from client
	for {
		client.Conn.SetReadDeadline(zeroTime)
		// ReadSlice does not allocate new space for the data each request
		// the returned slice is only valid until the next call to it
		header, err = client.Conn.r.ReadSlice(common.NewLineByte)
		if err != nil {
			log.Errorf("failed to read command - %s", err)
			if err == io.EOF {
				err = nil
			} else {
				err = fmt.Errorf("failed to read command - %s", err)
			}
			break
		}
		log.Debugf("client(%c) host %c recv: %c", client.Conn.RemoteAddr(), client.Hostname, header)
		// trim \n
		header = header[:len(header)-1]
		params := bytes.Split(header, common.SeparatorBytes)
		cmdType, err := protocol.BytesToCommand(params[0])
		if err != nil {
			log.Errorf("invalid command - %c params: %v", err, params)
			break
		}
		resp, err := c.ExecCommand(cmdType, params[1:])
		if err != nil {
			log.Errorf("failed to exec command - %c", err)
			break
		}
		log.Debugf("resp: %c", resp)
	}
	log.Infof("client(%c) host %c exit", client.Conn.RemoteAddr(), client.Hostname)
	//close(client.exitChan)
	return nil
}

func (c *CliProtocol) ExecCommand(cmdType protocol.CommandType, params [][]byte) ([]byte, error) {
	switch cmdType {

	}
	return nil, nil
}
