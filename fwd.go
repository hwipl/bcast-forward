package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"golang.org/x/net/ipv4"
)

var (
	bcast = net.IPv4(255, 255, 255, 255)
	dport uint16
	dests []net.IP
)

func main() {
	// parse command line arguments
	var port = 6112
	var dest = ""
	flag.IntVar(&port, "p", port,
		"only forward packets with this destination `port`")
	flag.StringVar(&dest, "d", dest, "forward broadcast packets to "+
		"this comma-separated list of `IPs`, "+
		"e.g., \"192.168.1.1,192.168.1.2\"")
	flag.Parse()

	// make sure port is valid
	if port < 1 || port > 65535 {
		log.Fatal("invalid port")
	}
	dport = uint16(port)

	// make sure destination IPs are present and valid
	if dest == "" {
		log.Fatal("you must specify a destination IP")
	}
	for _, d := range strings.Split(dest, ",") {
		dests = append(dests, net.ParseIP(d))
	}

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

	fmt.Printf("Forwarding broadcast packets with destination port %d "+
		"to IPs %v\n", dport, dests)

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
