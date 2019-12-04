package network

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"path"
	"strings"

	"github.com/Sirupsen/logrus"
)

const ipamDefaultAllocatorPath = "/var/run/mydocker/network/ipam/subnet.json"

type IPAM struct {
	SubnetAllocatorPath string
	Subnet              *map[string]string
}

var ipAllocator = &IPAM{
	SubnetAllocatorPath: ipamDefaultAllocatorPath,
}

func (ipam *IPAM) load() error {
	// check the file store the network allocator infomation
	if _, err := os.Stat(ipam.SubnetAllocatorPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	// open and read network info
	subnetConfigFile, err := os.Open(ipam.SubnetAllocatorPath)
	defer subnetConfigFile.Close()
	if err != nil {
		return err
	}

	subnetJson := make([]byte, 1024)
	n, err := subnetConfigFile.Read(subnetJson)
	if err != nil {
		return err
	}

	err = json.Unmarshal(subnetJson[:n], ipam.Subnet)
	if err != nil {
		logrus.Errorf("error dump allocation info: %v", err)
		return err
	}
	return nil
}

func (ipam *IPAM) dump() error {
	ipamConfigFileDir, _ := path.Split(ipam.SubnetAllocatorPath)
	if _, err := os.Stat(ipamConfigFileDir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(ipamConfigFileDir, 0755)
		} else {
			return err
		}
	}

	// open or create new file
	subnetConfigFile, err := os.OpenFile(ipamDefaultAllocatorPath,
		os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	ipamConfigJson, err := json.Marshal(ipam.Subnet)
	if err != nil {
		return err
	}
	_, err = subnetConfigFile.Write(ipamConfigJson)
	if err != nil {
		return err
	}
	return nil
}

func (ipam *IPAM) Allocate(subnet *net.IPNet) (ip net.IP, err error) {
	ipam.Subnet = &map[string]string{}
	err := ipam.load()
	if err != nil {
		logrus.Errorf("error load allocation file:%v", err)
	}
	one, size := subnet.Mask.Size()

	// check if the net has been allocate
	if _, exist := (*ipam.Subnet)[subnet.String()]; !exist {
		(*ipam.Subnet)[subnet.String()] = strings.Repeat("0", 1<<uint8(size-one))
	}

	for c := range (*ipam.Subnet)[subnet.String()] {
		if (*ipam.Subnet)[subnet.String()][c] == '0' {
			ipalloc := []byte((*ipam.Subnet)[subnet.String()])
			ipalloc[c] = '1'
			(*ipam.Subnet)[subnet.String()] = string(ipalloc)
			ip = subnet.IP
		}
		for t := uint(4); t > 0; t -= 1 {
			[]byte(ip)[4-t] += uint8(c >> ((t - 1) * 8))
		}
		ip[3] += 1
		break
	}

	ipam.dump()
	return

}

func (ipam *IPAM) Release(subnet *net.IPNet, ipaddr *net.IP) error {
	ipam.Subnet = &map[string]string{}
	err := ipam.load()
	if err != nil {
		logrus.Errorf("error dump allocation info, %v", err)
	}

	c := 0
	releaseIP := ipaddr.To4()
	releaseIP[3] -= 1
	for t := uint(4); t > 0; t -= 1 {
		c += int(releaseIP[t-1]-subnet.IP[t-1]) << ((4 - t) * 8)
	}

	ipalloc := []byte(*ipam.Subnet)[subnet.String()]
	ipalloc[c] = '0'
	(*ipam.Subnet)[subnet.String()] = string(ipalloc)
	ipam.dump()
	return nil
}
