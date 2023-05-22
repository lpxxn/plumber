package httpredirect

import (
	"bufio"
	"bytes"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"
)

type testConn struct {
	r bytes.Buffer
	w bytes.Buffer
}

func (c *testConn) Read(b []byte) (int, error)  { return c.r.Read(b) }
func (c *testConn) Write(b []byte) (int, error) { return c.w.Write(b) }
func (*testConn) Close() error                  { return nil }

func (*testConn) LocalAddr() net.Addr                { return &net.TCPAddr{Port: 0, Zone: "", IP: net.IPv4zero} }
func (*testConn) RemoteAddr() net.Addr               { return &net.TCPAddr{Port: 0, Zone: "", IP: net.IPv4zero} }
func (*testConn) SetDeadline(_ time.Time) error      { return nil }
func (*testConn) SetReadDeadline(_ time.Time) error  { return nil }
func (*testConn) SetWriteDeadline(_ time.Time) error { return nil }

func TestHttpRedirectPeekHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/test/op", nil)
	// Dump raw http request
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		t.Fatal(err)
	}
	conn := new(testConn)

	if _, err := conn.r.Write(dump); err != nil {
		t.Fatal(err)
	}

	httpRedirectCon := &httpRedirectConn{
		Conn: conn,
		r:    bufio.NewReader(conn),
	}
	idx, err := httpRedirectCon.Read([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(idx)
}

func TestReverse(t *testing.T) {
	// Target URL
	target, err := url.Parse("http://example.com")
	if err != nil {
		panic(err)
	}

	// Proxy handler
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Serve HTTP
	http.ListenAndServe(":8000", proxy)
}

func TestReverse2(t *testing.T) {
	dialer := &net.Dialer{
		Timeout: time.Duration(time.Second * 30),
	}
	dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, address)
	}
	rt := &http.Transport{
		Proxy:       http.ProxyFromEnvironment,
		DialContext: dialContext,
	}

}
