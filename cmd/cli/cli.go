package main

import (
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":7070")
	if err != nil {
		log.Fatalf("Error starting TCP server: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection2(conn)
	}
}

func handleConnection2(conn net.Conn) {
	defer conn.Close()

	sshConn, err := net.Dial("tcp", "127.0.0.1:22")
	if err != nil {
		log.Printf("Error connecting to local SSH port: %v", err)
		return
	}
	defer sshConn.Close()

	go io.Copy(conn, sshConn)
	io.Copy(sshConn, conn)
}
