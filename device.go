package stagelinq

import "net"

// DeviceState represents a device's state in the network.
// Possible values are DevicePresent and DeviceLeaving.
type DeviceState byte

const (
	// DevicePresent indicates that a device is actively announcing itself to the network.
	DevicePresent DeviceState = iota

	// DeviceLeaving indicates that a device has announced that it is leaving the network. It will no longer send announcements after this.
	DeviceLeaving
)

// Device presents information about a discovered StagelinQ device on the network.
type Device struct {
	port  uint16
	token [16]byte

	IP              net.IP
	Name            string
	SoftwareName    string
	SoftwareVersion string
}

// IsEqual checks if this device has the same address and values as the other given device.
func (device *Device) IsEqual(anotherDevice *Device) bool {
	return device.token == anotherDevice.token &&
		device.Name == anotherDevice.Name &&
		device.SoftwareName == anotherDevice.SoftwareName &&
		device.SoftwareVersion == anotherDevice.SoftwareVersion
}

func newDeviceFromDiscovery(addr *net.UDPAddr, msg *DiscoveryMessage) *Device {
	return &Device{
		port:  msg.Port,
		token: msg.Token,

		IP:              addr.IP,
		Name:            msg.Source,
		SoftwareName:    msg.SoftwareName,
		SoftwareVersion: msg.SoftwareVersion,
	}
}
