package stagelinq

import (
	"net"
	"sync"
	"time"
)

// Service contains information about a data service a device provides.
type Service struct {
	Name string
	Port uint16
}

// MainConnection represents a connection to the main TCP port of a StagelinQ device.
type MainConnection struct {
	lock sync.Mutex

	token   Token
	msgConn *messageConnection

	offeredServices []*Service

	servicesC                 chan *Service
	atLeastOneServiceReceived bool

	errorC chan error

	reference int64
}

var mainConnectionMessageSet = newDeviceConnMessageSet([]message{
	&serviceAnnouncementMessage{},
	&referenceMessage{},
	&servicesRequestMessage{},
})

// newMainConnection wraps an existing network connection to communicate StagelinQ main connection messages with it.
func newMainConnection(conn net.Conn, token Token, targetToken Token, offeredServices []*Service) (retval *MainConnection, err error) {
	msgConn := newMessageConnection(conn, mainConnectionMessageSet)

	mainConn := &MainConnection{
		token:           token,
		msgConn:         msgConn,
		errorC:          make(chan error, 1),
		offeredServices: offeredServices,
	}

	go func() {
		for {
			<-time.After(250 * time.Millisecond)

			// TODO - we're always returning zero as timestamp here just like SoundSwitch does, we still need to implement the behavior Resolume Arena
			mainConn.lock.Lock()
			ref := mainConn.reference
			mainConn.lock.Unlock()
			if err = mainConn.msgConn.WriteMessage(&referenceMessage{
				tokenPrefixedMessage: tokenPrefixedMessage{
					Token: mainConn.token,
				},
				Token2:    targetToken,
				Reference: ref,
			}); err != nil {
				return
			}
		}
	}()

	go func() {
		var err error
		defer func() {
			if err != nil {
				mainConn.errorC <- err
				close(mainConn.errorC)
			}
			if mainConn.servicesC != nil {
				close(mainConn.servicesC)
				mainConn.servicesC = nil
			}
		}()
		for {
			var msg message
			msg, err = mainConn.msgConn.ReadMessage()
			if err != nil {
				return
			}

			func() {
				mainConn.lock.Lock()
				defer mainConn.lock.Unlock()

				switch v := msg.(type) {
				case *serviceAnnouncementMessage:
					if mainConn.servicesC == nil {
						err = ErrInvalidMessageReceived
						break
					}
					mainConn.servicesC <- &Service{
						Name: v.Service,
						Port: v.Port,
					}
				case *referenceMessage:
					if mainConn.servicesC != nil {
						close(mainConn.servicesC)
						mainConn.servicesC = nil
					}
					// TODO - not sure what else to actually do with this information yet
					// mainConn.reference = v.Reference
				case *servicesRequestMessage:
					for _, service := range mainConn.offeredServices {
						if err = mainConn.announceService(service.Name, service.Port); err != nil {
							return
						}
					}
				}
			}()
		}
	}()

	retval = mainConn

	return
}

// Close terminates the connection.
func (conn *MainConnection) Close() (err error) {
	return conn.msgConn.conn.Close()
}

// RequestServices asks the device to return other TCP ports it is listening on and which services it provides on them.
func (conn *MainConnection) RequestServices() (retval []*Service, err error) {
	if err = conn.requestServices(); err != nil {
		return
	}

	conn.lock.Lock()
	serviceC := make(chan *Service)
	conn.servicesC = serviceC
	conn.atLeastOneServiceReceived = false
	services := []*Service{}
	conn.lock.Unlock()

	for service := range serviceC {
		services = append(services, service)
	}
	select {
	case err = <-conn.errorC:
	default:
	}

	// the message reading loop already will set conn.servicesC back to nil

	retval = services

	return
}

func (conn *MainConnection) requestServices() (err error) {
	if err = conn.msgConn.WriteMessage(&servicesRequestMessage{
		tokenPrefixedMessage: tokenPrefixedMessage{
			Token: conn.token,
		},
	}); err != nil {
		return
	}
	return
}

// announceService tells the device about a service we provide.
func (conn *MainConnection) announceService(name string, port uint16) error {
	return conn.msgConn.WriteMessage(&serviceAnnouncementMessage{
		tokenPrefixedMessage: tokenPrefixedMessage{
			Token: conn.token,
		},
		Service: name,
		Port:    port,
	})
}
