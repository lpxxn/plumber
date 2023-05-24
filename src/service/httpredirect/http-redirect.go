package httpredirect

import (
	"bufio"
	"net"
	"sync"
)

type httpRedirectListener struct {
	net.Listener
}

func (l *httpRedirectListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	return &httpRedirectConn{
		Conn: c,
		r:    bufio.NewReader(c),
	}, nil
}

type httpRedirectConn struct {
	net.Conn
	once sync.Once
	r    *bufio.Reader
}

func (c *httpRedirectConn) Read(p []byte) (int, error) {
	return c.r.Read(p)
}

func (c *httpRedirectConn) CheckIsHttp() bool {
	firstBytes, err := c.r.Peek(5)
	if err != nil {
		return false
	}

	// If the request doesn't look like HTTP, then it's probably
	// TLS bytes and we don't need to do anything.
	if !firstBytesLookLikeHTTP(firstBytes) {
		return false
	}
	// Parse the HTTP request, so we can get the Host and URL to redirect to.
	//req, err := http.ReadRequest(c.r)
	//if err != nil {
	//	return 0, err
	//}

	return true
}

func firstBytesLookLikeHTTP(hdr []byte) bool {
	switch string(hdr[:5]) {
	case "GET /", "HEAD ", "POST ", "PUT /", "OPTIO":
		return true
	}
	return false
}
