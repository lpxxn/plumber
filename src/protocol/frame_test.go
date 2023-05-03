package protocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFrameData(t *testing.T) {
	f := FrameResp{
		Code: Success,
		Msg:  "success",
	}
	dummyWriter := bytes.NewBuffer(nil)
	_, err := SendFrameData(dummyWriter, &f)
	assert.Nil(t, err)

	rev, err := ReadFrameData(dummyWriter)
	assert.Nil(t, err)
	assert.Equal(t, f.Code, rev.Code)
	assert.Equal(t, f.Msg, rev.Msg)
}
