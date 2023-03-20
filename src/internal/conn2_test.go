package internal

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"testing"

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
