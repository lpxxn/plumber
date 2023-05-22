package httpredirect

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"sync"
)

type httpRedirectConn struct {
	net.Conn
	once sync.Once
	r    *bufio.Reader
}

// Read tries to peek at the first few bytes of the request, and if we get
// an error reading the headers, and that error was due to the bytes looking
// like an HTTP request, then we perform a HTTP->HTTPS redirect on the same
// port as the original connection.
func (c *httpRedirectConn) Read(p []byte) (int, error) {
	var errReturn error
	c.once.Do(func() {
		firstBytes, err := c.r.Peek(5)
		if err != nil {
			return
		}

		// If the request doesn't look like HTTP, then it's probably
		// TLS bytes and we don't need to do anything.
		if !firstBytesLookLikeHTTP(firstBytes) {
			return
		}

		// Parse the HTTP request, so we can get the Host and URL to redirect to.
		req, err := http.ReadRequest(c.r)
		if err != nil {
			return
		}

		// Build the redirect response, using the same Host and URL,
		// but replacing the scheme with https.
		headers := make(http.Header)
		headers.Add("Location", "https://"+req.Host+req.URL.String())
		resp := &http.Response{
			Proto:      "HTTP/1.0",
			Status:     "308 Permanent Redirect",
			StatusCode: 308,
			ProtoMajor: 1,
			ProtoMinor: 0,
			Header:     headers,
		}

		err = resp.Write(c.Conn)
		if err != nil {
			errReturn = fmt.Errorf("couldn't write HTTP->HTTPS redirect")
			return
		}

		errReturn = fmt.Errorf("redirected HTTP request on HTTPS port")
		c.Conn.Close()
	})

	if errReturn != nil {
		return 0, errReturn
	}

	return c.r.Read(p)
}

func firstBytesLookLikeHTTP(hdr []byte) bool {
	switch string(hdr[:5]) {
	case "GET /", "HEAD ", "POST ", "PUT /", "OPTIO":
		return true
	}
	return false
}
