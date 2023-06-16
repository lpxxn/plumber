package httpredirect

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"sync"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/common"
	"github.com/lpxxn/plumber/src/log"
	"github.com/lpxxn/plumber/src/protocol"
)

func NewHttpProxy(conf *config.SrvHttpProxyConf) (*HttpProxySrv, error) {
	ln, err := net.Listen("tcp", conf.LocalSrvAddress())
	if err != nil {
		return nil, err
	}
	log.Infof("http proxy listen on %s", conf.LocalSrvAddress())
	var srv = &HttpProxySrv{
		Listener:               ln,
		Router:                 nil,
		Conf:                   conf,
		HttpProxyClientConnMap: make(map[string]net.Conn),
	}
	if conf.DefaultForwardTo != "" {
		srv.DefaultForwardConn = srv.ForwardConn(conf.DefaultForwardTo)
	}
	srv.Router, err = srv.ParseRouter(conf)
	if err != nil {
		return nil, err
	}
	go srv.Handle()
	return srv, nil
}

type HttpProxySrv struct {
	net.Listener
	Router                 *Router
	Conf                   *config.SrvHttpProxyConf
	DefaultForwardConn     func() (net.Conn, error)
	HttpProxyClientConnMap map[string]net.Conn
	// lock
	lock sync.Mutex
}

func (l *HttpProxySrv) AddClient(identity *config.ClientHttpProxyConf, conn net.Conn) error {
	// use lock to avoid concurrent map write
	l.lock.Lock()
	defer l.lock.Unlock()
	if _, ok := l.HttpProxyClientConnMap[identity.UID]; ok {
		log.Infof("client %s already exists", identity.UID)
		return errors.New("client already exists")
	}
	log.Infof("add client %s", identity.UID)
	l.HttpProxyClientConnMap[identity.UID] = conn
	return nil
}

func (l *HttpProxySrv) GetClient(uid string) (net.Conn, bool) {
	l.lock.Lock()
	defer l.lock.Unlock()
	conn, ok := l.HttpProxyClientConnMap[uid]
	return conn, ok
}

func (l *HttpProxySrv) ParseRouter(conf *config.SrvHttpProxyConf) (*Router, error) {
	router := NewRouter()
	for _, item := range conf.Forwards {
		router.Add(item.Path)
		if item.ForwardTo != "" {
			router.ForwardConn = l.ForwardConn(item.ForwardTo)
		}
	}
	return router, nil
}

func (l *HttpProxySrv) ForwardConn(addr string) func() (net.Conn, error) {
	return func() (net.Conn, error) {
		if conn, ok := l.GetClient(addr); ok {
			return conn, nil
		}
		return net.Dial("tcp", addr)
	}
}

func (l *HttpProxySrv) Handle() {
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
			if isClient, err := hc.CheckReadRemoteClient(); isClient {
				log.Infof("is client %s", hc.RemoteAddr())
				cmdType, header, err := protocol.ReadCmdHeader(hc.r)
				if err != nil {
					log.Errorf("read cmd header error: %s", err.Error())
					return
				}
				log.Infof("cmdType: %s, header: %s", cmdType, header)
				identity, err := protocol.ReadClientHttpProxyCmd(hc)
				if err != nil {
					log.Errorf("read httpProxy command error: %s", err.Error())
					return
				}

				if err := l.AddClient(identity, hc); err != nil {
					log.Errorf("add client error: %s", err.Error())
					return
				}
				return
			} else if err != nil {
				log.Errorf("is client: %t, err: %v", isClient, err)
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
			resp, err := http.ReadResponse(bufio.NewReader(forwardConn), req)
			if err != nil {
				log.Errorf("copy error: %s", err.Error())
				return
			}

			if err := resp.Write(hc); err != nil {
				log.Errorf("write http response error: %s", err.Error())
				return
			}
		}(conn)
	}
}

func (l *HttpProxySrv) ConnByRouter(r *Route) (net.Conn, error) {
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

func (l *HttpProxySrv) Accept() (net.Conn, error) {
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

func (c *httpRedirectConn) CheckReadRemoteClient() (bool, error) {
	firstBytes, err := c.r.Peek(len(common.HttpMagicBytes))
	if err != nil {
		return false, nil
	}

	if string(firstBytes) != common.HttpMagicString {
		return false, nil
	}
	r := make([]byte, len(common.HttpMagicBytes))
	_, err = c.Read(r)
	if err != nil {
		return false, err
	}
	log.Infof("httpproxy client conn magic: %s", string(r))
	return true, nil
}

func firstBytesLookLikeHTTP(hdr []byte) bool {
	switch string(hdr[:5]) {
	case "GET /", "HEAD ", "POST ", "PUT /", "OPTIO":
		return true
	}
	return false
}
