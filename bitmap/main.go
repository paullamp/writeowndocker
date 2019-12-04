package main

import (
	"bitmap/network"
	// "fmt"
)

func main() {
	network.CreateNetwork("bridge", "192.168.0.0/16", "mydocker0")
	network.CreateNetwork("bridge", "10.1.1.0/24", "mydocker1")
	network.CreateNetwork("bridge", "172.16.0.0/24", "mydcoker172")
	// nw := network.Network{
	// 	Name: "mydocker0",
	// }
	// // nw.load("/tmp/mynet/")
	// fmt.Println(nw)
	// nw.Load("/tmp/mynet/")
	network.Init()
	network.ListNetwork()
	network.DeleteNetwork("mydocker1")
	network.Init()
	network.ListNetwork()

}
