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
	srcIP net.IP
	dests []net.IP
)

// parse_command_line parses the command line arguments
func parse_command_line() {
	var port = 6112
	var dest = ""
	var src = ""

	// set command line arguments
	flag.IntVar(&port, "p", port,
		"only forward packets with this destination `port`")
	flag.StringVar(&src, "s", src, "rewrite source address to this `IP`")
	flag.StringVar(&dest, "d", dest, "forward broadcast packets to "+
		"this comma-separated list of `IPs`, "+
		"e.g., \"192.168.1.1,192.168.1.2\"")
	flag.Parse()

	// make sure port is valid
	if port < 0 || port > 65535 {
		log.Fatal("invalid port")
	}
	dport = uint16(port)

	// parse source IP
	if src != "" {
		srcIP = net.ParseIP(src).To4()
		if srcIP == nil {
			log.Fatal("invalid source IP: ", src)
		}
	}

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
}

// print_info prints an info/settings header to the console
func print_info() {
	sep := strings.Repeat("-", 70)
	port := fmt.Sprintf("%d", dport)
	if dport == 0 {
		port = "any"
	}

	pFmt := "Receiving broadcast packets with destination port:    %s\n"
	dFmt := "Forwarding packets to IP:                             %s\n"
	sFmt := "Rewriting source address to IP:                       %s\n"

	fmt.Println(sep)
	fmt.Printf(pFmt, port)
	for _, ip := range dests {
		fmt.Printf(dFmt, ip)
	}
	if srcIP != nil {
		fmt.Printf(sFmt, srcIP)
	}
	fmt.Println(sep)
}

// run_socket_loop runs the main socket loop, reading packets from the socket
// and forwarding them to destination ip addresses
func run_socket_loop() {
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
	print_info()

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
		for _, ip := range dests {
			// set new source and destination ip and send packet
			if srcIP != nil {
				header.Src = srcIP
			}
			header.Dst = ip
			err = raw.WriteTo(header, payload, nil)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func main() {
	parse_command_line()
	run_socket_loop()
}
