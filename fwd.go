package main

import (
	"fmt"
	"golang.org/x/net/ipv4"
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

	raw, err := ipv4.NewRawConn(conn)
	if err != nil {
		log.Fatal(err)
	}

	// create packet buffer and start reading packets from raw socket
	buf := make([]byte, 2048)
	for {
		header, payload, controlMsg, err := raw.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(header, payload, controlMsg)
	}
}
