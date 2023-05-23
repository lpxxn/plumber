package httpredirect

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	go func() {
		go func() {
			http.ListenAndServe(":5678", &TestHandler{t: t})
		}()
		http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("my-header", "my-header-value")
			w.Write([]byte("hello2222"))
		})
		http.ListenAndServe(":5679", nil)
	}()
	time.Sleep(time.Second)
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:5678/test", nil)
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

}

type TestHandler struct {
	t *testing.T
}

func (h *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := h.t
	dialer := &net.Dialer{
		Timeout: time.Duration(time.Second * 30),
	}
	conn, err := dialer.Dial("tcp", ":5679")
	assert.Nil(t, err)
	assert.NotNil(t, conn)

	dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
		t.Logf("dialContext: %s, %s", network, address)
		return dialer.DialContext(ctx, network, ":5679")
	}
	rt := &http.Transport{
		Proxy:       http.ProxyFromEnvironment,
		DialContext: dialContext,
	}

	newUrl, err := url.Parse("http://127.0.0.1:5679")
	assert.Nil(t, err)
	r.URL.Scheme = newUrl.Scheme
	r.URL.Host = newUrl.Host
	resp, err := rt.RoundTrip(r)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	defer resp.Body.Close()

	// Copy response headers
	for k, v := range resp.Header {
		//w.Header()[k] = v
		w.Header().Set(k, v[0])
	}

	// Copy response status code
	//err = resp.Write(w)
	_, err = io.Copy(w, resp.Body)
	assert.Nil(t, err)
}
