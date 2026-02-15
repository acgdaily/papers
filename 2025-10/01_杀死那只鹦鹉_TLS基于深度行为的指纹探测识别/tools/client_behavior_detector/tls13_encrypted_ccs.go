package main

import (
	"crypto/tls"
	"net"
)

func acceptTLS13_Encrypted_ChangeCipherSpec(c net.Conn) {
	defer c.Close()

	tlsConn := tls.Server(c, tlsConfig)
	defer tlsConn.Close()

	err := tlsConn.Handshake()
	if err != nil {
		return
	}

	for range 100 {
		delay()

		err = sendChangeCipherSpec(tlsConn)
		if err != nil {
			return
		}
	}
}
