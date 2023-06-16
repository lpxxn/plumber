package client

import (
	"fmt"
	"net"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/common"
	"github.com/lpxxn/plumber/src/log"
	"github.com/lpxxn/plumber/src/protocol"
)

type HttpProxy struct {
	Conf *config.ClientHttpProxyConf

	HttpProxyConn   net.Conn
	ForwardHttpConn net.Conn
	ServIP          string
}

func NewHttpProxy(servIP string, conf *config.ClientHttpProxyConf) *HttpProxy {
	return &HttpProxy{
		ServIP: servIP,
		Conf:   conf,
	}
}

func (h *HttpProxy) InitConn() error {
	localHttpProxyConn, err := h.ConnForwardHttpSrv()
	if err != nil {
		return err
	}
	if err := h.ConnRemoteSrv(); err != nil {
		return err
	}
	h.ForwardHttpConn = localHttpProxyConn
	return nil
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

func (h *HttpProxy) ConnRemoteSrv() error {
	httpConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", h.ServIP, h.Conf.RemotePort))
	if err != nil {
		log.Errorf("connect to http proxy failed: %s", err.Error())
		return err
	}
	_, err = httpConn.Write(common.HttpMagicBytes)
	if err != nil {
		log.Errorf("send http magic bytes failed: %s", err.Error())
		return err
	}
	httpCmd, err := protocol.HttpProxyCmd(h.Conf)
	if err != nil {
		log.Errorf("create http proxy cmd failed: %s", err.Error())
		return err
	}
	_, err = httpCmd.Write(httpConn)
	if err != nil {
		log.Errorf("send http proxy cmd failed: %s", err.Error())
		return err
	}

	h.HttpProxyConn = httpConn
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
