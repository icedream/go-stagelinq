package eaas

import (
	"bytes"
	"context"
	"errors"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/icedream/go-stagelinq/internal/socket"
)

// Announcer listens on UDP port 11224 for EAAS clients and announces itself to them.
type Announcer struct {
	softwareVersion   string
	hostname          string
	packetConn        net.PacketConn
	token             Token
	port              uint16
	shutdownCond      *sync.Cond
	shutdownWaitGroup sync.WaitGroup
}

// Token returns our token that is being announced to the EAAS network. Use this
// token for further communication with services on other devices.
func (l *Announcer) Token() Token {
	return l.token
}

// Close shuts down the listener.
func (l *Announcer) Close() error {
	// notify goroutines we are going to shut down and wait for them to finish
	l.shutdownCond.Broadcast()
	l.shutdownWaitGroup.Wait()

	return l.packetConn.Close()
}

// Announce announces this EAAS listener to the network.
//
// This function should be called before actually listening in for devices to
// allow them to pick up our token for communication immediately.
func (l *Announcer) Announce() error {
	return l.announce()
}

// AnnounceEvery will start a goroutine which calls the Announce function at given interval.
// It will automatically terminate once this listener is shut down.
// A recommended value for the interval is 1 second.
func (l *Announcer) AnnounceEvery(interval time.Duration) {
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

func (l *Announcer) announce() (err error) {
	// TODO - optimization: cache the built message because it will be sent repeatedly?
	m := &eaasDiscoveryRequestMessage{}
	b := new(bytes.Buffer)
	err = m.WriteMessageTo(b)
	if err != nil {
		return
	}
	finalBytes := b.Bytes()
	ips, err := socket.GetAllBroadcastIPs()
	if err != nil {
		return
	}
	for _, ip := range ips {
		addr := makeEAASDiscoveryBroadcastAddress(ip)
		packetConn, err := net.DialUDP("udp", nil, addr)
		if err == nil {
			_, _ = packetConn.Write(finalBytes)
			packetConn.Close()
		}
	}

	return
}

// Announce sets up a EAAS announcer.
func Announce() (announcer *Announcer, err error) {
	return AnnounceWithConfiguration(nil)
}

var zeroToken = Token{}

// AnnounceWithConfiguration sets up a EAAS announcer with the given configuration.
func AnnounceWithConfiguration(announcerConfig *AnnouncerConfiguration) (announcer *Announcer, err error) {
	// Use empty configuration if no configuration object was passed
	if announcerConfig == nil {
		announcerConfig = new(AnnouncerConfiguration)
	}

	// Initialize token if none was configured
	token := announcerConfig.Token
	if bytes.Equal(announcerConfig.Token[:], zeroToken[:]) {
		if _, err = rand.Read(token[:]); err != nil {
			return
		}
	}

	// Use background context if none was configured
	ctx := announcerConfig.Context
	if ctx == nil {
		ctx = context.Background()
	}

	// We are setting up a shared UDP address socket here to allow other applications to still listen for EAAS discovery messages
	config := &net.ListenConfig{
		Control: socket.SetSocketControlForReusePort,
	}
	packetConn, err := config.ListenPacket(ctx, eaasDiscoveryNetwork, eaasDiscoveryAddressString)
	if err != nil {
		return
	}

	return &Announcer{
		hostname:        announcerConfig.Name,
		packetConn:      packetConn,
		softwareVersion: announcerConfig.SoftwareVersion,
		token:           token,
		shutdownCond:    sync.NewCond(&sync.Mutex{}),
	}, nil
}
