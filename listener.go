package stagelinq

import (
	"bytes"
	"context"
	"errors"
	"math/rand"
	"net"
	"sync"
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

const (
	stagelinqDiscoveryNetwork       = "udp"
	stagelinqDiscoveryAddressString = ":51337"
)

var magicBytes = []byte("airD")

func makeStagelinqDiscoveryBroadcastAddress(ip net.IP) *net.UDPAddr {
	return &net.UDPAddr{
		IP:   ip,
		Port: 51337,
	}
}

func getAllBroadcastIPs() (retval []net.IP, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	ips := []net.IP{}
addrsLoop:
	for _, addr := range addrs {
		var ip net.IP
		var mask net.IPMask
		switch v := addr.(type) {
		case *net.IPAddr:
			ip = v.IP
			mask = v.IP.DefaultMask()
		case *net.IPNet:
			ip = v.IP
			mask = v.Mask
		}
		if ip == nil {
			continue
		}

		// prevent addresses from being added multiple times (for example zeroconf)
		bip := makeBroadcastIP(ip, mask)
		for _, alreadyAddedIP := range ips {
			if alreadyAddedIP.Equal(bip) {
				continue addrsLoop
			}
		}

		ips = append(ips, bip)
	}

	retval = ips
	return
}

// Listener listens on UDP port 51337 for StagelinQ devices and announces itself in the same way.
type Listener struct {
	softwareName      string
	softwareVersion   string
	name              string
	packetConn        net.PacketConn
	token             Token
	port              uint16
	shutdownCond      *sync.Cond
	shutdownWaitGroup sync.WaitGroup
}

// Token returns our token that is being announced to the StagelinQ network.
// Use this token for further communication with services on other devices.
func (l *Listener) Token() Token {
	return l.token
}

// Close shuts down the listener.
func (l *Listener) Close() error {
	// notify goroutines we are going to shut down and wait for them to finish
	l.shutdownCond.Broadcast()
	l.shutdownWaitGroup.Wait()

	return l.packetConn.Close()
}

// Announce announces this StagelinQ listener to the network.
// This function should be called before actually listening in for devices to allow them to pick up our token for communication immediately.
func (l *Listener) Announce() error {
	return l.announce(discovererHowdy)
}

// AnnounceEvery will start a goroutine which calls the Announce function at given interval.
// It will automatically terminate once this listener is shut down.
// A recommended value for the interval is 1 second.
func (l *Listener) AnnounceEvery(interval time.Duration) {
	shutdownC := make(chan interface{}, 1)

	// make Close() wait for us
	l.shutdownWaitGroup.Add(1)

	// listen for shutdown signal broadcast, forward it to our own channel
	go func() {
		l.shutdownCond.L.Lock()
		defer l.shutdownCond.L.Unlock()
		l.shutdownCond.Wait()
		shutdownC <- nil
	}()

	go func() {
		defer l.shutdownWaitGroup.Done()

		// timestamp for when to send next announcement
		ticker := time.NewTicker(interval)

		// do first announcement immediately
		l.Announce()
		defer l.Unannounce()

		for {
			select {
			case <-ticker.C: // next interval - announcement
				if err := l.Announce(); errors.Is(err, net.ErrClosed) {
					return
				}
				// NOTE - Considering AnnounceEvery is a fire-and-forget command we're ignoring other errors here for now. Not sure how to properly handle them otherwise atm.
			case <-shutdownC:
				return
			}
		}
	}()
}

// Unannounce announces this StagelinQ listener leaving from the network.
// Call this before closing the listener!
func (l *Listener) Unannounce() error {
	return l.announce(discovererExit)
}

func (l *Listener) announce(action discovererMessageAction) (err error) {
	// TODO - optimization: cache the built message because it will be sent repeatedly?
	m := &discoveryMessage{
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
	finalBytes := b.Bytes()
	ips, err := getAllBroadcastIPs()
	if err != nil {
		return
	}
	for _, ip := range ips {
		addr := makeStagelinqDiscoveryBroadcastAddress(ip)
		packetConn, err := net.DialUDP("udp", nil, addr)
		if err == nil {
			_, _ = packetConn.Write(finalBytes)
			packetConn.Close()
		}
	}

	return
}

// Discover listens for any StagelinQ devices announcing to the network.
// If no device is found within the given timeout or any non-StagelinQ message has been received, nil is returned for the device.
// If a device has been discovered before, the returned device object is not going to be the same as when the device was previously discovered.
// Use device.IsEqual for such comparison.
func (l *Listener) Discover(timeout time.Duration) (device *Device, deviceState DeviceState, err error) {
	b := make([]byte, 8*1024)

	if timeout != 0 {
		l.packetConn.SetReadDeadline(time.Now().Add(timeout))
	}

readLoop:
	for {
		var n int
		var src net.Addr
		n, src, err = l.packetConn.ReadFrom(b)
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				// ignore i/o timeout since we set the timeout ourself
				err = nil
			}
			return
		}

		// message smaller than expected magic bytes?
		if n < 4 {
			err = ErrTooShortDiscoveryMessageReceived
			return
		}

		// decode message
		r := bytes.NewReader(b)
		m := new(discoveryMessage)
		if err = m.readFrom(r); err != nil {
			return
		}

		device = newDeviceFromDiscovery(src.(*net.UDPAddr), m)

		// is this just ourself?
		if bytes.Equal(device.token[:], l.token[:]) &&
			device.Name == l.name &&
			device.SoftwareName == l.softwareName &&
			device.SoftwareVersion == l.softwareVersion {
			// ignore
			continue readLoop
		}

		switch m.Action {
		case discovererExit:
			deviceState = DeviceLeaving
		case discovererHowdy:
			deviceState = DevicePresent
		default:
			err = ErrInvalidDiscovererActionReceived
			return
		}

		return
	}
}

// Listen sets up a StagelinQ listener.
func Listen() (listener *Listener, err error) {
	return ListenWithConfiguration(nil)
}

var zeroToken = Token{}

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
		token:           token,
		shutdownCond:    sync.NewCond(&sync.Mutex{}),
	}

	return
}
