package protocol

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/common"
	"github.com/lpxxn/plumber/src/log"
)

type Command struct {
	Type   CommandType
	Params [][]byte
	Body   []byte
}

func (c *Command) Write(w io.Writer) (int64, error) {
	var total int64
	n, err := w.Write(CommandToBytes(c.Type))
	total += int64(n)
	if err != nil {
		return total, err
	}
	for _, param := range c.Params {
		n, err := w.Write(common.SeparatorBytes)
		total += int64(n)
		if err != nil {
			return total, err
		}
		n, err = w.Write(param)
		total += int64(n)
		if err != nil {
			return total, err
		}
	}
	n, err = w.Write(common.NewLineBytes)
	total += int64(n)
	if err != nil {
		return total, err
	}

	if c.Body != nil {
		bodyLen := [4]byte{}
		binary.BigEndian.PutUint32(bodyLen[:], uint32(len(c.Body)))
		n, err := w.Write(bodyLen[:])
		total += int64(n)
		if err != nil {
			return total, err
		}
		n, err = w.Write(c.Body)
		total += int64(n)
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

func NewCommand(cmdType CommandType, params [][]byte, body []byte) *Command {
	return &Command{Type: cmdType, Params: params, Body: body}
}

var ErrInvalidCommand = errors.New("invalid command")

func ReadIdentifyCommand(r io.Reader) (*Identify, error) {
	body, err := ReadCommandData(r)
	if err != nil {
		return nil, err
	}
	identity := &Identify{}
	return identity, json.Unmarshal(body, identity)
}

func ReadSSHProxyCommand(params [][]byte, r io.Reader) (*config.SSHConf, error) {
	body, err := ReadCommandData(r)
	if err != nil {
		return nil, err
	}
	sshConf := &config.SSHConf{}
	return sshConf, json.Unmarshal(body, sshConf)
}

// ReadCommandData reads a single command from the provided io.Reader
// eg: | 4 byte length | body |
func ReadCommandData(r io.Reader) ([]byte, error) {
	bodyLen := [4]byte{}
	_, err := io.ReadFull(r, bodyLen[:])
	if err != nil {
		return nil, err
	}
	lenVal := binary.BigEndian.Uint32(bodyLen[:])
	body := make([]byte, lenVal)
	_, err = io.ReadFull(r, body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func IdentifyCmd(i *Identify) (*Command, error) {
	body, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return &Command{Type: IdentifyCommand, Body: body}, nil
}

func SSHProxyCmd(s *config.SSHConf) (*Command, error) {
	body, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return &Command{Type: SSHProxyCommand, Body: body}, nil
}

func HttpProxyCmd(s *config.ClientHttpProxyConf) (*Command, error) {
	body, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return &Command{Type: HttpProxyCommand, Body: body}, nil
}

func ReadClientHttpProxyCmd(r io.Reader) (*config.ClientHttpProxyConf, error) {
	body, err := ReadCommandData(r)
	if err != nil {
		return nil, err
	}
	httpConf := &config.ClientHttpProxyConf{}
	return httpConf, json.Unmarshal(body, httpConf)
}

func ReadCmdHeader(r *bufio.Reader) (CommandType, []byte, error) {
	header, err := r.ReadSlice(common.NewLineByte)
	if err != nil {
		log.Errorf("failed to read command - %s", err)
		if err == io.EOF {
			err = nil
		} else {
			err = fmt.Errorf("failed to read command - %s", err)
		}
		return 0, nil, err
	}
	// trim \n
	header = header[:len(header)-1]
	params := bytes.Split(header, common.SeparatorBytes)
	cmdType, err := BytesToCommand(params[0])
	if err != nil {
		log.Errorf("invalid command - %s params: %v", err, params)
		return 0, nil, err
	}
	return cmdType, header, nil
}

func PingCmd() *Command {
	return &Command{Type: PingCommand}
}
