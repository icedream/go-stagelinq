package stagelinq

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/icedream/go-stagelinq/internal/messages"
)

type Token messages.Token

type serviceAnnouncementMessage struct {
	messages.TokenPrefixedMessage
	Service string
	Port    uint16
}

func (m *serviceAnnouncementMessage) CheckMatch(r *bufio.Reader) (ok bool, err error) {
	id, err := messages.PeekMessageID(r)
	if err != nil {
		return
	}
	ok = id == 0x00000000
	return
}

func (m *serviceAnnouncementMessage) ReadMessageFrom(r io.Reader) (err error) {
	messageID, err := messages.ReadMessageID(r)
	if err != nil {
		return
	} else if messageID != 0x00000000 {
		err = ErrInvalidMessageReceived
		return
	}
	if err = m.TokenPrefixedMessage.ReadMessageFrom(r); err != nil {
		return
	}
	if err = messages.ReadUTF16NetworkString(r, &m.Service); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &m.Port); err != nil {
		return
	}
	return
}

func (m *serviceAnnouncementMessage) WriteMessageTo(w io.Writer) (err error) {
	if err = messages.WriteMessageID(w, 0x00000000); err != nil {
		return
	}
	if err = m.TokenPrefixedMessage.WriteMessageTo(w); err != nil {
		return
	}
	if err = messages.WriteUTF16NetworkString(w, m.Service); err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, m.Port)
	return
}

type referenceMessage struct {
	messages.TokenPrefixedMessage
	Token2    messages.Token
	Reference int64
}

func (m *referenceMessage) CheckMatch(r *bufio.Reader) (ok bool, err error) {
	id, err := messages.PeekMessageID(r)
	if err != nil {
		return
	}
	ok = id == 0x00000001
	return
}

func (m *referenceMessage) ReadMessageFrom(r io.Reader) (err error) {
	messageID, err := messages.ReadMessageID(r)
	if err != nil {
		return
	} else if messageID != 0x00000001 {
		err = ErrInvalidMessageReceived
		return
	}
	if err = m.TokenPrefixedMessage.ReadMessageFrom(r); err != nil {
		return
	}
	if _, err = r.Read(m.Token2[:]); err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &m.Reference)
	return
}

func (m *referenceMessage) WriteMessageTo(w io.Writer) (err error) {
	if err = messages.WriteMessageID(w, 0x00000001); err != nil {
		return
	}
	if err = m.TokenPrefixedMessage.WriteMessageTo(w); err != nil {
		return
	}
	if _, err = w.Write(m.Token2[:]); err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, m.Reference)
	return
}

type servicesRequestMessage struct {
	messages.TokenPrefixedMessage
}

func (m *servicesRequestMessage) CheckMatch(r *bufio.Reader) (ok bool, err error) {
	id, err := messages.PeekMessageID(r)
	if err != nil {
		return
	}
	ok = id == 0x00000002
	return
}

func (m *servicesRequestMessage) ReadMessageFrom(r io.Reader) (err error) {
	messageID, err := messages.ReadMessageID(r)
	if err != nil {
		return
	} else if messageID != 0x00000002 {
		err = ErrInvalidMessageReceived
		return
	}

	err = m.TokenPrefixedMessage.ReadMessageFrom(r)
	return
}

func (m *servicesRequestMessage) WriteMessageTo(w io.Writer) (err error) {
	if err = messages.WriteMessageID(w, 0x00000002); err != nil {
		return
	}

	err = m.TokenPrefixedMessage.WriteMessageTo(w)
	return
}

// TODO - StateSubscribeMessage.Interval: Check what Interval actually is, it seems to either be 00 00 00 00 (Resolume Arena) or 00 00 00 0a (SoundSwitch)

var smaaMagicBytes = []byte{0x73, 0x6d, 0x61, 0x61}

func checkSmaa(r *bufio.Reader, id int32) (ok bool, err error) {
	// peek length bytes and smaa magic bytes
	b, err := r.Peek(4 + 4 + 4)
	if err != nil {
		return
	}

	// check smaa magic bytes
	if ok = bytes.Equal(b[4:8], smaaMagicBytes); !ok {
		return
	}

	// check id
	if int32(binary.BigEndian.Uint32(b[8:12])) != id {
		ok = false
	}

	return
}

type stateSubscribeMessage struct {
	// Length uint32
	// Unknown []byte = {0x73,0x6d,0x61,0x61}
	// Unknown2 []byte = {0x00,0x00,0x07,0xd2}
	Name     string
	Interval uint32
}

func (m *stateSubscribeMessage) CheckMatch(r *bufio.Reader) (ok bool, err error) {
	return checkSmaa(r, 0x000007d2)
}

