package stagelinq

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"reflect"

	"github.com/icedream/go-stagelinq/internal/messages"
)

type messageSet struct {
	messages []reflect.Type
}

func newDeviceConnMessageSet(messageObjects []messages.Message) *messageSet {
	messages := make([]reflect.Type, len(messageObjects))
	for i, messageObject := range messageObjects {
		// .Elem() because type will be a pointer-to-type but we want to create instances of the type itself later
		messages[i] = reflect.TypeOf(messageObject).Elem()
	}
	return &messageSet{messages}
}

func (ms *messageSet) Messages() []reflect.Type {
	return ms.messages
}

type messageConnection struct {
	conn             net.Conn
	bufferedReader   *bufio.Reader
	expectedMessages *messageSet
}

func newMessageConnection(conn net.Conn, expectedMessages *messageSet) *messageConnection {
	if conn == nil {
		panic("conn must not be nil")
	}
	if expectedMessages == nil {
		panic("expectedMessages must not be nil")
	}
	if len(expectedMessages.Messages()) <= 0 {
		panic("expectedMessages must not be empty")
	}
	return &messageConnection{
		conn:             conn,
		bufferedReader:   bufio.NewReader(conn),
		expectedMessages: expectedMessages,
	}
}

func (s *messageConnection) WriteMessage(msg messages.Message) (err error) {
	buf := new(bytes.Buffer)

	// write message parts into buffer
	if err = msg.WriteMessageTo(buf); err != nil {
		return
	}

	// write the whole thing out as one message to the device
	_, err = s.conn.Write(buf.Bytes())

	// if err == nil {
	// 	log.Printf("SEND: %s", spew.Sdump(msg))
	// }

	return
}

func (s *messageConnection) ReadMessage() (msg messages.Message, err error) {
	var targetMsg messages.Message
	var ok bool
	for _, messageType := range s.expectedMessages.Messages() {
		targetMsg = reflect.New(messageType).Interface().(messages.Message)
		ok, err = targetMsg.CheckMatch(s.bufferedReader)
		if err != nil {
			return
		}
		if ok {
			break
		}
	}

	if !ok {
		b, _ := s.bufferedReader.Peek(s.bufferedReader.Buffered())
		err = fmt.Errorf("%w: buffered bytes:\n%s", ErrInvalidMessageReceived, hex.Dump(b))
		return
	}

	err = targetMsg.ReadMessageFrom(s.bufferedReader)
	if err == nil {
		msg = targetMsg
		// log.Printf("RECV: %s", spew.Sdump(msg))
	}

	return
}
