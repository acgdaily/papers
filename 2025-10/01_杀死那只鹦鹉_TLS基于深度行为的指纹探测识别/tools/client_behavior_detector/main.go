package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type acceptHandler func(net.Conn)

var (
	tlsConfig = &tls.Config{}
)

func init() {
	logger, err := os.OpenFile("ssl.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("open 'ssl.log' failed: %s\n", err)
		return
	}

	tlsConfig.KeyLogWriter = logger
}

func main() {
	var cert, key string
	flag.StringVar(&cert, "cert", "cert.pem", "The server certificates")
	flag.StringVar(&key, "key", "key.pem", "The server private key")

	var beginPort uint
	flag.UintVar(&beginPort, "begin-port", 10000, "Listening start port")

	help := flag.Bool("h", false, "Show help")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	acceptHandlers := []struct {
		name    string
		handler acceptHandler
	}{
		{
			name:    "TLS Alert BadProtocolVersion",
			handler: acceptTLS_Alert_BadProtocolVersion,
		},
		{
			name:    "TLS WarningAlert",
			handler: acceptTLS_WarningAlert,
		},
		{
			name:    "TLSv1.3 Plaintext ChangeCiperSpec",
			handler: acceptTLS13_Plaintext_ChangeCipherSpec,
		},
		{
			name:    "TLSv1.3 Finished Plaintext ChangeCiperSpec",
			handler: acceptTLS13_Finished_Plaintext_ChangeCipherSpec,
		},
		{
			name:    "TLSv1.3 Encrypted ChangeCiperSpec",
			handler: acceptTLS13_Encrypted_ChangeCipherSpec,
		},
		{
			name:    "TLSv1.3 Encrypted Alert FatalAsWarning",
			handler: acceptTLS13_Encrypted_Alert_FatalAsWarning,
		},
		{
			name:    "TLSv1.3 Encrypted WarningAlert UserCancel",
			handler: acceptTLS13_Encrypted_WarningAlert_UserCancel,
		},
		{
			name:    "TLSv1.3 Encrypted FatalAlert UserCancel",
			handler: acceptTLS13_Encrypted_FatalAlert_UserCancel,
		},
		{
			name:    "TLSv1.3 Encrypted WarningAlert NoRenegotiation",
			handler: acceptTLS13_Encrypted_WarningAlert_NoRenegotiation,
		},
	}

	if int(beginPort) > 65535-len(acceptHandlers) {
		log.Fatalf("begin port too high!\n")
		return
	}

	tlsCert, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Fatalf("load certificate failed: %s\n", err)
		return
	}

	tlsConfig.Certificates = []tls.Certificate{tlsCert}

	port := beginPort
	for _, acceptHandler := range acceptHandlers {
	Listen:
		if port > 65535 {
			log.Fatalf("no more port available!\n")
			return
		}

		ln, err := net.ListenTCP("tcp", &net.TCPAddr{
			Port: int(port),
		})
		if err != nil {
			port++
			goto Listen
		}

		log.Printf("Test <%s> listening %s\n", acceptHandler.name, ln.Addr())
		go acceptLoop(ln, acceptHandler.handler)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Signal received, exiting")
}

func acceptLoop(ln net.Listener, handler acceptHandler) {
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Temporary() {
				continue
			}
			log.Printf("accept %s failed: %s\n", ln.Addr(), err)
			return
		}

		go handler(conn)
	}
}

func delay() {
	time.Sleep(time.Second / 4)
}
