package eaas

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/icedream/go-stagelinq/internal/messages"
)

var eaasDiscoveryMagic = []byte("EAAS")

func checkEAASDiscoveryMagic(r *bufio.Reader, req [2]byte) (ok bool, err error) {
	var readMagic []byte
	if readMagic, err = r.Peek(6); err != nil {
		return
	}
	ok = bytes.Equal(readMagic, append(eaasDiscoveryMagic, req[:]...))
	return
}

type eaasDiscoveryRequestMessage struct{}

func (m *eaasDiscoveryRequestMessage) CheckMatch(r *bufio.Reader) (ok bool, err error) {
	return checkEAASDiscoveryMagic(r, [2]byte{1, 0})
}

func (m *eaasDiscoveryRequestMessage) ReadMessageFrom(r io.Reader) (err error) {
	readMagic := make([]byte, 6)
	if _, err = r.Read(readMagic); err != nil {
		return err
	} else if !bytes.Equal(readMagic, append(eaasDiscoveryMagic, 1, 0)) {
		err = ErrInvalidMessageReceived
	}
	return err
}

func (m *eaasDiscoveryRequestMessage) WriteMessageTo(w io.Writer) (err error) {
	_, err = w.Write(append(eaasDiscoveryMagic, 1, 0))
	return err
}

type eaasDiscoveryResponseMessage struct {
	messages.TokenPrefixedMessage
	Hostname        string
	SoftwareVersion string
	URL             string
	Extra           string // usually just _
}

func (m *eaasDiscoveryResponseMessage) CheckMatch(r *bufio.Reader) (ok bool, err error) {
	return checkEAASDiscoveryMagic(r, [2]byte{1, 1})
}

func (m *eaasDiscoveryResponseMessage) ReadMessageFrom(r io.Reader) (err error) {
	readMagic := make([]byte, 6)
	if _, err = r.Read(readMagic); err != nil {
		return err
	} else if !bytes.Equal(readMagic, append(eaasDiscoveryMagic, 1, 1)) {
		err = ErrInvalidMessageReceived
		return err
	}
	if err = m.TokenPrefixedMessage.ReadMessageFrom(r); err != nil {
		return fmt.Errorf("failed to read token: %w", err)
	}
	if err = messages.ReadUTF16NetworkString(r, &m.Hostname); err != nil {
		return fmt.Errorf("failed to read hostname string: %w", err)
	}
	if err = messages.ReadNetworkStringWithEncoding(r, &m.URL, messages.UTF8); err != nil {
		return fmt.Errorf("failed to read URL string: %w", err)
	}
	if err = messages.ReadUTF16NetworkString(r, &m.SoftwareVersion); err != nil {
		return fmt.Errorf("failed to read software version string: %w", err)
	}
	// TODO - there is an extra 0x01 here which idk what to do with
	if _, err = r.Read([]byte{0x01}); err != nil {
		return err
	}
	// TODO - is this really a string?
	if err = messages.ReadUTF16NetworkString(r, &m.Extra); err != nil {
		return fmt.Errorf("failed to read extra network string: %w", err)
	}
	return err
}

func (m *eaasDiscoveryResponseMessage) WriteMessageTo(w io.Writer) (err error) {
	if _, err = w.Write(append(eaasDiscoveryMagic, 1, 1)); err != nil {
		return fmt.Errorf("failed to write EAAS magic: %w", err)
	}
	if err = m.TokenPrefixedMessage.WriteMessageTo(w); err != nil {
		return fmt.Errorf("failed to write token: %w", err)
	}
	if err = messages.WriteUTF16NetworkString(w, m.Hostname); err != nil {
		return fmt.Errorf("failed to write hostname string: %w", err)
	}
	if err = messages.WriteNetworkStringWithEncoding(w, m.URL, messages.UTF8); err != nil {
		return fmt.Errorf("failed to write URL string: %w", err)
	}
	if err = messages.WriteUTF16NetworkString(w, m.SoftwareVersion); err != nil {
		return fmt.Errorf("failed to write software version string: %w", err)
	}
	if _, err = w.Write([]byte{0x01}); err != nil {
		return err
	}
	if err = messages.WriteUTF16NetworkString(w, m.Extra); err != nil {
		return fmt.Errorf("failed to write extra network string: %w", err)
	}
	return err
}
