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
	if port < 0 || port > 65535 {
		log.Fatal("invalid port")
	}
	dport = uint16(port)

	// make sure destination IPs are present and valid
	if dest == "" {
		log.Fatal("you must specify a destination IP")
	}
	for _, d := range strings.Split(dest, ",") {
		if d == "" {
			continue
		}
		ip := net.ParseIP(d).To4()
		if ip == nil {
			log.Fatal("invalid destination IP: ", d)
		}
		dests = append(dests, ip)
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

	// print some info before entering main loop
	fmt.Printf("Receiving broadcast packets with destination port %d.\n",
		dport)
	fmt.Printf("Forwarding received packets to IPs %s\n", dests)

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
