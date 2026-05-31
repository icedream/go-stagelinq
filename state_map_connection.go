package stagelinq

import (
	"encoding/json"
	"net"
	"strings"
	"time"

	"github.com/icedream/go-stagelinq/internal/messages"
	"github.com/icedream/go-stagelinq/internal/socket"
)

// State represents a received state value.
type State struct {
	Name  string
	Value map[string]interface{}
}

// StateMapConnection provides functionality to communicate with the StateMap data source.
type StateMapConnection struct {
	conn   *messageConnection
	errC   chan error
	stateC chan *State
}

var stateMapConnectionMessageSet = newDeviceConnMessageSet([]messages.Message{
	&stateEmitMessage{},
	&stateEmitResponseMessage{},
})

// NewStateMapConnection wraps an existing network connection and returns a StateMapConnection, providing the functionality to subscribe to and receive changes of state values.
// You need to pass the token that you have announced for your own device on the network.
func NewStateMapConnection(conn net.Conn, token Token) (smc *StateMapConnection, err error) {
	msgConn := newMessageConnection(conn, stateMapConnectionMessageSet)

	errC := make(chan error, 1)
	stateC := make(chan *State, 1)

	stateMapConn := &StateMapConnection{
		conn:   msgConn,
		errC:   errC,
		stateC: stateC,
	}

	// Announce our TCP source port to the device before subscribing. This
	// registers the port the device should use to push state updates back
	// to us (the callback/return channel).
	msgConn.WriteMessage(&serviceAnnouncementMessage{
		TokenPrefixedMessage: messages.TokenPrefixedMessage{
			Token: messages.Token(token),
		},
		Service: "StateMap",
		Port:    socket.GetPort(conn.LocalAddr()),
	})

	go func() {
		var err error
		defer func() {
			if err != nil {
				stateMapConn.errC <- err
				close(stateMapConn.errC)
			}
			close(stateMapConn.stateC)
		}()
		for {
			var msg messages.Message
			msg, err = msgConn.ReadMessage()
			if err != nil {
				return
			}

			switch v := msg.(type) {
			case *stateEmitMessage:
				state := &State{
					Name: v.Name,
				}
				err = json.NewDecoder(strings.NewReader(v.JSON)).Decode(&state.Value)
				if err != nil {
					return
				}
				stateC <- state
			}
		}
	}()

	smc = stateMapConn

	return
}

// StateMapSubscriptionOption represents an option for state map subscriptions.
type StateMapSubscriptionOption func(*stateSubscribeMessage)

// WithInterval sets a minimum time interval at which the event value should be
// resent on the wire. The default value of 0 means values are emitted only when
// changed as an event.
func WithInterval(d time.Duration) StateMapSubscriptionOption {
	return func(m *stateSubscribeMessage) {
		m.Interval = uint32(d.Milliseconds())
	}
}

// Subscribe tells the StagelinQ device to send us updates for the given state
// value path.
func (smc *StateMapConnection) Subscribe(event string, opts ...StateMapSubscriptionOption) error {
	m := &stateSubscribeMessage{
		Name: event,
	}
	for _, o := range opts {
		o(m)
	}
	return smc.conn.WriteMessage(m)
}

func (smc *StateMapConnection) Emit(state *State) error {
	jsonBytes, err := json.Marshal(state.Value)
	if err != nil {
		return err
	}
	return smc.conn.WriteMessage(&stateEmitMessage{
		Name: state.Name,
		JSON: string(jsonBytes),
	})
}

// StateC returns the channel via which state changes will be returned for this connection.
func (smc *StateMapConnection) StateC() <-chan *State {
	return smc.stateC
}

// ErrorC returns the channel via which connectionrerors will be returned for this connection.
func (smc *StateMapConnection) ErrorC() <-chan error {
	return smc.errC
}
