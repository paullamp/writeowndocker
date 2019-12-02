package network

import (
	"net"
)

func Allocate(cidr *net.IPNet) (net.IP, error) {
	return cidr.IP, nil
}
