package cmd

import (
	"flag"
	"log"
	"net"
	"strings"
)

// global variables set via command line arguments:
var (
	// dport is the destination port for packet matching
	dport uint16

	// srcIP is the source IP used for source IP rewriting
	srcIP net.IP

	// dests is the list of IPs to forward the packets to
	dests []*dest
)

// parseCommandLine parses the command line arguments
func parseCommandLine() {
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
		dst := newDest(d)
		if dst == nil {
			log.Fatal("invalid destination IP: ", d)
		}
		dests = append(dests, dst)
	}
}

// Run is the main entry point
func Run() {
	parseCommandLine()
	runSocketLoop()
}
