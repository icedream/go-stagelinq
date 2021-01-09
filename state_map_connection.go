package stagelinq

import (
	"encoding/json"
	"net"
)

// State represents a received state value.
type State struct {
	Name  string
	Value json.RawMessage
}

// StateMapConnection provides functionality to communicate with the StateMap data source.
type StateMapConnection struct {
	conn   *messageConnection
	errC   chan error
	stateC chan *State
}

var stateMapConnectionMessageSet = newDeviceConnMessageSet([]message{
	&stateEmitMessage{},
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

	// Before we do anything else, we announce our TCP source port in-protocol.
	// I have observed SoundSwitch and Resolume doing this, don't know what the purpose is though.
	msgConn.WriteMessage(&serviceAnnouncementMessage{
		tokenPrefixedMessage: tokenPrefixedMessage{
			Token: token,
		},
		Service: "StateMap",
		Port:    uint16(getPort(conn.LocalAddr())),
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
			var msg message
			msg, err = msgConn.ReadMessage()
			if err != nil {
				return
			}

			switch v := msg.(type) {
			case *stateEmitMessage:
				stateC <- &State{
					Name:  v.Name,
					Value: json.RawMessage(v.JSON),
				}
			}
		}
	}()

	smc = stateMapConn

	return
}

// Subscribe tells the StagelinQ device to let us know about changes for the given state value.
func (smc *StateMapConnection) Subscribe(event string) error {
	// TODO - check what to do with the int field in the state subscribe message, what is that?
	return smc.conn.WriteMessage(&stateSubscribeMessage{
		Name: event,
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
