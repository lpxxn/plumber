package proxy

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/lpxxn/plumber/src/common"
)

func TestTCPServer(t *testing.T) {
	hostname, err := os.Hostname()
	if err != nil {
		t.Fatalf("failed to get hostname: %v", err)
	}
	t.Logf("hostname: %s", hostname)
	// Listen for incoming connections
	//ln, err := net.Listen("tcp", "localhost:8080")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Listening on :8080")

	for {
		// Accept incoming connections
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Set the keepalive period
		tcpConn := conn.(*net.TCPConn)
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(60 * time.Second)

		// Handle incoming messages
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	defer conn.Close()
	clientIP := common.ClientIP(conn)
	fmt.Println("New client connected:", clientIP)
	fmt.Println("New client connected:", conn.RemoteAddr(), conn.LocalAddr())

	// Send a welcome message to the client
	conn.Write([]byte("Welcome to the server!"))

	// Read messages from the client
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			break
		}

		// Process the received message
		msg := string(buf[:n])
		fmt.Println("Received message from client:", msg)

		// Echo the message back to the client
		if _, err := conn.Write([]byte("Echo: " + msg)); err != nil {
			fmt.Println("Error writing to client:", err)
			break
		}
	}

	fmt.Println("Client disconnected:", conn.RemoteAddr())
}

func TestTCPClient(t *testing.T) {
	for {
		// Try to connect to the server
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			fmt.Println("Failed to connect:", err)
			time.Sleep(time.Second)
			continue
		}

		fmt.Println("Connected to server")

		// Keep the connection alive
		tcpConn := conn.(*net.TCPConn)
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(60 * time.Second)

		go func() {
			for {
				time.Sleep(2 * time.Second)
				if _, err := conn.Write([]byte("Hello")); err != nil {
					fmt.Println("Error writing to server:", err)
					break
				}
			}
		}()
		// Read messages from the server
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading from server:", err)
				break
			}

			// Process the received message
			msg := string(buf[:n])
			fmt.Println("Received message:", msg)
		}
		// Close the connection and retry after a delay
		fmt.Println("Disconnected from server")
		conn.Close()
		time.Sleep(time.Second)
	}
}

func TestAddr(t *testing.T) {
	h, p, err := net.SplitHostPort("127.0.0.1")
	t.Logf("h=%s, p=%s, err=%v", h, p, err)
	h, p, err = net.SplitHostPort("127.0.1:123")
	t.Logf("h=%s, p=%s, err=%v", h, p, err)
	h, p, err = net.SplitHostPort("127.0.0.1:80")
	t.Logf("h=%s, p=%s, err=%v", h, p, err)

	fc := func(addr string) {
		tcpAddr, err := common.TcpAddr(addr)
		t.Log(tcpAddr, err)
		h, p, err = net.SplitHostPort(tcpAddr.String())
		t.Logf("h=%s, p=%s, err=%v", h, p, err)
	}
	fc("")
	fc(":123")
	fc("127.0.1:123")
	fc("127.0.0.1")
	fc("127.0.0.1:22")
}
