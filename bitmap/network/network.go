package network

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"

	"github.com/Sirupsen/logrus"

	"github.com/vishvananda/netlink"
)

const (
	defaultNetworkPath = "/tmp/mynet/"
)

type Network struct {
	Name    string
	IpRange *net.IPNet
	Driver  string
}

type Endpoint struct {
	ID          string           `json:"id"`
	Device      netlink.Veth     `json:"dev"`
	IPAddress   net.IP           `json:"ip"`
	MacAddress  net.HardwareAddr `json:"mac"`
	PortMapping []string         `json:"portmapping"`
	Network     *Network
}

type NetworkDriver interface {
	Name() string
	Create(subnet string, name string) (*Network, error)
	Delete(network Network) error
	Connect(network *Network, endpoint *Endpoint) error
	Disconnect(network Network, endpoint *Endpoint) error
}

func Create(subnet, name string) (*Network, error) {
	ip, ipRange, err := net.ParseCIDR(subnet)
	ipRange.IP = ip
	n := &Network{
		Name:    name,
		IpRange: ipRange,
	}
	return n, err
}

//mydocker netowrk create --subnet 192.168.0.0/24 --driver bridge testbridge
func CreateNetwork(driver, subnet, name string) error {
	_, cidr, _ := net.ParseCIDR(subnet)
	gatewayIP, err := Allocate(cidr)
	if err != nil {
		return err
	}
	cidr.IP = gatewayIP

	// nw, err := drivers[driver].Create(cidr.String(), name)
	nw, err := Create(subnet, name)
	if err != nil {
		return err
	}
	nw.Driver = "bridge"
	return nw.dump(defaultNetworkPath)
}

func (nw *Network) dump(dumpPath string) error {
	if _, err := os.Stat(dumpPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dumpPath, 0755)
		} else {
			return err
		}
	}

	nwPath := path.Join(dumpPath, nw.Name)
	nwFile, err := os.OpenFile(nwPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logrus.Errorf("error :", err)
		return err
	}

	defer nwFile.Close()

	nwJson, err := json.Marshal(nw)
	if err != nil {
		logrus.Errorf("error:", err)
		return err
	}

	_, err = nwFile.Write(nwJson)
	if err != nil {
		logrus.Errorf("error: ", err)
		return err
	}
	return nil
}

func (nw *Network) Load(dumpPath string) error {
	fullpath := dumpPath + nw.Name
	nwConfigFile, err := os.Open(fullpath)
	defer nwConfigFile.Close()
	if err != nil {
		return err
	}
	nwJson := make([]byte, 2000)
	n, err := nwConfigFile.Read(nwJson)
	if err != nil {
		return err
	}
	err = json.Unmarshal(nwJson[:n], nw)
	if err != nil {
		logrus.Errorf("error load nw info", err)
		return err
	}

	fmt.Println(nw)
	fmt.Println(string(nwJson[:n]))
	return nil
}
