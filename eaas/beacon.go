package eaas

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"runtime"
	"sync"

	"github.com/icedream/go-stagelinq/internal/messages"
	"github.com/icedream/go-stagelinq/internal/socket"
	"golang.org/x/net/ipv4"
)

// Beacon listens on UDP port 11224 for EAAS clients and announces itself to them.
type Beacon struct {
	softwareVersion   string
	hostname          string
	packetConn4       *ipv4.PacketConn
	token             Token
	grpcHost          string
	grpcPort          uint16
	shutdownWaitGroup sync.WaitGroup
}

// Token returns our token that is being announced to the EAAS network. Use this
// token for further communication with services on other devices.
func (l *Beacon) Token() Token {
	return l.token
}

// Shutdown shuts down the listener.
func (l *Beacon) Shutdown() error {
	err := l.packetConn4.Close()

	// wait for Listen goroutine to finish
	l.shutdownWaitGroup.Wait()

	return err
}

// List will start a goroutine which waits for EAAS clients to announce back to
// them. It will automatically terminate once this listener is shut down.
func (l *Beacon) listen() {
	// make Close() wait for us
	l.shutdownWaitGroup.Add(1)

	go func() {
		defer l.shutdownWaitGroup.Done()

		b := make([]byte, 8)
		for {
			n, cm, addr, err := l.packetConn4.ReadFrom(b)
			if errors.Is(err, net.ErrClosed) {
				break
			}
			if err != nil {
				// TODO - log this somehow
				continue
			}
			if err = l.handleIncomingIPv4Packet(b[0:n], cm, addr); err != nil {
				// TODO - log this somehow
				continue
			}

		}
	}()
}

func (l *Beacon) getGRPCURL(ip net.IP) string {
	var host string
	switch {
	case len(l.grpcHost) > 0:
		host = l.grpcHost
	default:
		host = ip.String()
	}
	return fmt.Sprintf("grpc://%s:%d", host, l.grpcPort)
}

func (l *Beacon) replyIPv4(controlMessage *ipv4.ControlMessage, srcAddr net.Addr) error {
	var ip net.IP
	if controlMessage != nil && controlMessage.IfIndex > 0 {
		// Short path - figure out the an IP that we could be reachable from
		// based on the interface the broadcast came from
		netInterface, err := net.InterfaceByIndex(controlMessage.IfIndex)
		if err != nil {
			return err
		}
		interfaceAddresses, err := netInterface.Addrs()
		if err != nil {
			return err
		}
		ip = socket.GetIP(interfaceAddresses[0])
	} else {
		// Long path - figure out any IP-based path that works for the source to
		// reach us
		netInterfaces, err := net.Interfaces()
		if err != nil {
			return err
		}
		allInterfaceAddresses := []net.Addr{}
		for _, netInterface := range netInterfaces {
			interfaceAddresses, err := netInterface.Addrs()
			if err != nil {
				// TODO - Log this somehow? For now let's just move on...
				continue
			}
			allInterfaceAddresses = append(allInterfaceAddresses, interfaceAddresses...)
		}
		// Pick the first interface IP that could talk straight to the source
		srcIP := socket.GetIP(srcAddr)
		for _, interfaceAddr := range allInterfaceAddresses {
			interfaceIP := socket.GetIP(interfaceAddr)
			interfaceMask := socket.GetMaskFromAddress(interfaceAddr)
			intfnet := &net.IPNet{
				IP:   interfaceIP,
				Mask: interfaceMask,
			}
			if intfnet.Contains(srcIP) {
				ip = interfaceIP
				break
			}
		}
	}
	// TODO - optimization: cache the built message because it will be sent repeatedly?
	responseMessage := &eaasDiscoveryResponseMessage{
		TokenPrefixedMessage: messages.TokenPrefixedMessage{
			Token: messages.Token(l.token),
		},
		Hostname:        l.hostname,
		SoftwareVersion: l.softwareVersion,
		URL:             l.getGRPCURL(ip),
		Extra:           "_",
	}
	b := new(bytes.Buffer)
	if err := responseMessage.WriteMessageTo(b); err != nil {
		return err
	}
	var replyControlMessage *ipv4.ControlMessage
	if controlMessage != nil {
		replyControlMessage = &ipv4.ControlMessage{
			IfIndex: controlMessage.IfIndex,
		}
	}
	if _, err := l.packetConn4.WriteTo(b.Bytes(), replyControlMessage, srcAddr); err != nil {
		return err
	}
	return nil
}

func (l *Beacon) handleIncomingIPv4Packet(b []byte, cm *ipv4.ControlMessage, srcAddr net.Addr) error {
	// decode message
	r := bytes.NewReader(b)
	m := new(eaasDiscoveryRequestMessage)
	if err := m.ReadMessageFrom(r); err != nil {
		return err
	}

	return l.replyIPv4(cm, srcAddr)
}

// StartBeacon sets up an EAAS beacon.
func StartBeacon() (*Beacon, error) {
	return StartBeaconWithConfiguration(nil)
}

var zeroToken = Token{}

// StartBeaconWithConfiguration sets up a EAAS announcer with the given configuration.
func StartBeaconWithConfiguration(beaconConfig *BeaconConfiguration) (beacon *Beacon, err error) {
	// Use empty configuration if no configuration object was passed
	if beaconConfig == nil {
		beaconConfig = new(BeaconConfiguration)
	}

	// Use background context if none was configured
	ctx := beaconConfig.Context
	if ctx == nil {
		ctx = context.Background()
	}

	// Initialize token if none was configured
	token := beaconConfig.Token
	if bytes.Equal(beaconConfig.Token[:], zeroToken[:]) {
		if _, err = rand.Read(token[:]); err != nil {
			return
		}
	}

	// Use default EAAS gRPC port if none was set
	grpcPort := beaconConfig.GRPCPort
	if grpcPort == 0 {
		grpcPort = DefaultEAASGRPCPort
	}

	config := &net.ListenConfig{
		Control: socket.SetSocketControlForReusePort,
	}
	packetConn, err := config.ListenPacket(
		ctx,
		"udp4",
		makeEAASDiscoveryAddress(net.IPv4zero).String())
	if err != nil {
		return
	}

	ipv4PacketConn := ipv4.NewPacketConn(packetConn)
	// NOTE - this part only works on Linux and is unimplemented on Windows...
	//
	// This is however necessary so we return the correct IP for Engine software
	// to connect to.
	switch runtime.GOOS {
	case "windows":
		// No control message flags implemented, gotta sort out by source IP
	default:
		// Everything else can use the short path of just using interface IP
		if err := ipv4PacketConn.SetControlMessage(ipv4.FlagInterface, true); err != nil {
			err = errors.Join(
				fmt.Errorf(
					"failed to set control message flags: %w",
					err),
				packetConn.Close())
			return nil, err
		}
	}

	b := &Beacon{
		packetConn4:     ipv4PacketConn,
		hostname:        beaconConfig.Name,
		softwareVersion: beaconConfig.SoftwareVersion,
		token:           token,
		grpcHost:        beaconConfig.GRPCHost,
		grpcPort:        grpcPort,
	}
	go b.listen()

	return b, nil
}
