package main

import (
	"crypto/tls"
	"errors"
	"log"
	"net"
)

type plaintextChangeCipherSpecConn struct {
	net.Conn
	tlsConn              *tls.Conn
	clientHelloSent      bool
	changeCipherSpecSent bool
}

func (c *plaintextChangeCipherSpecConn) Write(b []byte) (int, error) {
	if !c.clientHelloSent {
		c.clientHelloSent = true
		return c.Conn.Write(b)
	}

	if connectionState(c.tlsConn).Version != tls.VersionTLS13 {
		return 0, errors.New("peer unsupported TLS v1.3")
	}

	if !c.changeCipherSpecSent {
		c.changeCipherSpecSent = true

		for range 100 {
			delay()

			_, err := c.Write(changeCipherSpec13)
			if err != nil {
				return 0, err
			}
		}
	}

	return c.Conn.Write(b)
}

func testTLS13_Plaintext_ChangeCipherSpec() {
	c, err := dial()
	if err != nil {
		log.Printf("dial peer failed, error: %s\n", err)
		return
	}
	defer c.Close()

	wrappedConn := &plaintextChangeCipherSpecConn{Conn: c}

	tlsConn := tls.Client(wrappedConn, tlsConfig)
	defer tlsConn.Close()

	wrappedConn.tlsConn = tlsConn

	err = tlsConn.Handshake()
	if err != nil {
		if !wrappedConn.changeCipherSpecSent {
			log.Printf("handshake failed, error: %s\n", err)
		}
		return
	}
}
