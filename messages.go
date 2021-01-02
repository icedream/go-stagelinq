package stagelinq

import (
	"encoding/binary"
	"io"
)

// Token contains the identifying token for a device in the StagelinQ network.
type Token [16]byte

type message interface {
	readFrom(io.Reader) error
	writeTo(io.Writer) error
	id() int32
}

type tokenPrefixedMessage struct {
	Token Token
}

func (m *tokenPrefixedMessage) readFrom(r io.Reader) (err error) {
	_, err = r.Read(m.Token[:])
	return
}

func (m *tokenPrefixedMessage) writeTo(w io.Writer) (err error) {
	_, err = w.Write(m.Token[:])
	return
}

type ServiceAnnouncementMessage struct {
	tokenPrefixedMessage
	Service string
	Port    uint16
}

func (m *ServiceAnnouncementMessage) id() int32 {
	return 0x00000000
}

func (m *ServiceAnnouncementMessage) readFrom(r io.Reader) (err error) {
	if err = m.tokenPrefixedMessage.readFrom(r); err != nil {
		return
	}
	if err = readNetworkString(r, &m.Service); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &m.Port); err != nil {
		return
	}
	return
}

func (m *ServiceAnnouncementMessage) writeTo(w io.Writer) (err error) {
	if err = m.tokenPrefixedMessage.writeTo(w); err != nil {
		return
	}
	if err = writeNetworkString(w, m.Service); err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, m.Port)
	return
}

type PingMessage struct {
	tokenPrefixedMessage
	Token2 Token
	Data   []byte
}

func (m *PingMessage) id() int32 {
	return 0x00000001
}

func (m *PingMessage) readFrom(r io.Reader) (err error) {
	if err = m.tokenPrefixedMessage.readFrom(r); err != nil {
		return
	}
	if _, err = r.Read(m.Token2[:]); err != nil {
		return
	}
	buf := make([]byte, 8)
	n, err := r.Read(buf)
	if err != nil {
		return
	}
	m.Data = buf[0:n]
	return
}

func (m *PingMessage) writeTo(w io.Writer) (err error) {
	if err = m.tokenPrefixedMessage.writeTo(w); err != nil {
		return
	}
	if _, err = w.Write(m.Token2[:]); err != nil {
		return
	}
	_, err = w.Write(m.Data)
	return
}

type EmptyMessage struct {
	tokenPrefixedMessage
}

func (m *EmptyMessage) id() int32 {
	return 0x00000002
}

// func (m *EmptyMessage) readFrom(r io.Reader) (err error) {
// 	if err = m.tokenPrefixedMessage.readFrom(r); err != nil {
// 		return
// 	}
// 	return
// }

// func (m *EmptyMessage) writeTo(w io.Writer) (err error) {
// 	if err = m.tokenPrefixedMessage.writeTo(w); err != nil {
// 		return
// 	}
// 	return
// }

// DiscovererMessageAction is the action taken by a device as part of StagelinQ device discovery.
// Possible values are DiscovererHowdy or DiscovererExit.
type DiscovererMessageAction string

const (
	// DiscovererHowdy is the value set on the Action field of a DiscoveryMessage when a StagelinQ-compatible device announces itself in the network.
	DiscovererHowdy DiscovererMessageAction = "DISCOVERER_HOWDY_"

	// DiscovererExit is the value set on the Action field of a DiscoveryMessage when a StagelinQ-compatible device leaves the network.
	DiscovererExit DiscovererMessageAction = "DISCOVERER_EXIT_"
)

// DiscoveryMessage contains the data carried in the message payload for device trying to handshake the StagelinQ protocol to any other device in the network.
type DiscoveryMessage struct {
	tokenPrefixedMessage
	Source          string
	Action          DiscovererMessageAction
	SoftwareName    string
	SoftwareVersion string
	Port            uint16
}

func (m *DiscoveryMessage) id() int32 {
	return 0x61697244
}

func (m *DiscoveryMessage) readFrom(r io.Reader) (err error) {
	if err = m.tokenPrefixedMessage.readFrom(r); err != nil {
		return
	}
	if err = readNetworkString(r, &m.Source); err != nil {
		return
	}
	actionString := ""
	if err = readNetworkString(r, &actionString); err != nil {
		return
	}
	m.Action = DiscovererMessageAction(actionString)
	if err = readNetworkString(r, &m.SoftwareName); err != nil {
		return
	}
	if err = readNetworkString(r, &m.SoftwareVersion); err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &m.Port)
	return
}

func (m *DiscoveryMessage) writeTo(w io.Writer) (err error) {
	if err = m.tokenPrefixedMessage.writeTo(w); err != nil {
		return
	}
	if err = writeNetworkString(w, m.Source); err != nil {
		return
	}
	if err = writeNetworkString(w, string(m.Action)); err != nil {
		return
	}
	if err = writeNetworkString(w, m.SoftwareName); err != nil {
		return
	}
	if err = writeNetworkString(w, m.SoftwareVersion); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, m.Port); err != nil {
		return
	}
	return
}

var tcpMessageMap = map[int32]func() message{
	0x00000000: func() message { return new(ServiceAnnouncementMessage) },
	0x00000001: func() message { return new(PingMessage) },
	0x00000002: func() message { return new(EmptyMessage) },
}

var udpMessageMap = map[int32]func() message{
	0x61697244 /* "airD" */ : func() message { return new(DiscoveryMessage) },
}
