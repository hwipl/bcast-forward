# bcast-forward

bcast-forward is a Linux command line tool that forwards UDP broadcast packets
with destination IP `255.255.255.255` to a specified list of unicast addresses.
For example, it can be used for playing old LAN games, that use broadcasts to
discover game servers, over a VPN tunnel.

## Installation

You can download and install bcast-forward with its dependencies to your GOPATH
or GOBIN with the go tool:

```console
$ go get github.com/hwipl/bcast-forward/cmd/bcast-forward
```

## Usage

You can run `bcast-forward` with the following command line arguments:

```
  -d IPs
	forward broadcast packets to this comma-separated list of IPs, e.g.,
        "192.168.1.1,192.168.1.2"
  -p port
        only forward packets with this destination port (default 6112)
  -s IP
        rewrite source address to this IP
```
