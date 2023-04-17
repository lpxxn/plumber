package protocol

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/common"
)

type Command struct {
	Type   CommandType
	Params [][]byte
	Body   []byte
}

func (c *Command) Write(w io.Writer) (int64, error) {
	var total int64
	n, err := w.Write([]byte{byte(c.Type)})
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

var ErrInvalidCommand = errors.New("invalid command")

//func GetCmd(client *client) (*Command, error) {
//	header, err := r.ReadSlice(common.NewLineByte)
//	if err != nil {
//		log.Errorf("failed to read command - %s", err)
//		if err == io.EOF {
//			err = nil
//		} else {
//			err = fmt.Errorf("failed to read command - %s", err)
//		}
//		return nil, err
//	}
//	log.Debugf("client(%s) host %s recv: %s", , line)
//	// trim \n
//	header = header[:len(header)-1]
//	cmdType := CommandType(header[0])
//	switch cmdType {
//	case IdentifyCommand:
//		return getIdentifyCmd(r)
//	case SSHProxyCommand:
//		return getSSHProxyCmd(r)
//	default:
//		return nil, ErrInvalidCommand
//	}
//}

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
