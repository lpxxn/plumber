package protocol

type CommonResultType int32

const (
	Success CommonResultType = iota
	Failed
)

type Identify struct {
	Hostname string `json:"hostname"`
	LocalIP  string `json:"localIP"`
}

type CommonResult struct {
	Code CommonResultType `json:"code"`
	Msg  string           `json:"msg"`
}
