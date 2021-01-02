package stagelinq

import (
	"encoding/binary"
	"io"

	"golang.org/x/text/encoding/unicode"
)

var networkStringEncoding = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)

func writeNetworkString(w io.Writer, v string) (err error) {
	converted, err := networkStringEncoding.NewEncoder().Bytes([]byte(v))
	if err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint32(len(converted))); err != nil {
		return
	}
	_, err = w.Write(converted)
	return
}

func readNetworkString(r io.Reader, v *string) (err error) {
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
	vBytes, err := networkStringEncoding.NewDecoder().Bytes(buf)
	if err != nil {
		return
	}
	*v = string(vBytes)
	return
}
