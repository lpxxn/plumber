package protocol

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/lpxxn/plumber/src/common"
	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	iCmd, err := IdentifyCmd(&Identify{
		Hostname: "abcdefghijklmnopqrstuvwsyz",
		LocalIP:  "127.0.0.1",
	})
	assert.Nil(t, err)
	//w := os.Stdout
	w2 := bytes.NewBuffer(nil)
	buf := bufio.NewReader(w2)
	n, err := iCmd.Write(w2)
	assert.Nil(t, err)
	t.Log(n)
	commandStr, err := buf.ReadSlice(common.NewLineByte)
	assert.Nil(t, err)
	t.Log(commandStr)
	// trim the last byte \n
	commandStr = commandStr[:len(commandStr)-1]

	param := bytes.Split(commandStr, common.SeparatorBytes)
	t.Log(param)
	t.Logf("command: %s", CommandType(param[0][0]))

	//bodyLen := [4]byte{}
	//_, err = io.ReadFull(buf, bodyLen[:])
	//assert.Nil(t, err)
	//t.Log(bodyLen)
	//lenVal := binary.BigEndian.Uint32(bodyLen[:])
	//t.Log(lenVal)
	//body := make([]byte, lenVal)
	//_, err = io.ReadFull(buf, body)
	//assert.Nil(t, err)
	//t.Log(body)
	//identity := &Identify{}
	//err = json.Unmarshal(body, identity)
	//assert.Nil(t, err)
	identity, err := ReadIdentifyCommand(param[1:], buf)
	assert.Nil(t, err)
	t.Log(identity)
}

func TestIP(t *testing.T) {
	ip, err := common.LocalPrivateIPV4()
	assert.Nil(t, err)
	t.Log(ip, ip.String())
	a := []string{ip.String()}
	t.Log(a[0])
	t.Log(a[1:])
}
