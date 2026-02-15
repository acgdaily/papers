package main

import (
	"crypto/tls"
	"log"
)

func testTLS13_Encrypted_FatalAlert_UserCancel() {
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

	var data [2]byte
	data[0] = alertLevelError
	data[1] = byte(alertUserCanceled)

	for range 100 {
		delay()

		_, err = writeRecord(tlsConn, recordTypeAlert, data[:])
		if err != nil {
			return
		}
	}
}
