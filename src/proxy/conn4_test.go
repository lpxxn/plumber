package proxy

import (
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/hashicorp/yamux"
)

// forwards datas from localhost:8881 to localhost:80
func TestForward(t *testing.T) {
	listener, err := net.Listen("tcp", ":8881")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
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

func TestYamux1(t *testing.T) {
	listener, err := net.Listen("tcp", ":8881")
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}
			session, err := yamux.Server(conn, nil)
			if err != nil {
				panic(err)
			}
			// accept a stream
			stream, err := session.Accept()
			if err != nil {
				panic(err)
			}
			buf := make([]byte, 4)
			_, err = stream.Read(buf)
			if err != nil {
				panic(err)
			}
			t.Logf("serv read: %s", buf)
			if _, err = stream.Write([]byte("pong")); err != nil {
				panic(err)
			}
		}
	}()

	clientFunc := func() {
		// client
		conn, err := net.Dial("tcp", ":8881")
		if err != nil {
			panic(err)
		}
		// setup client side of yamux
		session, err := yamux.Client(conn, nil)
		if err != nil {
			panic(err)
		}

		// open a new stream
		stream, err := session.Open()
		if err != nil {
			panic(err)
		}
		if _, err = stream.Write([]byte("ping")); err != nil {
			panic(err)
		}
		buf := make([]byte, 4)
		_, err = stream.Read(buf)
		if err != nil {
			panic(err)
		}
		t.Logf("client read: %s", buf)
	}
	clientFunc()
	clientFunc()

}

func TestYamux2(t *testing.T) {
	listener, err := net.Listen("tcp", ":8881")
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}
			session, err := yamux.Client(conn, nil)
			if err != nil {
				panic(err)
			}
			// accept a stream
			stream, err := session.Open()
			if err != nil {
				panic(err)
			}

			if _, err = stream.Write([]byte("pong")); err != nil {
				panic(err)
			}
			buf := make([]byte, 4)
			_, err = stream.Read(buf)
			if err != nil {
				panic(err)
			}
			t.Logf("serv read: %s", buf)
		}
	}()

	clientFunc := func() {
		// client
		conn, err := net.Dial("tcp", ":8881")
		if err != nil {
			panic(err)
		}
		// setup client side of yamux
		session, err := yamux.Server(conn, nil)
		if err != nil {
			panic(err)
		}

		// open a new stream
		stream, err := session.AcceptStream()
		if err != nil {
			panic(err)
		}
		t.Logf("StreamID: %d ", stream.StreamID())
		buf := make([]byte, 4)
		_, err = stream.Read(buf)
		if err != nil {
			panic(err)
		}
		if _, err = stream.Write([]byte("ping")); err != nil {
			panic(err)
		}
		t.Logf("client read: %s", buf)
	}
	clientFunc()
	clientFunc()

}
