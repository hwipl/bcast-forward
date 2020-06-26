package cmd

import "net"

// dest stores information about a forwarding destination
type dest struct {
	ip net.IP
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

	return &dest
}
