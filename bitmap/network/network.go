package network

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"

	"github.com/Sirupsen/logrus"

	"github.com/vishvananda/netlink"
)

const (
	defaultNetworkPath = "/tmp/mynet/"
)

var (
	networks = map[string]*Network{}
	drivers  = map[string]NetworkDriver{}
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
	// fullpath := dumpPath + nw.Name
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
	err = json.Unmarshal(nwJson[:n], nw)
	if err != nil {
		logrus.Errorf("error load nw info", err)
		return err
	}

	fmt.Println(nw)
	fmt.Println(string(nwJson[:n]))
	return nil
}

func Init() error {
	fmt.Println("Iinit will called ?")
	var bridgeDriver = BridgeNetworkDriver{}
	drivers[bridgeDriver.Name()] = &bridgeDriver

	//check if network exist
	if _, err := os.Stat(defaultNetworkPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(defaultNetworkPath, 0755)
		} else {
			return err
		}
	}

	//check all network config files
	filepath.Walk(defaultNetworkPath, func(nwPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		_, nwName := path.Split(nwPath)
		nw := &Network{
			Name: nwName,
		}
		if err := nw.Load(nwPath); err != nil {
			logrus.Errorf("error load network: %s", err)
		}
		networks[nwName] = nw
		return nil
	})
	return nil
}

func ListNetwork() {
	fmt.Println("######################")
	for _, nw := range networks {
		fmt.Println(nw)
	}
	fmt.Println("######################")
}

func DeleteNetwork(networkName string) error {
	nw, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("no such network:%s", networkName)
	}
	return nw.remove(defaultNetworkPath)
}

func (nw *Network) remove(dumpPath string) error {
	if _, err := os.Stat(path.Join(dumpPath, nw.Name)); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	} else {
		delete(networks, nw.Name)
		return os.Remove(path.Join(dumpPath, nw.Name))
	}
}
