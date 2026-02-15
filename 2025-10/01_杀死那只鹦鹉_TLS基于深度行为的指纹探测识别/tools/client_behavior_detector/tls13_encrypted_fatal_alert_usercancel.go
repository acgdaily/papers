package main

import (
	"crypto/tls"
	"net"
)

func acceptTLS13_Encrypted_FatalAlert_UserCancel(c net.Conn) {
	defer c.Close()

	tlsConn := tls.Server(c, tlsConfig)
	defer tlsConn.Close()

	err := tlsConn.Handshake()
	if err != nil {
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
