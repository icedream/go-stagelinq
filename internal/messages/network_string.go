package messages

import (
	"encoding/binary"
	"io"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

var (
	UTF16 = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	UTF8  = unicode.UTF8
)

func WriteNetworkStringWithEncoding(w io.Writer, v string, enc encoding.Encoding) (err error) {
	converted, err := enc.NewEncoder().Bytes([]byte(v))
	if err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint32(len(converted))); err != nil {
		return
	}
	_, err = w.Write(converted)
	return
}

func ReadNetworkStringWithEncoding(r io.Reader, v *string, enc encoding.Encoding) (err error) {
	var expectedLength uint32
	if err = binary.Read(r, binary.BigEndian, &expectedLength); err != nil {
		return
	}
	expectedLengthInt := int(expectedLength)
	buf := make([]byte, expectedLengthInt)
	offset := 0
	// I don't know if this is necessary, but I'm doing it anyways because #DefensiveProgrammingIsGood
	for offset < expectedLengthInt {
		var n int
		n, err = r.Read(buf[offset:])
		if err != nil {
			return
		}
		offset += n
	}
	vBytes, err := enc.NewDecoder().Bytes(buf)
	if err != nil {
		return
	}
	*v = string(vBytes)
	return
}

func WriteUTF16NetworkString(w io.Writer, v string) (err error) {
	return WriteNetworkStringWithEncoding(w, v, UTF16)
}

func ReadUTF16NetworkString(r io.Reader, v *string) (err error) {
	return ReadNetworkStringWithEncoding(r, v, UTF16)
}
