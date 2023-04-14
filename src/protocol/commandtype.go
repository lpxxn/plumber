package protocol

//go:generate stringer -type CommandType
type CommandType int32

const (
	Nop CommandType = iota
	IdentifyCommand
	SSHProxyCommand
	HttpProxyCommand
)
