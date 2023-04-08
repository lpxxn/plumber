package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

func TestSSHCopy(t *testing.T) {
	localEndpoint := &Endpoint{
		Host: "localhost",
		Port: 7900,
	}

	serverEndpoint := &Endpoint{
		Host: "localhost",
		Port: 6700,
	}

	remoteEndpoint := &Endpoint{
		Host: "172.17.0.5",
		Port: 22,
	}

	// Read the private key file
	privateKeyBytes, err := os.ReadFile("test.pem")
	if err != nil {
		log.Fatalf("Failed to read private key: %s", err)
	}

	// Parse the private key
	privateKey, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %s", err)
	}

	// Initialize SSH client configuration with the private key
	sshConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privateKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	tunnel := &SSHtunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}

	tunnel.Start()
}

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

type SSHtunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint

	Config *ssh.ClientConfig
}

func (tunnel *SSHtunnel) Start() error {
	fmt.Printf("Starting tunnel %s -> %s -> %s", tunnel.Local, tunnel.Server, tunnel.Remote)
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go tunnel.forward(conn)
	}
}

func (tunnel *SSHtunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		fmt.Printf("Server dial error: %s\n", err)
		return
	}

	//remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		fmt.Printf("Remote dial error: %s\n", err)
		return
	}

	// Set keepalive parameters
	if tcpConn, ok := remoteConn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(5 * time.Minute)
	}

	copyConn := func(writer, reader net.Conn) {
		defer writer.Close()
		defer reader.Close()

		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

func TestSSH2(t *testing.T) {
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
	dialer := net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	sshConn, err := dialer.Dial("tcp", "127.0.0.1:22")
	if err != nil {
		log.Printf("Error connecting to local SSH port: %v", err)
		return
	}

	tcpConn, ok := conn.(*net.TCPConn)
	if ok {
		tcpConn.SetKeepAlivePeriod(60 * time.Second)
	}

	defer sshConn.Close()

	go io.Copy(conn, sshConn)
	io.Copy(sshConn, conn)
}

func TestSSH3(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:7070")
	if err != nil {
		log.Printf("Error connecting to local SSH port: %v", err)
		return
	}
	defer conn.Close()
	// Read messages from the client
	go func() {
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
	}()
	if _, err := conn.Write([]byte("ls -l\r\n")); err != nil {
		log.Printf("Error writing to SSH port: %v", err)
		return
	}
	time.Sleep(5 * time.Second)
}
