package stagelinq

import (
	"net"
	"testing"

	"github.com/icedream/go-stagelinq/internal/messages"
	"github.com/stretchr/testify/require"
)

var testToken = messages.Token{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

func Test_MessageConnection_Read(t *testing.T) {
	testMessages := []struct {
		Message messages.Message
		Bytes   []byte
	}{
		{
			Message: &serviceAnnouncementMessage{
				TokenPrefixedMessage: messages.TokenPrefixedMessage{
					Token: messages.Token{0xf4, 0x05, 0xdc, 0x14, 0x02, 0x23, 0x47, 0xf5, 0x8b, 0x79, 0x2c, 0x8c, 0x49, 0x33, 0x52, 0x76},
				},
				Service: "StateMap",
				Port:    0xb1d7,
			},
			Bytes: []byte{
				/*
					from Prime4
					00000014  00 00 00 00 f4 05 dc 14  02 23 47 f5 8b 79 2c 8c   ........ .#G..y,.
					00000024  49 33 52 76 00 00 00 10  00 53 00 74 00 61 00 74   I3Rv.... .S.t.a.t
					00000034  00 65 00 4d 00 61 00 70  b1 d7                     .e.M.a.p ..
				*/
				0x00, 0x00, 0x00, 0x00, 0xf4, 0x05, 0xdc, 0x14,
				0x02, 0x23, 0x47, 0xf5, 0x8b, 0x79, 0x2c, 0x8c,
				0x49, 0x33, 0x52, 0x76, 0x00, 0x00, 0x00, 0x10,
				0x00, 0x53, 0x00, 0x74, 0x00, 0x61, 0x00, 0x74,
				0x00, 0x65, 0x00, 0x4d, 0x00, 0x61, 0x00, 0x70,
				0xb1, 0xd7,
			},
		},
		{
			Message: &referenceMessage{
				TokenPrefixedMessage: messages.TokenPrefixedMessage{
					Token: messages.Token{0xf4, 0x05, 0xdc, 0x14, 0x02, 0x23, 0x47, 0xf5, 0x8b, 0x79, 0x2c, 0x8c, 0x49, 0x33, 0x52, 0x76},
				},
				Reference: 0x000009ed4f310604,
			},
			Bytes: []byte{
				/*
					from Prime4
					00000106  00 00 00 01 f4 05 dc 14  02 23 47 f5 8b 79 2c 8c   ........ .#G..y,.
					00000116  49 33 52 76 00 00 00 00  00 00 00 00 00 00 00 00   I3Rv.... ........
					00000126  00 00 00 00 00 00 09 ed  4f 31 06 04               ........ O1..
				*/
				0x00, 0x00, 0x00, 0x01, 0xf4, 0x05, 0xdc, 0x14,
				0x02, 0x23, 0x47, 0xf5, 0x8b, 0x79, 0x2c, 0x8c,
				0x49, 0x33, 0x52, 0x76, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09, 0xed,
				0x4f, 0x31, 0x06, 0x04,
			},
		},
		{
			Message: &servicesRequestMessage{
				TokenPrefixedMessage: messages.TokenPrefixedMessage{
					Token: messages.Token{0xf4, 0x05, 0xdc, 0x14, 0x02, 0x23, 0x47, 0xf5, 0x8b, 0x79, 0x2c, 0x8c, 0x49, 0x33, 0x52, 0x76},
				},
			},
			Bytes: []byte{
				/* from Prime4 */
				0x00, 0x00, 0x00, 0x02, 0xf4, 0x05, 0xdc, 0x14,
				0x02, 0x23, 0x47, 0xf5, 0x8b, 0x79, 0x2c, 0x8c,
				0x49, 0x33, 0x52, 0x76,
			},
		},
		{
			Message: &serviceAnnouncementMessage{
				TokenPrefixedMessage: messages.TokenPrefixedMessage{
					Token: messages.Token{0x52, 0x3e, 0x67, 0x9d, 0xa4, 0x18, 0x4d, 0x1e, 0x83, 0xd0, 0xc7, 0x52, 0xcf, 0xca, 0x8f, 0xf7},
				},
				Service: "DirectoryService",
				Port:    0xe190,
			},
			Bytes: []byte{
				/*
					from Resolume
					00000000  00 00 00 00 52 3e 67 9d  a4 18 4d 1e 83 d0 c7 52   ....R>g. ..M....R
					00000010  cf ca 8f f7 00 00 00 20  00 44 00 69 00 72 00 65   .......  .D.i.r.e
					00000020  00 63 00 74 00 6f 00 72  00 79 00 53 00 65 00 72   .c.t.o.r .y.S.e.r
					00000030  00 76 00 69 00 63 00 65  e1 90                     .v.i.c.e ..
				*/
				0x00, 0x00, 0x00, 0x00, 0x52, 0x3e, 0x67, 0x9d,
				0xa4, 0x18, 0x4d, 0x1e, 0x83, 0xd0, 0xc7, 0x52,
				0xcf, 0xca, 0x8f, 0xf7, 0x00, 0x00, 0x00, 0x20,
				0x00, 0x44, 0x00, 0x69, 0x00, 0x72, 0x00, 0x65,
				0x00, 0x63, 0x00, 0x74, 0x00, 0x6f, 0x00, 0x72,
				0x00, 0x79, 0x00, 0x53, 0x00, 0x65, 0x00, 0x72,
				0x00, 0x76, 0x00, 0x69, 0x00, 0x63, 0x00, 0x65,
				0xe1, 0x90,
			},
		},
		{
			Message: &servicesRequestMessage{
				TokenPrefixedMessage: messages.TokenPrefixedMessage{
					Token: messages.Token{0x52, 0x3e, 0x67, 0x9d, 0xa4, 0x18, 0x4d, 0x1e, 0x83, 0xd0, 0xc7, 0x52, 0xcf, 0xca, 0x8f, 0xf7},
				},
			},
			Bytes: []byte{
				/*
					from Resolume
					0000003A  00 00 00 02 52 3e 67 9d  a4 18 4d 1e 83 d0 c7 52   ....R>g. ..M....R
					0000004A  cf ca 8f f7                                        ....
				*/
				0x00, 0x00, 0x00, 0x02, 0x52, 0x3e, 0x67, 0x9d,
				0xa4, 0x18, 0x4d, 0x1e, 0x83, 0xd0, 0xc7, 0x52,
				0xcf, 0xca, 0x8f, 0xf7,
			},
		},
		{
			Message: &referenceMessage{
				Token2:    messages.Token{0x52, 0x3e, 0x67, 0x9d, 0xa4, 0x18, 0x4d, 0x1e, 0x83, 0xd0, 0xc7, 0x52, 0xcf, 0xca, 0x8f, 0xf7},
				Reference: 0x000009ed4f310604,
			},
			Bytes: []byte{
				/*
					from Resolume
					0000004E  00 00 00 01 00 00 00 00  00 00 00 00 00 00 00 00   ........ ........
					0000005E  00 00 00 00 52 3e 67 9d  a4 18 4d 1e 83 d0 c7 52   ....R>g. ..M....R
					0000006E  cf ca 8f f7 00 00 09 ed  4f 31 06 04               ........ O1..
				*/
				0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x52, 0x3e, 0x67, 0x9d,
				0xa4, 0x18, 0x4d, 0x1e, 0x83, 0xd0, 0xc7, 0x52,
				0xcf, 0xca, 0x8f, 0xf7, 0x00, 0x00, 0x09, 0xed,
				0x4f, 0x31, 0x06, 0x04,
			},
		},
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Can't set up test listener: %s", err.Error())
		return
	}

	{
		conn, err := net.Dial("tcp", listener.Addr().String())
		if err != nil {
			t.Fatalf("Failed to set up test connection: %s", err.Error())
			return
		}
		go func() {
			for _, testMessage := range testMessages {
				_, err := conn.Write(testMessage.Bytes)
				if err != nil {
					t.Logf("Failed to write bytes: %s", err.Error())
				}
			}
		}()
	}

	conn, err := listener.Accept()
	if err != nil {
		t.Fatalf("Failed to accept test connection: %s", err.Error())
	}

	messageObjects := []messages.Message{}
	for _, testMessage := range testMessages {
		messageObjects = append(messageObjects, testMessage.Message)
	}
	msgConn := newMessageConnection(conn, newDeviceConnMessageSet(messageObjects))

	for _, expectedMessage := range testMessages {
		message, err := msgConn.ReadMessage()
		require.Nil(t, err)
		require.Equal(t, expectedMessage.Message, message)
	}
}

func Test_MessageConnection(t *testing.T) {
	testMessages := []messages.Message{
		&serviceAnnouncementMessage{
			TokenPrefixedMessage: messages.TokenPrefixedMessage{
				Token: testToken,
			},
			Service: "test",
			Port:    0x1234,
		},
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Can't set up test listener: %s", err.Error())
		return
	}

	{
		conn, err := net.Dial("tcp", listener.Addr().String())
		if err != nil {
			t.Fatalf("Failed to set up test connection: %s", err.Error())
			return
		}
		go func() {
			msgConn := newMessageConnection(conn, newDeviceConnMessageSet(testMessages))

			for _, testMessage := range testMessages {
				err := msgConn.WriteMessage(testMessage)
				require.Nil(t, err)
			}
		}()
	}

	conn, err := listener.Accept()
	if err != nil {
		t.Fatalf("Failed to accept test connection: %s", err.Error())
	}
	msgConn := newMessageConnection(conn, newDeviceConnMessageSet(testMessages))
	for _, expectedMessage := range testMessages {
		message, err := msgConn.ReadMessage()
		require.Nil(t, err)
		require.Equal(t, expectedMessage, message)
	}
}
