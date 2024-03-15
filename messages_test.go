package stagelinq

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

var testMessages = []struct {
	Name          string
	Message       message
	CreateMessage func() message
	Bytes         []byte
}{
	{
		Name: "Discovery",
		Bytes: []byte{
			0x61, 0x69, 0x72, 0x44, 0xf4, 0x05, 0xdc, 0x14,
			0x02, 0x23, 0x47, 0xf5, 0x8b, 0x79, 0x2c, 0x8c,
			0x49, 0x33, 0x52, 0x76, 0x00, 0x00, 0x00, 0x0c,
			0x00, 0x70, 0x00, 0x72, 0x00, 0x69, 0x00, 0x6d,
			0x00, 0x65, 0x00, 0x34, 0x00, 0x00, 0x00, 0x22,
			0x00, 0x44, 0x00, 0x49, 0x00, 0x53, 0x00, 0x43,
			0x00, 0x4f, 0x00, 0x56, 0x00, 0x45, 0x00, 0x52,
			0x00, 0x45, 0x00, 0x52, 0x00, 0x5f, 0x00, 0x48,
			0x00, 0x4f, 0x00, 0x57, 0x00, 0x44, 0x00, 0x59,
			0x00, 0x5f, 0x00, 0x00, 0x00, 0x08, 0x00, 0x4a,
			0x00, 0x43, 0x00, 0x31, 0x00, 0x31, 0x00, 0x00,
			0x00, 0x0a, 0x00, 0x31, 0x00, 0x2e, 0x00, 0x35,
			0x00, 0x2e, 0x00, 0x32, 0x84, 0x03,
		},
		CreateMessage: func() message { return new(discoveryMessage) },
		Message: &discoveryMessage{
			tokenPrefixedMessage: tokenPrefixedMessage{
				Token: Token{0xf4, 0x05, 0xdc, 0x14, 0x02, 0x23, 0x47, 0xf5, 0x8b, 0x79, 0x2c, 0x8c, 0x49, 0x33, 0x52, 0x76},
			},
			Source:          "prime4",
			Action:          discovererHowdy,
			SoftwareName:    "JC11",
			SoftwareVersion: "1.5.2",
			Port:            0x8403,
		},
	},
	{
		Name: "Service announcement",
		Bytes: []byte{
			0x00, 0x00, 0x00, 0x00, 0x52, 0x3e, 0x67, 0x9d,
			0xa4, 0x18, 0x4d, 0x1e, 0x83, 0xd0, 0xc7, 0x52,
			0xcf, 0xca, 0x8f, 0xf7, 0x00, 0x00, 0x00, 0x10,
			0x00, 0x53, 0x00, 0x74, 0x00, 0x61, 0x00, 0x74,
			0x00, 0x65, 0x00, 0x4d, 0x00, 0x61, 0x00, 0x70,
			0xe1, 0x96,
		},
		CreateMessage: func() message { return new(serviceAnnouncementMessage) },
		Message: &serviceAnnouncementMessage{
			tokenPrefixedMessage: tokenPrefixedMessage{
				Token: Token{0x52, 0x3e, 0x67, 0x9d, 0xa4, 0x18, 0x4d, 0x1e, 0x83, 0xd0, 0xc7, 0x52, 0xcf, 0xca, 0x8f, 0xf7},
			},
			Service: "StateMap",
			Port:    0xe196,
		},
	},
	{
		Name: "Services request",
		Bytes: []byte{
			0x00, 0x00, 0x00, 0x02, 0xf4, 0x05, 0xdc, 0x14,
			0x02, 0x23, 0x47, 0xf5, 0x8b, 0x79, 0x2c, 0x8c,
			0x49, 0x33, 0x52, 0x76,
		},
		CreateMessage: func() message { return new(servicesRequestMessage) },
		Message: &servicesRequestMessage{
			tokenPrefixedMessage: tokenPrefixedMessage{
				Token: Token{0xf4, 0x05, 0xdc, 0x14, 0x02, 0x23, 0x47, 0xf5, 0x8b, 0x79, 0x2c, 0x8c, 0x49, 0x33, 0x52, 0x76},
			},
		},
	},
	{
		Name: "Reference",
		Bytes: []byte{
			0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0xf4, 0x05, 0xdc, 0x14,
			0x02, 0x23, 0x47, 0xf5, 0x8b, 0x79, 0x2c, 0x8c,
			0x49, 0x33, 0x52, 0x76, 0x00, 0x00, 0x09, 0xed,
			0x4f, 0x31, 0x06, 0x04,
		},
		CreateMessage: func() message { return new(referenceMessage) },
		Message: &referenceMessage{
			tokenPrefixedMessage: tokenPrefixedMessage{
				Token: Token{},
			},
			Token2:    Token{0xf4, 0x05, 0xdc, 0x14, 0x02, 0x23, 0x47, 0xf5, 0x8b, 0x79, 0x2c, 0x8c, 0x49, 0x33, 0x52, 0x76},
			Reference: 0x000009ed4f310604,
		},
	},
	// {
	// 	Bytes: []byte{
	// 		0x00, 0x00, 0x00, 0x00, 0x06, 0xd2, 0x3b, 0xe4,
	// 		0x8e, 0xb2, 0x4f, 0xc7, 0x8f, 0x03, 0xc6, 0xcc,
	// 		0x70, 0x70, 0x30, 0x1b, 0x00, 0x00, 0x00, 0x10,
	// 		0x00, 0x53, 0x00, 0x74, 0x00, 0x61, 0x00, 0x74,
	// 		0x00, 0x65, 0x00, 0x4d, 0x00, 0x61, 0x00, 0x70,
	// 		0xff, 0xbd,
	// 	},
	// },
	{
		Name: "State subscribe",
		Bytes: []byte{
			0x00, 0x00, 0x00, 0x44, 0x73, 0x6d, 0x61, 0x61,
			0x00, 0x00, 0x07, 0xd2, 0x00, 0x00, 0x00, 0x34,
			0x00, 0x2f, 0x00, 0x43, 0x00, 0x6c, 0x00, 0x69,
			0x00, 0x65, 0x00, 0x6e, 0x00, 0x74, 0x00, 0x2f,
			0x00, 0x50, 0x00, 0x72, 0x00, 0x65, 0x00, 0x66,
			0x00, 0x65, 0x00, 0x72, 0x00, 0x65, 0x00, 0x6e,
			0x00, 0x63, 0x00, 0x65, 0x00, 0x73, 0x00, 0x2f,
			0x00, 0x50, 0x00, 0x6c, 0x00, 0x61, 0x00, 0x79,
			0x00, 0x65, 0x00, 0x72, 0x00, 0x00, 0x00, 0x00,
		},
		CreateMessage: func() message { return new(stateSubscribeMessage) },
		Message: &stateSubscribeMessage{
			Name: ClientPreferencesPlayer,
		},
	},
	{
		Name: "State emit",
		Bytes: []byte{
			0x00, 0x00, 0x00, 0x72, 0x73, 0x6d, 0x61, 0x61,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34,
			0x00, 0x2f, 0x00, 0x43, 0x00, 0x6c, 0x00, 0x69,
			0x00, 0x65, 0x00, 0x6e, 0x00, 0x74, 0x00, 0x2f,
			0x00, 0x50, 0x00, 0x72, 0x00, 0x65, 0x00, 0x66,
			0x00, 0x65, 0x00, 0x72, 0x00, 0x65, 0x00, 0x6e,
			0x00, 0x63, 0x00, 0x65, 0x00, 0x73, 0x00, 0x2f,
			0x00, 0x50, 0x00, 0x6c, 0x00, 0x61, 0x00, 0x79,
			0x00, 0x65, 0x00, 0x72, 0x00, 0x00, 0x00, 0x2e,
			0x00, 0x7b, 0x00, 0x22, 0x00, 0x73, 0x00, 0x74,
			0x00, 0x72, 0x00, 0x69, 0x00, 0x6e, 0x00, 0x67,
			0x00, 0x22, 0x00, 0x3a, 0x00, 0x22, 0x00, 0x31,
			0x00, 0x22, 0x00, 0x2c, 0x00, 0x22, 0x00, 0x74,
			0x00, 0x79, 0x00, 0x70, 0x00, 0x65, 0x00, 0x22,
			0x00, 0x3a, 0x00, 0x34, 0x00, 0x7d,
		},
		CreateMessage: func() message { return new(stateEmitMessage) },
		Message: &stateEmitMessage{
			Name: ClientPreferencesPlayer,
			JSON: `{"string":"1","type":4}`,
		},
	},
	{
		Name: "State emit response",
		Bytes: []byte{
			0x00, 0x00, 0x00, 0x44, 0x73, 0x6d, 0x61, 0x61,
			0x00, 0x00, 0x07, 0xd1, 0x00, 0x00, 0x00, 0x34,
			0x00, 0x2f, 0x00, 0x43, 0x00, 0x6c, 0x00, 0x69,
			0x00, 0x65, 0x00, 0x6e, 0x00, 0x74, 0x00, 0x2f,
			0x00, 0x50, 0x00, 0x72, 0x00, 0x65, 0x00, 0x66,
			0x00, 0x65, 0x00, 0x72, 0x00, 0x65, 0x00, 0x6e,
			0x00, 0x63, 0x00, 0x65, 0x00, 0x73, 0x00, 0x2f,
			0x00, 0x50, 0x00, 0x6c, 0x00, 0x61, 0x00, 0x79,
			0x00, 0x65, 0x00, 0x72, 0xff, 0xff, 0xff, 0xff,
		},
		CreateMessage: func() message { return new(stateEmitResponseMessage) },
		Message: &stateEmitResponseMessage{
			Name:     ClientPreferencesPlayer,
			Interval: 0xffffffff,
		},
	},
	{
		Name: "Beat info start stream",
		Bytes: []byte{
			0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00,
		},
		CreateMessage: func() message { return new(beatInfoStartStreamMessage) },
		Message:       &beatInfoStartStreamMessage{},
	},
	{
		Name: "Beat info stop stream",
		Bytes: []byte{
			0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01,
		},
		CreateMessage: func() message { return new(beatInfoStopStreamMessage) },
		Message:       &beatInfoStopStreamMessage{},
	},
	{
		Name: "Beat emit",
		Bytes: []byte{
			0x00, 0x00, 0x00, 0x90, 0x00, 0x00, 0x00, 0x02,
			0x00, 0x00, 0x06, 0x73, 0xfc, 0x64, 0x81, 0xac,
			0x00, 0x00, 0x00, 0x04, 0x40, 0x71, 0xd6, 0xa3,
			0x0e, 0xf9, 0xc6, 0x44, 0x40, 0x79, 0x0f, 0xe1,
			0xe8, 0x2d, 0x23, 0xbd, 0x40, 0x5b, 0x80, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x40, 0x71, 0xb9, 0x22,
			0x53, 0x6d, 0xc5, 0x20, 0x40, 0x80, 0x61, 0xb1,
			0x01, 0x76, 0x7d, 0xce, 0x40, 0x5a, 0x40, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x40, 0x5e, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x40, 0x5e, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x41, 0x5a, 0x30, 0x9c,
			0xe0, 0x16, 0xd3, 0x64, 0x41, 0x5b, 0x5b, 0x1c,
			0x8b, 0xd2, 0x15, 0xf2, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
		},
		CreateMessage: func() message { return new(beatEmitMessage) },
		Message: &beatEmitMessage{
			Clock: 7095225450924,
			Players: []PlayerInfo{
				{
					Beat:       285.41480920379786,
					TotalBeats: 400.99265306122453,
					Bpm:        110,
				},
				{
					Beat:       283.57088034514345,
					TotalBeats: 524.2114285714285,
					Bpm:        105,
				},
				{
					Beat:       0,
					TotalBeats: 0,
					Bpm:        120,
				},
				{
					Beat:       0,
					TotalBeats: 0,
					Bpm:        120,
				},
			},
			Timelines: []float64{6865523.501393173, 7171186.184697615, 0, 0},
		},
	},
	{
		Name: "EAAS discovery request",
		Bytes: []byte{
			0x45, 0x41, 0x41, 0x53, 0x01, 0x00,
		},
		CreateMessage: func() message { return new(eaasDiscoveryRequestMessage) },
		Message:       &eaasDiscoveryRequestMessage{},
	},
	{
		Name: "EAAS discovery response",
		Bytes: []byte{
			0x45, 0x41, 0x41, 0x53, 0x01, 0x01, 0x79, 0x9b,
			0x2d, 0xab, 0xf7, 0xc7, 0x43, 0x63, 0xb4, 0x9c,
			0x59, 0xe1, 0x91, 0x16, 0x89, 0x9e, 0x00, 0x00,
			0x00, 0x24, 0x00, 0x69, 0x00, 0x63, 0x00, 0x65,
			0x00, 0x64, 0x00, 0x72, 0x00, 0x65, 0x00, 0x61,
			0x00, 0x6d, 0x00, 0x2d, 0x00, 0x66, 0x00, 0x72,
			0x00, 0x61, 0x00, 0x6d, 0x00, 0x65, 0x00, 0x77,
			0x00, 0x6f, 0x00, 0x72, 0x00, 0x6b, 0x00, 0x00,
			0x00, 0x1c, 0x67, 0x72, 0x70, 0x63, 0x3a, 0x2f,
			0x2f, 0x31, 0x39, 0x32, 0x2e, 0x31, 0x36, 0x38,
			0x2e, 0x31, 0x38, 0x38, 0x2e, 0x31, 0x32, 0x30,
			0x3a, 0x35, 0x30, 0x30, 0x31, 0x30, 0x00, 0x00,
			0x00, 0x20, 0x00, 0x33, 0x00, 0x2e, 0x00, 0x34,
			0x00, 0x2e, 0x00, 0x30, 0x00, 0x2e, 0x00, 0x66,
			0x00, 0x36, 0x00, 0x62, 0x00, 0x33, 0x00, 0x64,
			0x00, 0x63, 0x00, 0x32, 0x00, 0x63, 0x00, 0x32,
			0x00, 0x30, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00,
			0x5f,
		},
		CreateMessage: func() message { return new(eaasDiscoveryResponseMessage) },
		Message: &eaasDiscoveryResponseMessage{
			tokenPrefixedMessage: tokenPrefixedMessage{
				Token: Token{0x79, 0x9b, 0x2d, 0xab, 0xf7, 0xc7, 0x43, 0x63, 0xb4, 0x9c, 0x59, 0xe1, 0x91, 0x16, 0x89, 0x9e},
			},
			Hostname:        "icedream-framework",
			SoftwareVersion: "3.4.0.f6b3dc2c20",
			URL:             "grpc://192.168.188.120:50010",
			Extra:           "_",
		},
	},
}

func Test_Messages_Read(t *testing.T) {
	for _, test := range testMessages {
		def := test
		t.Run(test.Name, func(t *testing.T) {
			r := bytes.NewReader(def.Bytes)
			m := def.CreateMessage()
			err := m.readFrom(r)
			require.NoError(t, err)
			require.Equal(t, def.Message, m)
		})
	}
}

func Test_Messages_Write(t *testing.T) {
	for _, test := range testMessages {
		def := test
		t.Run(test.Name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := def.Message.writeTo(buf)
			require.NoError(t, err)
			resultBytes := buf.Bytes()
			require.Equal(t, def.Bytes, resultBytes)
		})
	}
}

func Test_Messages_CheckMatch(t *testing.T) {
	for _, test := range testMessages {
		def := test
		t.Run(test.Name, func(t *testing.T) {
			ok, err := def.Message.checkMatch(bufio.NewReader(bytes.NewReader(def.Bytes)))
			require.NoError(t, err)
			require.True(t, ok)
		})
	}
}
