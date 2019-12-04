package network

type BridgeNetworkDriver struct {
}

func (br *BridgeNetworkDriver) Name() string {
	return "bridge"
}

func (br *BridgeNetworkDriver) Create(subnet string, name string) (*Network, error) {
	return &Network{}, nil
}

func (br *BridgeNetworkDriver) Delete(network Network) error {
	return nil
}

func (br *BridgeNetworkDriver) Connect(network *Network, endpoint *Endpoint) error {
	return nil
}

func (br *BridgeNetworkDriver) Disconnect(network Network, endpoint *Endpoint) error {
	return nil
}
