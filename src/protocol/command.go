package protocol

//go:generate stringer -type Command
type Command int32

const (
	Nop Command = iota
	SSHProxyCommand
	HttpProxyCommand
)
