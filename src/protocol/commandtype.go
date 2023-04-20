package protocol

import "errors"

//go:generate stringer -type CommandType
type CommandType int32

const (
	Nop CommandType = iota
	IdentifyCommand
	SSHProxyCommand
	HttpProxyCommand
	ReadyCommand
)

func CommandToBytes(cmd CommandType) []byte {
	return []byte{byte(cmd)}
}

func BytesToCommand(b []byte) (CommandType, error) {
	if len(b) != 1 {
		return Nop, errors.New("invalid command")
	}
	return CommandType(b[0]), nil
}
