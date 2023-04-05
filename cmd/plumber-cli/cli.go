package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {

	// exit signal
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	<-ch
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