func (m *stateSubscribeMessage) ReadMessageFrom(r io.Reader) (err error) {
	var expectedLength uint32
	if err = binary.Read(r, binary.BigEndian, &expectedLength); err != nil {
		return
	}

	// read smaa magic bytes
	magicBytes := make([]byte, 4)
	if _, err = r.Read(magicBytes); err != nil {
		return
	}
	if !bytes.Equal(magicBytes, smaaMagicBytes) {
		err = errors.New("invalid smaa magic bytes")
		return
	}

	// TODO - figure this out
	if _, err = r.Read(magicBytes); err != nil {
		return
	}
	if !bytes.Equal(magicBytes, []byte{0x00, 0x00, 0x07, 0xd2}) {
		err = errors.New("invalid post-smaa magic bytes")
		return
	}

	// read value name
	if err = messages.ReadUTF16NetworkString(r, &m.Name); err != nil {
		return
	}

	// TODO - figure this out
	err = binary.Read(r, binary.BigEndian, &m.Interval)
	return
}

func (m *stateSubscribeMessage) WriteMessageTo(w io.Writer) (err error) {
	// write smaa magic bytes
	buf := new(bytes.Buffer)
	if _, err = buf.Write(smaaMagicBytes); err != nil {
		return
	}

	// TODO - figure this out
	if _, err = buf.Write([]byte{0x00, 0x00, 0x07, 0xd2}); err != nil {
		return
	}

	// write value name
	if err = messages.WriteUTF16NetworkString(buf, m.Name); err != nil {
		return
	}

	// TODO - figure this out
	if err = binary.Write(buf, binary.BigEndian, m.Interval); err != nil {
		return
	}

	// send message length over wire
	if err = binary.Write(w, binary.BigEndian, uint32(buf.Len())); err != nil {
		return
	}

	// send actual message over wire
	_, err = w.Write(buf.Bytes())
	return
}

type stateEmitResponseMessage struct {
	// Length uint32
	// Unknown []byte = {0x73,0x6d,0x61,0x61}
	// Unknown2 []byte = {0x00,0x00,0x07,0xd1}
	Name     string
	Interval uint32
}

func (m *stateEmitResponseMessage) CheckMatch(r *bufio.Reader) (ok bool, err error) {
	return checkSmaa(r, 0x000007d1)
}

func (m *stateEmitResponseMessage) ReadMessageFrom(r io.Reader) (err error) {
	var expectedLength uint32
	if err = binary.Read(r, binary.BigEndian, &expectedLength); err != nil {
		return
	}

	// read smaa magic bytes
	magicBytes := make([]byte, 4)
	if _, err = r.Read(magicBytes); err != nil {
		return
	}
	if !bytes.Equal(magicBytes, smaaMagicBytes) {
		err = errors.New("invalid smaa magic bytes")
		return
	}

	// TODO - figure this out
	if _, err = r.Read(magicBytes); err != nil {
		return
	}
	if !bytes.Equal(magicBytes, []byte{0x00, 0x00, 0x07, 0xd1}) {
		err = errors.New("invalid post-smaa magic bytes")
		return
	}

	// read value name
	if err = messages.ReadUTF16NetworkString(r, &m.Name); err != nil {
		return
	}

	// TODO - figure this out
	err = binary.Read(r, binary.BigEndian, &m.Interval)
	return
}

func (m *stateEmitResponseMessage) WriteMessageTo(w io.Writer) (err error) {
	// write smaa magic bytes
	buf := new(bytes.Buffer)
	if _, err = buf.Write(smaaMagicBytes); err != nil {
		return
	}

	// TODO - figure this out
	if _, err = buf.Write([]byte{0x00, 0x00, 0x07, 0xd1}); err != nil {
		return
	}

	// write value name
	if err = messages.WriteUTF16NetworkString(buf, m.Name); err != nil {
		return
	}

	// TODO - figure this out
	if err = binary.Write(buf, binary.BigEndian, m.Interval); err != nil {
		return
	}

	// send message length over wire
	if err = binary.Write(w, binary.BigEndian, uint32(buf.Len())); err != nil {
		return
	}

	// send actual message over wire
	_, err = w.Write(buf.Bytes())
	return
}

type stateEmitMessage struct {
	// Length uint32
	// Unknown []byte = {0x73,0x6d,0x61,0x61}
	// Unknown2 []byte = {0x00,0x00,0x00,0x00}
	Name string
	JSON string
}

func (m *stateEmitMessage) CheckMatch(r *bufio.Reader) (ok bool, err error) {
	return checkSmaa(r, 0x00000000)
}

