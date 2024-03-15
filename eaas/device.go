package eaas

// Device presents information about a discovered StagelinQ device on the network.
type Device struct {
	port  uint16
	token Token

	URL             string
	Hostname        string
	SoftwareVersion string
}

// TODO - func (device *Device) Dial() (conn net.Conn, err error)
// TODO - func (device *Device) Connect(token Token, offeredServices []*Service) (conn *MainConnection, err error) {

// IsEqual checks if this device has the same address and values as the other given device.
func (device *Device) IsEqual(anotherDevice *Device) bool {
	return device.token == anotherDevice.token &&
		device.Hostname == anotherDevice.Hostname &&
		device.URL == anotherDevice.URL &&
		device.SoftwareVersion == anotherDevice.SoftwareVersion
}

func newDeviceFromDiscovery(msg *eaasDiscoveryResponseMessage) *Device {
	return &Device{
		token: Token(msg.Token),

		URL:             msg.URL,
		Hostname:        msg.Hostname,
		SoftwareVersion: msg.SoftwareVersion,
	}
}
