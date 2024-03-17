package stagelinq

import (
	"net"
)

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
	token Token

	IP              net.IP
	Name            string
	SoftwareName    string
	SoftwareVersion string
}

// Dial starts a TCP connection with the device on the given port.
func (device *Device) Dial(port uint16) (conn net.Conn, err error) {
	ip := device.IP
	conn, err = net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   ip,
		Port: int(port),
	})

	return
}

// Connect starts a new main connection with the device.
// You need to pass the StagelinQ token announced for your own device.
// You also need to pass services you want to provide; if you don't have any, pass an empty array.
func (device *Device) Connect(token Token, offeredServices []*Service) (conn *MainConnection, err error) {
	tcpConn, err := device.Dial(device.port)
	if err != nil {
		return
	}
	conn, err = newMainConnection(tcpConn, token, device.token, offeredServices)
	return
}

// IsEqual checks if this device has the same address and values as the other given device.
func (device *Device) IsEqual(anotherDevice *Device) bool {
	return device.token == anotherDevice.token &&
		device.Name == anotherDevice.Name &&
		device.SoftwareName == anotherDevice.SoftwareName &&
		device.SoftwareVersion == anotherDevice.SoftwareVersion
}

func newDeviceFromDiscovery(addr *net.UDPAddr, msg *discoveryMessage) *Device {
	return &Device{
		port:  msg.Port,
		token: Token(msg.Token),

		IP:              addr.IP,
		Name:            msg.Source,
		SoftwareName:    msg.SoftwareName,
		SoftwareVersion: msg.SoftwareVersion,
	}
}
