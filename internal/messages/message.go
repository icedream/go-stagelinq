package messages

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

type Message interface {
	ReadMessageFrom(io.Reader) error
	WriteMessageTo(io.Writer) error

	// CheckMatch MUST use the Peek method to read any bytes needed to exactly
	// identify whether a message matches.
	//
	// It SHOULD not peek more bytes than are necessary to identify the message.
	//
	// The method MUST avoid Read to allow other message types to validate the
	// message properly.
	CheckMatch(*bufio.Reader) (bool, error)
}

func ReadMessageID(r io.Reader) (id int32, err error) {
	err = binary.Read(r, binary.BigEndian, &id)
	return
}

func PeekMessageID(r *bufio.Reader) (id int32, err error) {
	b, err := r.Peek(4)
	if err != nil {
		return
	}
	return ReadMessageID(bytes.NewReader(b))
}

func WriteMessageID(w io.Writer, id int32) (err error) {
	err = binary.Write(w, binary.BigEndian, id)
	return
}
