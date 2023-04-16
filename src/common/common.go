package common

import (
	"errors"
	"net"
	"sync"

	"github.com/lpxxn/plumber/src/log"
)

const (
	MagicString = "GoV1"
)

var SeparatorBytes = []byte(" ")
var NewLineByte = byte('\n')
var NewLineBytes = []byte{NewLineByte}

type WaitGroup struct {
	sync.WaitGroup
}

func (w *WaitGroup) WaitFunc(f func()) {
	w.Add(1)
	go func() {
		f()
		w.Done()
	}()
}

func TcpAddr(addrStr string) (*net.TCPAddr, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		log.Errorf("Error resolving TCP address: %v\n", err)
		return tcpAddr, err
	}
	log.Infof("TCP address is valid: %v\n", tcpAddr)
	return tcpAddr, nil
}

type Validator interface {
	// Validate validates the given data.
	Validate() error
}

func LocalPrivateIPV4() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}
	return nil, errors.New("no private ipv4 address found")
}
