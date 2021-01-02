package stagelinq

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

// ErrUnknownMessageID indicates that an invalid message was received from the connected device.
var ErrUnknownMessageID = errors.New("unknown message ID")

type deviceConn struct {
	net.Conn
}

func newDeviceConn(conn net.Conn) *deviceConn {
	return &deviceConn{conn}
}

func (s *deviceConn) Close() error {
	return s.Conn.Close()
}

func (s *deviceConn) WriteMessage(msg message) (err error) {
	buf := new(bytes.Buffer)

	// write message id (4 bytes)
	if err = binary.Write(buf, binary.BigEndian, msg.id()); err != nil {
		return
	}

	// write message itself
	if err = msg.writeTo(buf); err != nil {
		return
	}

	// write the whole thing out to the device
	_, err = s.Conn.Write(buf.Bytes())
	return
}

func (s *deviceConn) ReadMessage() (msg message, err error) {
	// read message id (4 bytes)
	var messageID int32
	if err = binary.Read(s.Conn, binary.BigEndian, &messageID); err != nil {
		return
	}

	// find associated function that creates message object
	messageObjectGeneratorFunction, ok := tcpMessageMap[messageID]
	if !ok {
		err = ErrUnknownMessageID
		err = fmt.Errorf("%s: %x", ErrUnknownMessageID.Error(), messageID)
		return
	}

	// create message object and decode message from device
	msg = messageObjectGeneratorFunction()
	if err = msg.readFrom(s.Conn); err != nil {
		return
	}

	return
}
