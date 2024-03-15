package eaas

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/icedream/go-stagelinq/internal/messages"
	"github.com/icedream/go-stagelinq/internal/socket"
)

type Token messages.Token

// ErrTooShortDiscoveryMessageReceived is returned by Listener.Discover if a
// too short message has been received on the EAAS discovery port.
// This would indicate another application using UDP port 11224 on the network
// for broadcasts.
var ErrTooShortDiscoveryMessageReceived = errors.New("too short discovery message received")

// ErrInvalidMessageReceived is returned by Listener.Discover if a message has
// been received but it is not a EAAS message.
// This would indicate another application using UDP port 11224 on the network
// for broadcasts.
var ErrInvalidMessageReceived = errors.New("invalid message received")

const (
	eaasDiscoveryNetwork       = "udp"
	eaasDiscoveryAddressString = ":11224"
)

func makeEAASDiscoveryBroadcastAddress(ip net.IP) *net.UDPAddr {
	return &net.UDPAddr{
		IP:   ip,
		Port: 11224,
	}
}

// Listener listens on UDP port 11224 for EAAS devices and announces itself in the same way.
type Listener struct {
	packetConn        net.PacketConn
	token             Token
	shutdownCond      *sync.Cond
	shutdownWaitGroup sync.WaitGroup
}

// Close shuts down the listener.
func (l *Listener) Close() error {
	// notify goroutines we are going to shut down and wait for them to finish
	l.shutdownCond.Broadcast()
	l.shutdownWaitGroup.Wait()

	return l.packetConn.Close()
}

// SendBeacon requests EAAS devices to send a response back on the network.
func (l *Listener) SendBeacon() error {
	return l.sendBeacon()
}

// SendBeaconEvery will start a goroutine which calls the [Listener.SendBeacon]
// function at given interval. It will automatically terminate once this
// listener is shut down. A recommended value for the interval is 5 seconds.
func (l *Listener) SendBeaconEvery(interval time.Duration) {
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
		l.SendBeacon()

		for {
			select {
			case <-ticker.C: // next interval - announcement
				if err := l.SendBeacon(); errors.Is(err, net.ErrClosed) {
					return
				}
				// NOTE - Considering AnnounceEvery is a fire-and-forget command we're ignoring other errors here for now. Not sure how to properly handle them otherwise atm.
			case <-shutdownC:
				return
			}
		}
	}()
}

func (l *Listener) sendBeacon() (err error) {
	// TODO - optimization: cache the built message because it will be sent repeatedly?
	m := &eaasDiscoveryRequestMessage{}
	b := new(bytes.Buffer)
	err = m.WriteMessageTo(b)
	if err != nil {
		log.Println(err)
		return
	}
	finalBytes := b.Bytes()
	ips, err := socket.GetAllBroadcastIPs()
	packetConn := l.packetConn
	for _, bcastIP := range ips {
		bcastAddr := &net.UDPAddr{
			IP:   bcastIP,
			Port: 11224,
		}
		_, err = packetConn.WriteTo(finalBytes, bcastAddr)
		if err != nil {
			log.Println(err)
		}
	}

	return
}

// Discover listens for any EAAS devices announcing to the network.
// If no device is found within the given timeout or any non-EAAS message has been received, nil is returned for the device.
// If a device has been discovered before, the returned device object is not going to be the same as when the device was previously discovered.
// Use device.IsEqual for such comparison.
func (l *Listener) Discover(timeout time.Duration) (device *Device, err error) {
	b := make([]byte, 8*1024)

	if timeout != 0 {
		l.packetConn.SetReadDeadline(time.Now().Add(timeout))
	}

	for {
		var n int
		n, _, err = l.packetConn.ReadFrom(b)
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
		m := new(eaasDiscoveryResponseMessage)
		if err = m.ReadMessageFrom(r); err != nil {
			return
		}

		device = newDeviceFromDiscovery(m)
		return device, nil
	}
}

// Listen sets up a EAAS listener.
func Listen() (listener *Listener, err error) {
	return ListenWithConfiguration(nil)
}

// ListenWithConfiguration sets up a EAAS listener with the given configuration.
func ListenWithConfiguration(listenerConfig *ListenerConfiguration) (listener *Listener, err error) {
	// Use empty configuration if no configuration object was passed
	if listenerConfig == nil {
		listenerConfig = new(ListenerConfiguration)
	}

	// Use background context if none was configured
	ctx := listenerConfig.Context
	if ctx == nil {
		ctx = context.Background()
	}

	// select random source port
	config := &net.ListenConfig{
		Control: socket.SetSocketControlForReusePort,
	}
	packetConn, err := config.ListenPacket(ctx, eaasDiscoveryNetwork, ":0")
	if err != nil {
		return
	}
	log.Println("Listener listening on", packetConn.LocalAddr())

	listener = &Listener{
		packetConn:   packetConn,
		shutdownCond: sync.NewCond(&sync.Mutex{}),
	}

	return
}
