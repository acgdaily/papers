package main

import (
	"bytes"
	"net"
)

func acceptTLS13_Plaintext_ChangeCipherSpec(c net.Conn) {
	defer c.Close()

	typ, clientHello, err := readRecord(c)
	if err != nil {
		return
	}

	if typ != recordTypeHandshake {
		return
	}

	if len(clientHello) < 39+32 || clientHello[38] != 32 {
		return
	}

	serverHelloX25519 := bytes.Clone(serverHelloX25519)
	copy(serverHelloX25519[44:44+32], clientHello[39:39+32]) // legacy session id

	c.Write(serverHelloX25519)

	for range 100 {
		delay()

		_, err = c.Write(changeCipherSpec13)
		if err != nil {
			return
		}
	}
}