func (m *stateEmitMessage) ReadMessageFrom(r io.Reader) (err error) {
	// read expected message length
	var expectedLength uint32
	if err = binary.Read(r, binary.BigEndian, &expectedLength); err != nil {
		return
	}

	// set up buffer to write message into
	msgBytes := make([]byte, int(expectedLength))
	msgBytesOffset := 0
	for msgBytesOffset < int(expectedLength) {
		var n int
		if n, err = r.Read(msgBytes[msgBytesOffset:]); err != nil {
			return
		}
		msgBytesOffset += n
	}
	msgReader := bytes.NewReader(msgBytes)

	// read smaa magic bytes
	magicBytes := make([]byte, 4)
	if _, err = msgReader.Read(magicBytes); err != nil {
		return
	}
	if !bytes.Equal(magicBytes, smaaMagicBytes) {
		err = errors.New("invalid smaa magic bytes")
		return
	}

	// TODO - figure this out
	if _, err = msgReader.Read(magicBytes); err != nil {
		return
	}
	if !bytes.Equal(magicBytes, []byte{0x00, 0x00, 0x00, 0x00}) {
		err = errors.New("invalid post-smaa magic bytes")
		return
	}

	// read value name
	if err = messages.ReadUTF16NetworkString(msgReader, &m.Name); err != nil {
		return
	}

	// read value JSON
	if err = messages.ReadUTF16NetworkString(msgReader, &m.JSON); err != nil {
		return
	}

	return
}

func (m *stateEmitMessage) WriteMessageTo(w io.Writer) (err error) {
	buf := new(bytes.Buffer)

	// write smaa magic bytes to message buffer
	if _, err = buf.Write(smaaMagicBytes); err != nil {
		return
	}

	// TODO - figure this out
	if _, err = buf.Write([]byte{0x00, 0x00, 0x00, 0x00}); err != nil {
		return
	}

	// write value name to message buffer
	if err = messages.WriteUTF16NetworkString(buf, m.Name); err != nil {
		return
	}

	// write value JSON to message buffer
	if err = messages.WriteUTF16NetworkString(buf, m.JSON); err != nil {
		return
	}

	// send message length over wire
	if err = binary.Write(w, binary.BigEndian, uint32(buf.Len())); err != nil {
		return
	}

	// send actual message over wire
	_, err = w.Write(buf.Bytes())
	return
}

// BeatInfo

var beatInfoStartStreamMagicBytes = []byte{0x0, 0x0, 0x0, 0x0}

type beatInfoStartStreamMessage struct{}

func (m *beatInfoStartStreamMessage) CheckMatch(r *bufio.Reader) (ok bool, err error) {
	// peek length bytes and magic bytes
	b, err := r.Peek(4 + 4)
	if err != nil {
		return
	}
	// check magic bytes
	if ok = bytes.Equal(b[4:8], beatInfoStartStreamMagicBytes); !ok {
		return
	}
	return
}

func (m *beatInfoStartStreamMessage) ReadMessageFrom(r io.Reader) (err error) {
	// read expected message length
	var expectedLength uint32
	if err = binary.Read(r, binary.BigEndian, &expectedLength); err != nil {
		return
	}

	// set up buffer to write message into
	msgBytes := make([]byte, int(expectedLength))
	msgBytesOffset := 0
	for msgBytesOffset < int(expectedLength) {
		var n int
		if n, err = r.Read(msgBytes[msgBytesOffset:]); err != nil {
			return
		}
		msgBytesOffset += n
	}
	msgReader := bytes.NewReader(msgBytes)

	// read beatInfoStartStream magic bytes
	magicBytes := make([]byte, 4)
	if _, err = msgReader.Read(magicBytes); err != nil {
		return
	}
	if !bytes.Equal(magicBytes, beatInfoStartStreamMagicBytes) {
		return errors.New("invalid magic bytes")
	}

	return
}

func (m *beatInfoStartStreamMessage) WriteMessageTo(w io.Writer) (err error) {
	buf := new(bytes.Buffer)

	payload_len := len(beatInfoStartStreamMagicBytes)
	if err = binary.Write(buf, binary.BigEndian, uint32(payload_len)); err != nil {
		return
	}

	// write BeatInfo "start stream" magic bytes
	if _, err = buf.Write(beatInfoStartStreamMagicBytes); err != nil {
		return
	}

	// send actual message over wire
	_, err = w.Write(buf.Bytes())
	return
}

var beatInfoStopStreamMagicBytes = []byte{0x0, 0x0, 0x0, 0x1}

type beatInfoStopStreamMessage struct{}

