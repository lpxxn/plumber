package internal

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestTCPServer(t *testing.T) {
	// Listen for incoming connections
	ln, err := net.Listen("tcp", "localhost:8080")
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

	fmt.Println("New client connected:", conn.RemoteAddr())

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
		conn.Write([]byte("Echo: " + msg))
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
