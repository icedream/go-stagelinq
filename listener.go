package stagelinq

import (
	"bytes"
	"context"
	"errors"
	"math/rand"
	"net"
	"syscall"
	"time"
)

// ErrTooShortDiscoveryMessageReceived is returned by Listener.Discover if a
// too short message has been received on the StagelinQ discovery port.
// This would indicate another application using UDP port 51337 on the network
// for broadcasts.
var ErrTooShortDiscoveryMessageReceived = errors.New("too short discovery message received")

// ErrInvalidMessageReceived is returned by Listener.Discover if a message has
// been received but it is not a StagelinQ message.
// This would indicate another application using UDP port 51337 on the network
// for broadcasts.
var ErrInvalidMessageReceived = errors.New("invalid message received")

// ErrInvalidDiscovererActionReceived is returned by Listener.Discover if a
// valid StagelinQ discovery message has been received by another device but it
// is reporting neither that it is leaving nor joining the network.
// This would indicate another application trying to speak via the StagelinQ
// protocol but it is sending invalid data.
// You can check the returned device object for the source address of the bad
// message.
var ErrInvalidDiscovererActionReceived = errors.New("invalid discoverer action received")

const stagelinqDiscoveryNetwork = "udp"
const stagelinqDiscoveryAddressString = "0.0.0.0:51337"

var stagelinqDiscoveryBroadcastAddress = &net.UDPAddr{
	IP:   net.IPv4(255, 255, 255, 255),
	Port: 51337,
}

var magicBytes = []byte("airD")

func setSocketControlForReusePort(network, address string, c syscall.RawConn) error {
	return c.Control(func(fd uintptr) {
		syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
	})
}

// Listener listens on UDP port 51337 for StagelinQ devices and announces itself in the same way.
type Listener struct {
	softwareName    string
	softwareVersion string
	name            string
	packetConn      net.PacketConn
	token           [16]byte
	timeout         time.Duration
	port            uint16
}

// Close shuts down the listener.
func (l *Listener) Close() error {
	return l.packetConn.Close()
}

// Announce announces this StagelinQ listener to the network.
// This function should be called before actually listening in for devices to allow them to pick up our token for communication immediately.
func (l *Listener) Announce() error {
	return l.announce(DiscovererHowdy)
}

func (l *Listener) announce(action DiscovererMessageAction) (err error) {
	// TODO - optimization: cache the built message because it will be sent repeatedly?
	m := &DiscoveryMessage{
		Source:          l.name,
		SoftwareName:    l.softwareName,
		SoftwareVersion: l.softwareVersion,
		tokenPrefixedMessage: tokenPrefixedMessage{
			Token: l.token,
		},
		Action: action,
		Port:   l.port,
	}
	b := new(bytes.Buffer)
	err = m.writeTo(b)
	if err != nil {
		return
	}

	_, err = l.packetConn.WriteTo(b.Bytes(), stagelinqDiscoveryBroadcastAddress)
	return
}

// Discover listens for any StagelinQ devices announcing to the network.
// If no device is found within the configured timeout or any non-StagelinQ message has been received, nil is returned for the device.
// If a device has been discovered before, the returned device object is not going to be the same as when the device was previously discovered.
// Use device.IsEqual for such comparison.
func (l *Listener) Discover() (device *Device, deviceState DeviceState, err error) {
	b := make([]byte, 8*1024)

	if l.timeout != 0 {
		l.packetConn.SetReadDeadline(time.Now().Add(l.timeout))
	}

	n, src, err := l.packetConn.ReadFrom(b)
	if err != nil {
		return
	}

	// message smaller than expected magic bytes?
	if n < 4 {
		err = ErrTooShortDiscoveryMessageReceived
		return
	}

	// do first bytes match expected magic bytes?
	if !bytes.Equal(b[0:4], magicBytes) {
		err = ErrInvalidMessageReceived
		return
	}

	// decode message
	r := bytes.NewReader(b[4:n])
	m := new(DiscoveryMessage)
	if err = m.readFrom(r); err != nil {
		return
	}

	device = newDeviceFromDiscovery(src.(*net.UDPAddr), m)

	switch m.Action {
	case DiscovererExit:
		deviceState = DeviceLeaving
	case DiscovererHowdy:
		deviceState = DevicePresent
	default:
		err = ErrInvalidDiscovererActionReceived
		return
	}

	return
}

// Listen sets up a StagelinQ listener.
func Listen() (listener *Listener, err error) {
	return ListenWithConfiguration(nil)
}

var zeroToken = [16]byte{}

// ListenWithConfiguration sets up a StagelinQ listener with the given configuration.
func ListenWithConfiguration(listenerConfig *ListenerConfiguration) (listener *Listener, err error) {
	// Use empty configuration if no configuration object was passed
	if listenerConfig == nil {
		listenerConfig = new(ListenerConfiguration)
	}

	// Initialize token if none was configured
	token := listenerConfig.Token
	if bytes.Equal(listenerConfig.Token[:], zeroToken[:]) {
		if _, err = rand.Read(token[:]); err != nil {
			return
		}
	}

	// Use background context if none was configured
	ctx := listenerConfig.Context
	if ctx == nil {
		ctx = context.Background()
	}

	// We are setting up a shared UDP address socket here to allow other applications to still listen for StagelinQ discovery messages
	config := &net.ListenConfig{
		Control: setSocketControlForReusePort,
	}
	packetConn, err := config.ListenPacket(ctx, stagelinqDiscoveryNetwork, stagelinqDiscoveryAddressString)
	if err != nil {
		return
	}

	listener = &Listener{
		name:            listenerConfig.Name,
		packetConn:      packetConn,
		softwareName:    listenerConfig.SoftwareName,
		softwareVersion: listenerConfig.SoftwareVersion,
		timeout:         listenerConfig.DiscoveryTimeout,
		token:           token,
	}

	return
}
