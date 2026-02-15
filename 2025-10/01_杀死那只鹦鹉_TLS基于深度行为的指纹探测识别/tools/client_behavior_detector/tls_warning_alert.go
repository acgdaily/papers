package main

import (
	"crypto/tls"
	"encoding/binary"
	"net"
)

func acceptTLS_WarningAlert(c net.Conn) {
	defer c.Close()

	typ, _, err := readRecord(c)
	if err != nil {
		return
	}

	if typ != recordTypeHandshake {
		return
	}

	var data [7]byte
	data[0] = byte(recordTypeAlert)
	binary.BigEndian.PutUint16(data[1:3], tls.VersionTLS10)
	binary.BigEndian.PutUint16(data[3:5], 2)
	data[5] = alertLevelWarning
	data[6] = byte(alertUserCanceled)

	for range 100 {
		delay()

		_, err := c.Write(data[:])
		if err != nil {
			return
		}
	}
}
