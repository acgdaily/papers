package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"os"
	"time"
)

var (
	peer      string
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
	flag.StringVar(&peer, "peer", "", "The server address (ip:port)")

	help := flag.Bool("h", false, "Show help")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	host, _, err := net.SplitHostPort(peer)
	if err != nil {
		log.Printf("parse peer failed, error: %s\n", err)
		return
	}
	tlsConfig.ServerName = host

	tests := []struct {
		name    string
		handler func()
	}{
		{
			name:    "TLS WarningAlert",
			handler: testTLS_WarningAlert,
		},
		{
			name:    "TLSv1.3 Plaintext ChangeCiperSpec",
			handler: testTLS13_Plaintext_ChangeCipherSpec,
		},
		{
			name:    "TLSv1.3 Finished Plaintext ChangeCiperSpec",
			handler: testTLS13_Finished_Plaintext_ChangeCipherSpec,
		},
		{
			name:    "TLSv1.3 Encrypted ChangeCiperSpec",
			handler: testTLS13_Encrypted_ChangeCipherSpec,
		},
		{
			name:    "TLSv1.3 Encrypted Alert FatalAsWarning",
			handler: testTLS13_Encrypted_Alert_FatalAsWarning,
		},
		{
			name:    "TLSv1.3 Encrypted WarningAlert UserCancel",
			handler: testTLS13_Encrypted_WarningAlert_UserCancel,
		},
		{
			name:    "TLSv1.3 Encrypted FatalAlert UserCancel",
			handler: testTLS13_Encrypted_FatalAlert_UserCancel,
		},
		{
			name:    "TLSv1.3 Encrypted WarningAlert NoRenegotiation",
			handler: testTLS13_Encrypted_WarningAlert_NoRenegotiation,
		},
	}

	for _, test := range tests {
		log.Printf("Test <%s> begin\n", test.name)
		test.handler()
		log.Printf("Test <%s> finished\n", test.name)

		time.Sleep(5 * time.Second)
	}

	log.Println("Test finished, exiting")
}

func dial() (net.Conn, error) {
	return net.Dial("tcp", peer)
}

func delay() {
	time.Sleep(time.Second / 4)
}
