### Go TLS Client-Server test
*Use both server and client authentication*

Generate root certificate:
```
$ go run gencert.go -type root
2014/10/25 14:27:07 certificate: root_cert.pem
2014/10/25 14:27:07 private key: root_key.pem
2014/10/25 14:27:07 dnsname:  root.testapplication
```

Generate server certificate signed by root:
```
$ go run gencert.go -type server
2014/10/25 14:27:21 loading root_cert.pem and root_key.pem
2014/10/25 14:27:21 certificate: server_cert.pem
2014/10/25 14:27:21 private key: server_key.pem
2014/10/25 14:27:21 dnsname:  server.testapplication
```

Generate client certificate signed by root:
```
$ go run gencert.go -type client -client-id=12345
2014/10/25 14:27:58 loading root_cert.pem and root_key.pem
2014/10/25 14:27:58 certificate: client_12345_cert.pem
2014/10/25 14:27:58 private key: client_12345_key.pem
2014/10/25 14:27:58 dnsname:  12345.client.testapplication
```

Server output:
```
$ go run server.go
2014/10/25 15:00:15 server started
2014/10/25 15:00:20 handle connection from client-id: 12345
2014/10/25 15:00:20 client:12345: read 13 bytes
2014/10/25 15:00:20 client:12345: echo: hello-goodbye
2014/10/25 15:00:20 client:12345: wrote 13 bytes
2014/10/25 15:00:20 client:12345: connection closed
```

Client output:
```
$ go run client.go
2014/10/25 15:00:20 handshake: true
2014/10/25 15:00:20 sending: hello-goodbye
2014/10/25 15:00:20 wrote 13 bytes
2014/10/25 15:00:20 read 13 bytes
2014/10/25 15:00:20 server reply: hello-goodbye
```