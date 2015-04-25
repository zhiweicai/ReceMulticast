package main

import (
	"encoding/hex"
	"log"
	"net"
	"os"
)

const (
	maxDatagramSize = 8192
)

func main() {
	if len(os.Args) > 1 {
		service := os.Args[1]
		var sourceinterface string
		if len(os.Args) > 2 {
			sourceinterface = os.Args[2]
		}
		serveMulticastUDP(service, sourceinterface, msgHandler)
	}
}

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	log.Println(n, "bytes read from", src)
	log.Println(hex.Dump(b[:n]))
}

func serveMulticastUDP(a string, b string, h func(*net.UDPAddr, int, []byte)) {
	addr, err := net.ResolveUDPAddr("udp", a)
	verify(err)

	chooseif, err := net.InterfaceByName(b)
	verify(err)

	// ifaces, err := net.Interfaces()
	// verify(err)

	// for _, i := range ifaces {
	// 	addrs, err := i.Addrs()
	// 	verify(err)

	// 	for _, addr := range addrs {
	// 		fmt.Println(addr.String())
	// 		if addr.String() == b {
	// 			chooseif = &i
	// 			break
	// 		}
	// 	}

	// }

	l, err := net.ListenMulticastUDP("udp", chooseif, addr)
	l.SetReadBuffer(maxDatagramSize)
	for {
		b := make([]byte, maxDatagramSize)
		n, src, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		h(src, n, b)
	}
}

func verify(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
