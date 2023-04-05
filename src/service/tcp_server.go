package service

import (
	"fmt"
	"net"
	"runtime"
	"strings"
	"sync"

	"github.com/lpxxn/plumber/src/log"
)

type TCPHandler interface {
	Handle(net.Conn)
}

func TCPServer(listener net.Listener, handler TCPHandler) error {
	log.Infof("TCPServer: listening on %s", listener.Addr())

	var wg sync.WaitGroup

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			// net.Error.Temporary() is deprecated, but is valid for accept
			// this is a hack to avoid a staticcheck error
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				log.Errorf("Temporary err: %s", err)
				runtime.Gosched()
				continue
			}
			// theres no direct way to detect this error because it is not exposed
			if !strings.Contains(err.Error(), "use of closed network connection") {
				return fmt.Errorf("listener.Accept() error - %s", err)
			}
			break
		}

		wg.Add(1)
		go func() {
			handler.Handle(clientConn)
			wg.Done()
		}()
	}
	wg.Wait()
	log.Infof("TCP: closing %s", listener.Addr())
	return nil
}
