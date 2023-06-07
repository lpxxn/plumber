package httpredirect

import (
	"bufio"
	"errors"
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/log"
)

type ClientConnections interface {
	All() []net.Conn
	GetByName(name string) net.Conn
}

func NewHttpProxy(clientConnections ClientConnections, conf *config.SrvHttpProxyConf) (*Listener, error) {
	ln, err := net.Listen("tcp", conf.LocalSrvAddress())
	if err != nil {
		return nil, err
	}
	var listener = &Listener{
		Listener:          ln,
		Router:            nil,
		Conf:              conf,
		ClientConnections: clientConnections,
	}
	if conf.DefaultForwardTo != "" {
		listener.DefaultForwardConn = listener.ForwardConn(conf.DefaultForwardTo)
	}
	listener.Router, err = listener.ParseRouter(conf)
	if err != nil {
		return nil, err
	}
	go listener.Handle()
	return listener, nil
}

type Listener struct {
	net.Listener
	Router             *Router
	Conf               *config.SrvHttpProxyConf
	ClientConnections  ClientConnections
	DefaultForwardConn func() (net.Conn, error)
}

func (l *Listener) ParseRouter(conf *config.SrvHttpProxyConf) (*Router, error) {
	router := NewRouter()
	for _, item := range conf.Forwards {
		router.Add(item.Path)
		if item.ForwardTo != "" {
			router.ForwardConn = l.ForwardConn(item.ForwardTo)
		}
	}
	return router, nil
}

func (l *Listener) ForwardConn(addr string) func() (net.Conn, error) {
	return func() (net.Conn, error) {
		if conn := l.ClientConnections.GetByName(addr); conn != nil {
			return conn, nil
		}
		return net.Dial("tcp", addr)
	}
}

func (l *Listener) Handle() {
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Errorf("accept error: %s", err.Error())
			continue
		}
		go func(conn net.Conn) {
			hc, ok := conn.(*httpRedirectConn)
			if !ok {
				log.Errorf("conn is not httpRedirectConn")
				return
			}
			req, err := hc.GetHttpRequest()
			if err != nil {
				log.Errorf("get http request error: %s", err.Error())
				return
			}
			r := l.Router.MatchRoute(req.URL.Path)
			forwardConn, err := l.ConnByRouter(r)
			if err != nil {
				log.Errorf("get forward conn error: %s", err.Error())
				return
			}

			if err := req.Write(forwardConn); err != nil {
				log.Errorf("write http request error: %s", err.Error())
				return
			}
			if _, err = io.Copy(conn, forwardConn); err != nil {
				log.Errorf("copy error: %s", err.Error())
				return
			}
		}(conn)
	}
}

func (l *Listener) ConnByRouter(r *Route) (net.Conn, error) {
	if r == nil {
		if l.DefaultForwardConn == nil {
			log.Errorf("no route matched and no default forward conn")
			return nil, errors.New("no route matched and no default forward conn")
		}
		return l.DefaultForwardConn()

	}
	if r.ForwardConn == nil {
		log.Errorf("route: %s no forward conn", r.OriginPath)
		return nil, errors.New("route: " + r.OriginPath + " no forward conn")
	}
	return r.ForwardConn()
}

func (l *Listener) Accept() (net.Conn, error) {
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
