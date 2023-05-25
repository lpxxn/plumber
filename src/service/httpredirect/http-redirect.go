package httpredirect

import (
	"bufio"
	"net"
	"net/http"
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

func (c *httpRedirectConn) GetHttpRequest() (*http.Request, error) {
	if !c.CheckIsHttp() {
		return nil, nil
	}
	// Parse the HTTP request, so we can get the Host and URL to redirect to.
	return http.ReadRequest(c.r)
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

	return true
}

func firstBytesLookLikeHTTP(hdr []byte) bool {
	switch string(hdr[:5]) {
	case "GET /", "HEAD ", "POST ", "PUT /", "OPTIO":
		return true
	}
	return false
}
