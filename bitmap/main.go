package main

import (
	"fmt"
	"net"
	"strings"
)

// type Network struct {
// 	net.IPMask
// }

func main() {
	str := strings.Repeat("0", 1<<8)
	fmt.Println(str)
	slic := []byte(str)
	slic[0] = '1'
	slic[1] = '1'

	fmt.Println(len(str))
	var key int
	for key = range slic {
		if slic[key] == '0' {
			fmt.Println(key)
			slic[key] = '1'
			break
		}
	}
	fmt.Println("The can use ip is :", key)
	fmt.Println(string(slic))

	ip, ipnet, err := net.ParseCIDR("192.168.0.18/24")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("192.168.0.18/24 's ip is : ", ip)
	fmt.Println("192.168.0.18/24 's net is :", ipnet)
	fmt.Println()
	fmt.Println(ipnet.Mask)
	fmt.Println(len(ipnet.Mask))
	ip_mask := net.IPv4Mask(1, 1, 1, 1)
	fmt.Println(ip_mask)
}
