package main

import (
	"crypto/tls"
	"net"
)

func acceptTLS13_Encrypted_Alert_FatalAsWarning(c net.Conn) {
	defer c.Close()

	tlsConn := tls.Server(c, tlsConfig)
	defer tlsConn.Close()

	err := tlsConn.Handshake()
	if err != nil {
		return
	}

	var data [2]byte
	data[0] = alertLevelWarning
	data[1] = byte(alertAccessDenied)

	for range 100 {
		delay()

		_, err = writeRecord(tlsConn, recordTypeAlert, data[:])
		if err != nil {
			return
		}
	}
}