func (m *beatInfoStopStreamMessage) CheckMatch(r *bufio.Reader) (ok bool, err error) {
	// peek length bytes and magic bytes
	b, err := r.Peek(4 + 4)
	if err != nil {
		return
	}
	// check magic bytes
	if ok = bytes.Equal(b[4:8], beatInfoStopStreamMagicBytes); !ok {
		return
	}
	return
}

func (m *beatInfoStopStreamMessage) ReadMessageFrom(r io.Reader) (err error) {
	// read expected message length
	var expectedLength uint32
	if err = binary.Read(r, binary.BigEndian, &expectedLength); err != nil {
		return
	}

	// set up buffer to write message into
	msgBytes := make([]byte, int(expectedLength))
	msgBytesOffset := 0
	for msgBytesOffset < int(expectedLength) {
		var n int
		if n, err = r.Read(msgBytes[msgBytesOffset:]); err != nil {
			return
		}
		msgBytesOffset += n
	}
	msgReader := bytes.NewReader(msgBytes)

	// read beatInfoStopStream magic bytes
	magicBytes := make([]byte, 4)
	if _, err = msgReader.Read(magicBytes); err != nil {
		return
	}
	if !bytes.Equal(magicBytes, beatInfoStopStreamMagicBytes) {
		return errors.New("invalid magic bytes")
	}

	return
}

func (m *beatInfoStopStreamMessage) WriteMessageTo(w io.Writer) (err error) {
	buf := new(bytes.Buffer)

	payload_len := len(beatInfoStopStreamMagicBytes)
	if err = binary.Write(buf, binary.BigEndian, uint32(payload_len)); err != nil {
		return
	}

	// write BeatInfo "stop stream" magic bytes
	if _, err = buf.Write(beatInfoStopStreamMagicBytes); err != nil {
		return
	}

	// send actual message over wire
	_, err = w.Write(buf.Bytes())
	return
}

var beatEmitMagicBytes = []byte{0x0, 0x0, 0x0, 0x2}

type PlayerInfo struct {
	Beat       float64
	TotalBeats float64
	Bpm        float64
}

type beatEmitMessage struct {
	// Length uint32
	// Magic []byte = {0x00,0x00,0x00,0x02}
	Clock     uint64
	Players   []PlayerInfo
	Timelines []float64
}

func (m *beatEmitMessage) CheckMatch(r *bufio.Reader) (ok bool, err error) {
	// peek length bytes and magic bytes
	b, err := r.Peek(4 + 4)
	if err != nil {
		return
	}
	// check magic bytes
	if ok = bytes.Equal(b[4:8], beatEmitMagicBytes); !ok {
		return
	}
	return
}

func (m *beatEmitMessage) ReadMessageFrom(r io.Reader) (err error) {
	// read expected message length
	var expectedLength uint32
	if err = binary.Read(r, binary.BigEndian, &expectedLength); err != nil {
		return
	}

	// set up buffer to write message into
	msgBytes := make([]byte, int(expectedLength))
	msgBytesOffset := 0
	for msgBytesOffset < int(expectedLength) {
		var n int
		if n, err = r.Read(msgBytes[msgBytesOffset:]); err != nil {
			return
		}
		msgBytesOffset += n
	}
	msgReader := bytes.NewReader(msgBytes)

	// read beatEmit magic bytes
	magicBytes := make([]byte, 4)
	if _, err = msgReader.Read(magicBytes); err != nil {
		return
	}
	if !bytes.Equal(magicBytes, beatEmitMagicBytes) {
		err = errors.New("invalid magic bytes")
		return
	}

	// read clock value
	if err = binary.Read(msgReader, binary.BigEndian, &m.Clock); err != nil {
		return
	}

	// read expected player records
	var expectedRecords uint32
	if err = binary.Read(msgReader, binary.BigEndian, &expectedRecords); err != nil {
		return
	}

	// bounds check
	// each playerInfo record is 24 bytes
	if msgReader.Len() < int(expectedRecords)*24 {
		err = errors.New("unknown packet format")
	}

	// loop through players records
	for i := 0; i < int(expectedRecords); i++ {
		var p PlayerInfo
		if err = binary.Read(msgReader, binary.BigEndian, &p.Beat); err != nil {
			return
		}
		if err = binary.Read(msgReader, binary.BigEndian, &p.TotalBeats); err != nil {
			return
		}
		if err = binary.Read(msgReader, binary.BigEndian, &p.Bpm); err != nil {
			return
		}
		m.Players = append(m.Players, p)
	}

	// bounds check
	// the rest of our payload should contain exactly enough bytes for the timeline records
	if msgReader.Len() == int(expectedRecords)*8 {
		err = errors.New("unknown packet format")
	}

	// loop through timelines
	for i := 0; i < int(expectedRecords); i++ {
		var t float64
		if err = binary.Read(msgReader, binary.BigEndian, &t); err != nil {
			return
		}
		m.Timelines = append(m.Timelines, t)
	}

	return
}

