package messages

import "io"

// Token contains the identifying Token for a device in the StagelinQ network.
type Token [16]byte

type TokenPrefixedMessage struct {
	Token Token
}

func (m *TokenPrefixedMessage) ReadMessageFrom(r io.Reader) (err error) {
	_, err = r.Read(m.Token[:])
	return
}

func (m *TokenPrefixedMessage) WriteMessageTo(w io.Writer) (err error) {
	_, err = w.Write(m.Token[:])
	return
}
