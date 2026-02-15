package main

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	_ "unsafe"
)

type recordType uint8

const (
	recordTypeChangeCipherSpec recordType = 20
	recordTypeAlert            recordType = 21
	recordTypeHandshake        recordType = 22
	recordTypeApplicationData  recordType = 23
)

type alert uint8

const (
	// alert level
	alertLevelWarning = 1
	alertLevelError   = 2
)

const (
	alertCloseNotify                  alert = 0
	alertUnexpectedMessage            alert = 10
	alertBadRecordMAC                 alert = 20
	alertDecryptionFailed             alert = 21
	alertRecordOverflow               alert = 22
	alertDecompressionFailure         alert = 30
	alertHandshakeFailure             alert = 40
	alertBadCertificate               alert = 42
	alertUnsupportedCertificate       alert = 43
	alertCertificateRevoked           alert = 44
	alertCertificateExpired           alert = 45
	alertCertificateUnknown           alert = 46
	alertIllegalParameter             alert = 47
	alertUnknownCA                    alert = 48
	alertAccessDenied                 alert = 49
	alertDecodeError                  alert = 50
	alertDecryptError                 alert = 51
	alertExportRestriction            alert = 60
	alertProtocolVersion              alert = 70
	alertInsufficientSecurity         alert = 71
	alertInternalError                alert = 80
	alertInappropriateFallback        alert = 86
	alertUserCanceled                 alert = 90
	alertNoRenegotiation              alert = 100
	alertMissingExtension             alert = 109
	alertUnsupportedExtension         alert = 110
	alertCertificateUnobtainable      alert = 111
	alertUnrecognizedName             alert = 112
	alertBadCertificateStatusResponse alert = 113
	alertBadCertificateHashValue      alert = 114
	alertUnknownPSKIdentity           alert = 115
	alertCertificateRequired          alert = 116
	alertNoApplicationProtocol        alert = 120
	alertECHRequired                  alert = 121
)

func readRecord(c net.Conn) (recordType, []byte, error) {
	var hdr [5]byte

	_, err := io.ReadFull(c, hdr[:])
	if err != nil {
		return 0, nil, err
	}

	typ := recordType(hdr[0])
	switch typ {
	case recordTypeAlert, recordTypeHandshake, recordTypeChangeCipherSpec, recordTypeApplicationData:
	default:
		return typ, nil, fmt.Errorf("unexpected record type %#x", typ)
	}

	version := binary.BigEndian.Uint16(hdr[1:3])
	if version < tls.VersionTLS10 || version > tls.VersionTLS13 {
		return typ, nil, fmt.Errorf("unexpected version %d", version)
	}

	length := binary.BigEndian.Uint16(hdr[3:5])

	payload := make([]byte, length)
	_, err = io.ReadFull(c, payload)
	if err != nil {
		return typ, nil, err
	}

	return typ, payload, nil
}

//go:linkname sendChangeCipherSpec crypto/tls.(*Conn).writeChangeCipherRecord
func sendChangeCipherSpec(c *tls.Conn) error

//go:linkname writeRecord crypto/tls.(*Conn).writeRecordLocked
func writeRecord(c *tls.Conn, typ recordType, data []byte) (int, error)
