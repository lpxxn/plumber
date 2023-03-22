package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

func main() {
	// Listen for incoming TCP connections on port 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer ln.Close()

	fmt.Println("TCP server listening on port 8080")

	// Loop forever to accept incoming client connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}

		// Handle incoming client connection in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("New client connected:", conn.RemoteAddr().String())

	// Read HTTP request from client
	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		fmt.Println("Error reading request:", err.Error())
		conn.Close()
		return
	}

	// Check if header contains "joy"
	if req.Header.Get("name") == "joy" {
		// Forward data to client1
		client1, err := net.Dial("tcp", "client1.example.com:8080")
		if err != nil {
			fmt.Println("Error connecting to client1:", err.Error())
			conn.Close()
			return
		}
		defer client1.Close()

		req.Write(client1)
	} else {
		// Check if request body contains "a": 123
		var body map[string]interface{}
		err := json.NewDecoder(req.Body).Decode(&body)
		if err == nil && body["a"] == 123 {
			// Forward data to client2
			client2, err := net.Dial("tcp", "client2.example.com:8080")
			if err != nil {
				fmt.Println("Error connecting to client2:", err.Error())
				return
			}
			defer client2.Close()
			req.Write(client2)
		}
	}
}
