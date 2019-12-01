package network

import (
	"encoding/json"
	"net"
	"os"
	"path"

	"github.com/Sirupsen/logrus"

	"github.com/vishvananda/netlink"
)

var (
	defaultNetworkPath = "/var/run/mydocker/network/network/"
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
	Disconnect(network *Network, endpoint *Endpoint) error
}

//mydocker network create --subnet 192.168.1.0/24 --driver bridge testbridge
//all the strings was transfored from command line
func CreateNetwork(driver, subnet, name string) error {
	_, cidr, _ := net.ParseCIDR(subnet)
	gatewayIP, err := ipAllocator.Allocate(cidr)
	if err != nil {
		return err
	}
	cidr.IP = gatewayIP

	nw, err := drivers[driver].Create(cidr.String(), name)
	if err != nil {
		return err
	}
	return nw.dump(defaultNetworkPath)
}

func (nw *Network) dump(dumpPath string) error {
	if _, err := os.Stat(dumpPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dumpPath, 0644)
		} else {
			return err
		}
	}

	// save network name as filename
	nwPath := path.Join(dumpPath, nw.Name)
	nwFile, err := os.OpenFile(nwPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 06444)
	if err != nil {
		logrus.Errorf("error: ", err)
		return err
	}

	defer nwFile.Close()

	//store json to file
	nwJson, err := json.Marshal(nw)
	if err != nil {
		logrus.Errorf("error: ", err)
		return err
	}

	_, err := nwFile.Write(nwJson)
	if err != nil {
		logrus.Errorf("error:", err)
	}
	return nil
}

func (nw *Network) load(dumpPath string) error {
	nwConfigFile, err := os.Open(dumpPath)
	defer nwConfigFile.Close()
	if err != nil {
		return err
	}
	nwJson := make([]byte, 2000)
	n, err := nwConfigFile.Read(nwJson)
	if err != nil {
		return err
	}
	err := json.Unmarshal(nwjson[:n], nw)
	if err != nil {
		logrus.Errorf("error load nw info", err)
		return err
	}
	return nil
}
