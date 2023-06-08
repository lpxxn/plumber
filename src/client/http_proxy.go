package client

import (
	"net"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/log"
)

type HttpProxy struct {
	Conf *config.ClientHttpProxyConf

	HttpProxyConn net.Conn
	ServIP        string
}

func NewHttpProxy(servIP string, conf *config.ClientHttpProxyConf) *HttpProxy {
	return &HttpProxy{
		ServIP: servIP,
		Conf:   conf,
	}
}

func (h *HttpProxy) Close() error {
	if h.HttpProxyConn != nil {
		h.HttpProxyConn.Close()
		h.HttpProxyConn = nil
	}
	return nil
}

func (h *HttpProxy) DryRun() error {
	return h.testConnection()
}

func (h *HttpProxy) testConnection() error {
	localHttpProxyConn, err := h.ConnForwardHttpSrv()
	if err != nil {
		return err
	}
	defer localHttpProxyConn.Close()
	return nil
}

func (h *HttpProxy) ConnForwardHttpSrv() (net.Conn, error) {
	proxyConn, err := net.Dial("tcp", h.Conf.LocalSrvAddr)
	if err != nil {
		log.Errorf("connect to local http server failed: %s", err.Error())
		return nil, err
	}
	log.Infof("connect to local http server [%s] success", proxyConn.LocalAddr())
	return proxyConn, nil
}
