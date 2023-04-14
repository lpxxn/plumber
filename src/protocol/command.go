package protocol

import (
	"encoding/binary"
	"encoding/json"
	"io"

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

func IdentifyCmd(i *Identify) (*Command, error) {
	body, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return &Command{Type: IdentifyCommand, Body: body}, nil
}
