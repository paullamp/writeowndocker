package main

import (
	"bitmap/network"
	"fmt"
)

func main() {
	network.CreateNetwork("bridge", "192.168.0.0/16", "mydocker0")
	nw := network.Network{
		Name: "mydocker0",
	}
	// nw.load("/tmp/mynet/")
	fmt.Println(nw)
	nw.Load("/tmp/mynet/")

}
