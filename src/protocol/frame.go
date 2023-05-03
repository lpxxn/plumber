package protocol

import (
	"encoding/binary"
	"encoding/json"
	"io"
)

type FrameDataType int32

const (
	Success FrameDataType = iota
	Failed
)

type FrameResp struct {
	Code FrameDataType `json:"code"`
	Msg  string        `json:"msg"`
}

func ReadFrameData(r io.Reader) (*FrameResp, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	dataSize := binary.BigEndian.Uint32(buf)
	data := make([]byte, dataSize)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}
	resp := &FrameResp{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func SendFrameData(w io.Writer, resp *FrameResp) (int, error) {
	buf := make([]byte, 4)
	data, err := json.Marshal(resp)
	if err != nil {
		return -1, err
	}
	dateSize := uint32(len(data))
	binary.BigEndian.PutUint32(buf, dateSize)
	n, err := w.Write(buf)
	if err != nil {
		return n, err
	}
	return w.Write(data)
}
