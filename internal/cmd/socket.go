package cmd

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/ipv4"
)

var (
	bcast = net.IPv4(255, 255, 255, 255)
)

// runSocketLoop runs the main socket loop, reading packets from the socket
// and forwarding them to destination ip addresses
func runSocketLoop() {
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

	// print some info before entering main loop
	printInfo()

	// create packet buffer and start reading packets from raw socket
	buf := make([]byte, 2048)
	for {
		header, payload, _, err := raw.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		// only handle broadcast traffic
		if !header.Dst.Equal(bcast) {
			continue
		}

		// only handle traffic to configured udp destination port
		destPort := binary.BigEndian.Uint16(payload[2:4])
		if dport > 0 && destPort != dport {
			continue
		}

		srcPort := binary.BigEndian.Uint16(payload[0:2])
		fmt.Printf("Got packet: %s:%d -> %s:%d\n", header.Src,
			srcPort, header.Dst, destPort)

		// remove udp header checksum in forwarded packets
		binary.BigEndian.PutUint16(payload[6:8], 0)

		// forward packet to configured destination IPs
		for _, d := range dests {
			// set new source and destination ip and send packet
			if srcIP != nil {
				header.Src = srcIP
			}
			header.Dst = d.ip
			err = raw.WriteTo(header, payload, nil)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
