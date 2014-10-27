package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
)

const appname = "testapp"

func main() {
	// load root certificate to verify server certificate
	rootPEM, err := ioutil.ReadFile("root_cert.pem")
	if err != nil {
		log.Fatalf("failed to read root certificate: %s", err)
	}
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		panic("failed to parse root certificate")
	}

	// load client certificate
	cert, err := tls.LoadX509KeyPair("client_1_cert.pem", "client_1_key.pem")
	if err != nil {
		log.Fatalf("failed to load client tls certificate: %s", err)
	}

	config := tls.Config{
		RootCAs:      roots,
		ServerName:   appname + "-server",
		Certificates: []tls.Certificate{cert},
	}

	conn, err := tls.Dial("tcp", "127.0.0.1:4444", &config)
	if err != nil {
		log.Fatalf("error: dial: %s", err)
	}
	defer conn.Close()

	message := "hello-goodbye"
	log.Printf("sending: %s", message)

	n, err := conn.Write([]byte(message))
	if err != nil {
		log.Fatalf("error: write: %s", err)
	}
	log.Printf("wrote %d bytes", n)

	buf := make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		log.Fatalf("error: read: %s", err)
	}
	log.Printf("read %d bytes", n)

	log.Printf("server reply: %s", string(buf[:n]))
}
