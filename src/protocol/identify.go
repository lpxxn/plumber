package protocol

type Identify struct {
	Hostname   string `json:"hostname"`
	LocalIP    string `json:"localIP"`
	ClientName string `json:"clientName"` // clientName is the name of the client, which is used to identify the client, unique
}