func (m *beatEmitMessage) WriteMessageTo(w io.Writer) (err error) {
	// sanity check number of records
	numRecords := len(m.Players)
	if numRecords != len(m.Timelines) {
		err = errors.New("number of player records must match number of timeline records")
		return
	}

	buf := new(bytes.Buffer)

	// write magic bytes to message buffer
	if _, err = buf.Write(beatEmitMagicBytes); err != nil {
		return
	}

	// write clock
	if err = binary.Write(buf, binary.BigEndian, m.Clock); err != nil {
		return
	}

	// write number of records
	if err = binary.Write(buf, binary.BigEndian, uint32(numRecords)); err != nil {
		return
	}

	// write beat info records
	for _, bi := range m.Players {
		if err = binary.Write(buf, binary.BigEndian, bi.Beat); err != nil {
			return
		}
		if err = binary.Write(buf, binary.BigEndian, bi.TotalBeats); err != nil {
			return
		}
		if err = binary.Write(buf, binary.BigEndian, bi.Bpm); err != nil {
			return
		}
	}

	// write timeline records
	for _, tl := range m.Timelines {
		if err = binary.Write(buf, binary.BigEndian, tl); err != nil {
			return
		}
	}

	// send message length over wire
	if err = binary.Write(w, binary.BigEndian, uint32(buf.Len())); err != nil {
		return
	}

	// send actual message over wire
	_, err = w.Write(buf.Bytes())
	return
}

// discovererMessageAction is the action taken by a device as part of StagelinQ device discovery.
// Possible values are DiscovererHowdy or DiscovererExit.
type discovererMessageAction string

const (
	// discovererHowdy is the value set on the Action field of a DiscoveryMessage when a StagelinQ-compatible device announces itself in the network.
	discovererHowdy discovererMessageAction = "DISCOVERER_HOWDY_"

	// discovererExit is the value set on the Action field of a DiscoveryMessage when a StagelinQ-compatible device leaves the network.
	discovererExit discovererMessageAction = "DISCOVERER_EXIT_"
)

// discoveryMessage contains the data carried in the message payload for device trying to handshake the StagelinQ protocol to any other device in the network.
type discoveryMessage struct {
	messages.TokenPrefixedMessage
	Source          string
	Action          discovererMessageAction
	SoftwareName    string
	SoftwareVersion string
	Port            uint16
}

var discoveryMagic = []byte("airD")

func (m *discoveryMessage) CheckMatch(r *bufio.Reader) (ok bool, err error) {
	var readMagic []byte
	if readMagic, err = r.Peek(4); err != nil {
		return
	}
	ok = bytes.Equal(readMagic, discoveryMagic)
	return
}

func (m *discoveryMessage) ReadMessageFrom(r io.Reader) (err error) {
	readMagic := make([]byte, 4)
	if _, err = r.Read(readMagic); err != nil {
		return
	} else if !bytes.Equal(readMagic, discoveryMagic) {
		err = ErrInvalidMessageReceived
		return
	}
	if err = m.TokenPrefixedMessage.ReadMessageFrom(r); err != nil {
		return
	}
	if err = messages.ReadUTF16NetworkString(r, &m.Source); err != nil {
		return
	}
	actionString := ""
	if err = messages.ReadUTF16NetworkString(r, &actionString); err != nil {
		return
	}
	m.Action = discovererMessageAction(actionString)
	if err = messages.ReadUTF16NetworkString(r, &m.SoftwareName); err != nil {
		return
	}
	if err = messages.ReadUTF16NetworkString(r, &m.SoftwareVersion); err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &m.Port)
	return
}

func (m *discoveryMessage) WriteMessageTo(w io.Writer) (err error) {
	if _, err = w.Write(discoveryMagic); err != nil {
		return
	}
	if err = m.TokenPrefixedMessage.WriteMessageTo(w); err != nil {
		return
	}
	if err = messages.WriteUTF16NetworkString(w, m.Source); err != nil {
		return
	}
	if err = messages.WriteUTF16NetworkString(w, string(m.Action)); err != nil {
		return
	}
	if err = messages.WriteUTF16NetworkString(w, m.SoftwareName); err != nil {
		return
	}
	if err = messages.WriteUTF16NetworkString(w, m.SoftwareVersion); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, m.Port); err != nil {
		return
	}
	return
}
