package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

const appdomain = "testapplication"

func main() {
	// load root certificate to verify client certificate
	rootPEM, err := ioutil.ReadFile("root_cert.pem")
	if err != nil {
		log.Fatalf("failed to read root certificate: %s", err)
	}
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		log.Fatalf("failed to parse root certificate")
	}

	// load server certificate
	cert, err := tls.LoadX509KeyPair("server_cert.pem", "server_key.pem")
	if err != nil {
		log.Fatalf("failed to load server tls certificate: %s", err)
	}

	config := tls.Config{
		ClientCAs:    roots,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", "127.0.0.1:4444", &config)
	if err != nil {
		log.Fatalf("error: listen: %s", err)
	}
	log.Print("server started")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: accept: %s", err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		log.Printf("error: not tls conn")
		return
	}

	err := tlsConn.Handshake()
	if err != nil {
		log.Printf("error: handshake: %s", err)
		return
	}

	clientID, err := getClientID(tlsConn)
	if err != nil {
		log.Printf("error: cannot get client-id: %s", err)
		return
	}

	log.Printf("handle connection from client-id: %s", clientID)

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("client:%s: error: read: %s", clientID, err)
			}
			break
		}
		log.Printf("client:%s: read %d bytes", clientID, n)

		log.Printf("client:%s: echo: %s", clientID, string(buf[:n]))

		n, err = conn.Write(buf[:n])
		if err != nil {
			log.Printf("client:%s: error: write: %s", clientID, err)
			break
		}
		log.Printf("client:%s: wrote %d bytes", clientID, n)
	}
	log.Printf("client:%s: connection closed", clientID)
}

func getClientID(tlsConn *tls.Conn) (string, error) {
	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return "", fmt.Errorf("client certificate not found")
	}
	cert := state.PeerCertificates[0]

	if len(cert.DNSNames) == 0 {
		return "", fmt.Errorf("client certificate dns name not found")
	}
	clientDNSName := cert.DNSNames[0]

	parts := strings.Split(clientDNSName, ".")
	if len(parts) != 3 || len(parts[0]) == 0 || parts[1] != "client" || parts[2] != appdomain {
		return "", fmt.Errorf("bad client dns name")
	}

	return parts[0], nil
}
