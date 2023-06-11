package protocol

type Identify struct {
	Hostname string `json:"hostname"`
	LocalIP  string `json:"localIP"`
	UID      string `json:"uid"`
}
