package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/ipv4"
)

var (
	bcast = net.IPv4(255, 255, 255, 255)
	dport uint16
	dests = []net.IP{
		net.IPv4(192, 168, 1, 1),
		net.IPv4(192, 168, 1, 2),
	}
)

func main() {
	// parse command line arguments
	var port = 6112
	flag.IntVar(&port, "p", port,
		"only forward packets with this destination `port`")
	flag.Parse()

	// make sure port is valid
	if port < 1 || port > 65535 {
		log.Fatal("invalid port")
	}
	dport = uint16(port)

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

	fmt.Printf("Forwarding broadcast packets with destination port %d\n",
		dport)

	// create packet buffer and start reading packets from raw socket
	buf := make([]byte, 2048)
	for {
		header, payload, controlMsg, err := raw.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		// only handle broadcast traffic
		if !header.Dst.Equal(bcast) {
			continue
		}

		// only handle traffic to configured udp destination port
		if binary.BigEndian.Uint16(payload[2:4]) != dport {
			continue
		}
		fmt.Println(header, payload, controlMsg)

		// forward packet to configured destination IPs
		for _, ip := range dests {
			// set new destination ip and send packet
			header.Dst = ip
			err = raw.WriteTo(header, payload, nil)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
