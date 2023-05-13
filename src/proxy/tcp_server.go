package proxy

import (
	"fmt"
	"net"
	"runtime"
	"strings"
	"sync"

	"github.com/lpxxn/plumber/src/log"
)

type TCPHandler func(net.Conn)

func TCPServer(listener net.Listener, handler TCPHandler) error {
	log.Infof("TCPServer: listening on %s", listener.Addr())

	var wg sync.WaitGroup

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			if err, ok := err.(interface{ Temporary() bool }); ok && err.Temporary() {
				log.Errorf("Temporary err: %s", err)
				runtime.Gosched()
				continue
			}
			if !strings.Contains(err.Error(), "use of closed network connection") {
				return fmt.Errorf("listener.Accept() error - %s", err)
			}
			break
		}

		wg.Add(1)
		go func() {
			defer clientConn.Close()
			handler(clientConn)
			wg.Done()
		}()
	}
	wg.Wait()
	log.Infof("TCP: closing %s", listener.Addr())
	return nil
}
