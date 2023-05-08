// Code generated by "stringer -type CommandType"; DO NOT EDIT.

package protocol

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Nop-0]
	_ = x[IdentifyCommand-1]
	_ = x[SSHProxyCommand-2]
	_ = x[HttpProxyCommand-3]
	_ = x[PingCommand-4]
}

const _CommandType_name = "NopIdentifyCommandSSHProxyCommandHttpProxyCommandPingCommand"

var _CommandType_index = [...]uint8{0, 3, 18, 33, 49, 60}

func (i CommandType) String() string {
	if i < 0 || i >= CommandType(len(_CommandType_index)-1) {
		return "CommandType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CommandType_name[_CommandType_index[i]:_CommandType_index[i+1]]
}
