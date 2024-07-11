package eaas

// Device presents information about a discovered EAAS device on the network.
type Device struct {
	port  uint16
	token Token

	URL             string
	Hostname        string
	SoftwareVersion string
}

// TODO - func (device *Device) Dial() - gRPC client

// IsEqual checks if this device has the same identifying token as the other
// given device.
func (device *Device) IsEqual(anotherDevice *Device) bool {
	return device.token == anotherDevice.token
}

func newDeviceFromDiscovery(msg *eaasDiscoveryResponseMessage) *Device {
	return &Device{
		token: Token(msg.Token),

		URL:             msg.URL,
		Hostname:        msg.Hostname,
		SoftwareVersion: msg.SoftwareVersion,
	}
}
