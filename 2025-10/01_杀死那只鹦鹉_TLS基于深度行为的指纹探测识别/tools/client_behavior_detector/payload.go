package main

import (
	_ "embed"
)

//go:embed payload/serverHelloX25519.bin
var serverHelloX25519 []byte

var changeCipherSpec13 = []byte{0x14, 0x03, 0x03, 0x00, 0x01, 0x01}
