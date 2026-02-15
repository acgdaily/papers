package main

import (
	"crypto/tls"
	"log"
)

func testTLS13_Finished_Plaintext_ChangeCipherSpec() {
	c, err := dial()
	if err != nil {
		log.Printf("dial peer failed, error: %s\n", err)
		return
	}
	defer c.Close()

	tlsConn := tls.Client(c, tlsConfig)
	defer tlsConn.Close()

	err = tlsConn.Handshake()
	if err != nil {
		log.Printf("handshake failed, error: %s\n", err)
		return
	}

	if tlsConn.ConnectionState().Version != tls.VersionTLS13 {
		log.Printf("peer unsupported TLS v1.3\n")
		return
	}

	for range 100 {
		delay()

		_, err = c.Write(changeCipherSpec13)
		if err != nil {
			return
		}
	}
}
