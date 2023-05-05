package proxy

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// orwards datas from localhost:8881 to localhost:80
func TestForward(t *testing.T) {
	ln, err := net.Listen("tcp", ":8881")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	fmt.Println("new client")

	proxy, err := net.Dial("tcp", "127.0.0.1:80")
	if err != nil {
		panic(err)
	}

	fmt.Println("proxy connected")
	go ioCopy(conn, proxy)
	go ioCopy(proxy, conn)
}

func ioCopy(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}
