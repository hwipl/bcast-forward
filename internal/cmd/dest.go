package cmd

import (
	"fmt"
	"log"
	"net"
)

// dest stores information about a forwarding destination
type dest struct {
	ip    net.IP
	srcIP net.IP
}

// getSourceIP gets the source IP used for the forwarding destination
func (d *dest) getSourceIP() {
	// create dummy connection to retrieve local address
	addr := fmt.Sprintf("%s:%d", d.ip, dport)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// get local address of connection
	laddr := conn.LocalAddr().(*net.UDPAddr)
	d.srcIP = laddr.IP
}

// newDest creates and returns a new dest
func newDest(addr string) *dest {
	var dest dest

	// parse IP address
	ip := net.ParseIP(addr).To4()
	if ip == nil {
		// invalid IP, stop here
		return nil
	}
	dest.ip = ip

	// get source IP
	dest.getSourceIP()

	return &dest
}
