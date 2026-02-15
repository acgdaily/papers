package main

import (
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var peer string
var ccs = []byte{0x14, 0x03, 0x03, 0x00, 0x01, 0x01}

func main() {
	flag.StringVar(&peer, "peer", "", "The server address (ip:port)")

	help := flag.Bool("h", false, "Show help")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	clientHello, err := os.ReadFile("clientHello.bin")
	if err != nil {
		fmt.Printf("load clientHello failed: %s\n", err)
		return
	}

	if len(clientHello) < 11+32 {
		fmt.Printf("unsupported clientHello, random not exists\n")
		return
	}

	replayBehavior, delay, err := replay(clientHello)
	if err != nil {
		fmt.Printf("replay failed: %s\n", err)
		return
	}

	if behavior := behavior(replayBehavior); behavior != "" {
		fmt.Printf("replay: peer server probably %s\n", behavior)
	}

	_, err = io.ReadFull(rand.Reader, clientHello[11:11+32])
	if err != nil {
		fmt.Printf("fill new random field to clientHello failed: %s\n", err)
		return
	}

	fetchBehavior, err := fetch(clientHello, delay)
	if err != nil {
		fmt.Printf("fetch failed: %s\n", err)
		return
	}

	if behavior := behavior(fetchBehavior); behavior != "" {
		fmt.Printf("fetch: peer server probably %s\n", behavior)
	}

	if replayBehavior != fetchBehavior {
		fmt.Printf("behavior mismatch: replay: %d, fetch: %d\n", replayBehavior, fetchBehavior)
		return
	}

	fmt.Printf("behavior matched: replay: %d, fetch: %d\n", replayBehavior, fetchBehavior)
}

func behavior(count int) string {
	switch count {
	case -1:
		return "gnutls"
	case 16:
		return "golang"
	case 32:
		return "openssl/boringssl"
	default:
		return ""
	}
}

func replay(clientHello []byte) (int, time.Duration, error) {
	c, err := net.Dial("tcp", peer)
	if err != nil {
		return 0, 0, err
	}
	defer c.Close()

	sentTime := time.Now()
	c.Write(clientHello) // send clientHello

	c.Read(make([]byte, 1))
	rtt := time.Since(sentTime) // calc RTT
	fmt.Printf("rtt %d ms\n", rtt.Milliseconds())

	if rtt > 500*time.Millisecond {
		return 0, 0, errors.New("rtt over 500ms maybe jitter, try again")
	}

	delay := rtt + 500*time.Millisecond // +500ms jitter

	go copy(io.Discard, c) // response TCP-FIN

	for i := 1; i <= 100; i++ {
		time.Sleep(delay)

		_, err = c.Write(ccs)
		if err != nil {
			return i - 2, delay, nil // write on i error means peer received i-1 records is over max value, max value is i-2
		}
	}

	return -1, rtt, nil
}

func fetch(clientHello []byte, delay time.Duration) (int, error) {
	c, err := net.Dial("tcp", peer)
	if err != nil {
		return 0, err
	}
	defer c.Close()

	c.Write(clientHello) // send clientHello

	go copy(io.Discard, c) // response TCP-FIN

	for i := 1; i <= 100; i++ {
		time.Sleep(delay)

		_, err = c.Write(ccs)
		if err != nil {
			return i - 2, nil // write on i error means peer received i-1 records is over max value, max value is i-2
		}
	}

	return -1, nil
}

func copy(dst io.Writer, src net.Conn) {
	defer src.Close()
	io.Copy(dst, src)
}
