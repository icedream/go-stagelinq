package stagelinq

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// Token contains the identifying Token for a device in the StagelinQ network.
type Token [16]byte

type message interface {
	readFrom(io.Reader) error
	writeTo(io.Writer) error

	// checkMatch MUST use the Peek method to read any bytes needed to exactly identify whether a message matches.
	// It SHOULD not peek more bytes than are necessary to identify the message.
	// The method MUST avoid Read to allow other message types to validate the message properly.
	checkMatch(*bufio.Reader) (bool, error)
}

func readMessageID(r io.Reader) (id int32, err error) {
	err = binary.Read(r, binary.BigEndian, &id)
	return
}

func peekMessageID(r *bufio.Reader) (id int32, err error) {
	b, err := r.Peek(4)
	if err != nil {
		return
	}
	return readMessageID(bytes.NewReader(b))
}

func writeMessageID(w io.Writer, id int32) (err error) {
	err = binary.Write(w, binary.BigEndian, id)
	return
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

type serviceAnnouncementMessage struct {
	tokenPrefixedMessage
	Service string
	Port    uint16
}

func (m *serviceAnnouncementMessage) checkMatch(r *bufio.Reader) (ok bool, err error) {
	id, err := peekMessageID(r)
	if err != nil {
		return
	}
	ok = id == 0x00000000
	return
}

func (m *serviceAnnouncementMessage) readFrom(r io.Reader) (err error) {
	messageID, err := readMessageID(r)
	if err != nil {
		return
	} else if messageID != 0x00000000 {
		err = ErrInvalidMessageReceived
		return
	}
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

func (m *serviceAnnouncementMessage) writeTo(w io.Writer) (err error) {
	if err = writeMessageID(w, 0x00000000); err != nil {
		return
	}
	if err = m.tokenPrefixedMessage.writeTo(w); err != nil {
		return
	}
	if err = writeNetworkString(w, m.Service); err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, m.Port)
	return
}

type referenceMessage struct {
	tokenPrefixedMessage
	Token2    Token
	Reference int64
}

func (m *referenceMessage) checkMatch(r *bufio.Reader) (ok bool, err error) {
	id, err := peekMessageID(r)
	if err != nil {
		return
	}
	ok = id == 0x00000001
	return
}

func (m *referenceMessage) readFrom(r io.Reader) (err error) {
	messageID, err := readMessageID(r)
	if err != nil {
		return
	} else if messageID != 0x00000001 {
		err = ErrInvalidMessageReceived
		return
	}
	if err = m.tokenPrefixedMessage.readFrom(r); err != nil {
		return
	}
	if _, err = r.Read(m.Token2[:]); err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &m.Reference)
	return
}

func (m *referenceMessage) writeTo(w io.Writer) (err error) {
	if err = writeMessageID(w, 0x00000001); err != nil {
		return
	}
	if err = m.tokenPrefixedMessage.writeTo(w); err != nil {
		return
	}
	if _, err = w.Write(m.Token2[:]); err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, m.Reference)
	return
}

type servicesRequestMessage struct {
	tokenPrefixedMessage
}

func (m *servicesRequestMessage) checkMatch(r *bufio.Reader) (ok bool, err error) {
	id, err := peekMessageID(r)
	if err != nil {
		return
	}
	ok = id == 0x00000002
	return
}

func (m *servicesRequestMessage) readFrom(r io.Reader) (err error) {
	messageID, err := readMessageID(r)
	if err != nil {
		return
	} else if messageID != 0x00000002 {
		err = ErrInvalidMessageReceived
		return
	}

	err = m.tokenPrefixedMessage.readFrom(r)
	return
}

func (m *servicesRequestMessage) writeTo(w io.Writer) (err error) {
	if err = writeMessageID(w, 0x00000002); err != nil {
		return
	}

	err = m.tokenPrefixedMessage.writeTo(w)
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
	//Length uint32
	//Unknown []byte = {0x73,0x6d,0x61,0x61}
	//Unknown2 []byte = {0x00,0x00,0x07,0xd2}
	Name     string
	Interval uint32
}

