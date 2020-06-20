package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// open raw socket
	conn, err := net.ListenPacket("ip4:udp", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// create packet buffer and start reading packets from raw socket
	buf := make([]byte, 2048)
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(n, addr, buf[:n])
	}
}
