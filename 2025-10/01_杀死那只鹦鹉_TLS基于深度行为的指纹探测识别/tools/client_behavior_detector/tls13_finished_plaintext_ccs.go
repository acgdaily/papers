package main

import (
	"crypto/tls"
	"net"
)

func acceptTLS13_Finished_Plaintext_ChangeCipherSpec(c net.Conn) {
	defer c.Close()

	tlsConn := tls.Server(c, tlsConfig)
	defer tlsConn.Close()

	err := tlsConn.Handshake()
	if err != nil {
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