func (m *stateSubscribeMessage) checkMatch(r *bufio.Reader) (ok bool, err error) {
	return checkSmaa(r, 0x000007d2)
}

func (m *stateSubscribeMessage) readFrom(r io.Reader) (err error) {
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
	if err = readNetworkString(r, &m.Name); err != nil {
		return
	}

	// TODO - figure this out
	err = binary.Read(r, binary.BigEndian, &m.Interval)
	return
}

func (m *stateSubscribeMessage) writeTo(w io.Writer) (err error) {
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
	if err = writeNetworkString(buf, m.Name); err != nil {
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
	//Length uint32
	//Unknown []byte = {0x73,0x6d,0x61,0x61}
	//Unknown2 []byte = {0x00,0x00,0x00,0x00}
	Name string
	JSON string
}

func (m *stateEmitMessage) checkMatch(r *bufio.Reader) (ok bool, err error) {
	return checkSmaa(r, 0x00000000)
}

func (m *stateEmitMessage) readFrom(r io.Reader) (err error) {
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
	if err = readNetworkString(msgReader, &m.Name); err != nil {
		return
	}

	// read value JSON
	if err = readNetworkString(msgReader, &m.JSON); err != nil {
		return
	}

	return
}

func (m *stateEmitMessage) writeTo(w io.Writer) (err error) {
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
	if err = writeNetworkString(buf, m.Name); err != nil {
		return
	}

	// write value JSON to message buffer
	if err = writeNetworkString(buf, m.JSON); err != nil {
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

type beatInfoStartStreamMessage struct {
}

func (m *beatInfoStartStreamMessage) checkMatch(r *bufio.Reader) (ok bool, err error) {
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

func (m *beatInfoStartStreamMessage) readFrom(r io.Reader) (err error) {
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

func (m *beatInfoStartStreamMessage) writeTo(w io.Writer) (err error) {
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

type beatInfoStopStreamMessage struct {
}

func (m *beatInfoStopStreamMessage) checkMatch(r *bufio.Reader) (ok bool, err error) {
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

func (m *beatInfoStopStreamMessage) readFrom(r io.Reader) (err error) {
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

func (m *beatInfoStopStreamMessage) writeTo(w io.Writer) (err error) {
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
	//Length uint32
	//Magic []byte = {0x00,0x00,0x00,0x02}
	Clock     uint64
	Players   []PlayerInfo
	Timelines []float64
}

func (m *beatEmitMessage) checkMatch(r *bufio.Reader) (ok bool, err error) {
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

func (m *beatEmitMessage) readFrom(r io.Reader) (err error) {
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

func (m *beatEmitMessage) writeTo(w io.Writer) (err error) {
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
	tokenPrefixedMessage
	Source          string
	Action          discovererMessageAction
	SoftwareName    string
	SoftwareVersion string
	Port            uint16
}

var discoveryMagic = []byte("airD")

func (m *discoveryMessage) checkMatch(r *bufio.Reader) (ok bool, err error) {
	var readMagic []byte
	if readMagic, err = r.Peek(4); err != nil {
		return
	}
	ok = bytes.Equal(readMagic, discoveryMagic)
	return
}

func (m *discoveryMessage) readFrom(r io.Reader) (err error) {
	readMagic := make([]byte, 4)
	if _, err = r.Read(readMagic); err != nil {
		return
	} else if !bytes.Equal(readMagic, discoveryMagic) {
		err = ErrInvalidMessageReceived
		return
	}
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
	m.Action = discovererMessageAction(actionString)
	if err = readNetworkString(r, &m.SoftwareName); err != nil {
		return
	}
	if err = readNetworkString(r, &m.SoftwareVersion); err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &m.Port)
	return
}

func (m *discoveryMessage) writeTo(w io.Writer) (err error) {
	if _, err = w.Write(discoveryMagic); err != nil {
		return
	}
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
