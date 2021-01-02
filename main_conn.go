package stagelinq

import "net"

type AnnouncedService struct {
	Name string
	Port uint16
}

type Reference []byte

type MainConnection struct {
	token [16]byte
	*deviceConn

	servicesC  chan *AnnouncedService
	referenceC chan *Reference
}

// ConnectToDevice starts a StagelinQ connection with the given device.
func ConnectToDevice(listener *Listener, device *Device) (conn *MainConnection, err error) {
	ip := device.IP
	port := device.port

	tcpConn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   ip,
		Port: int(port),
	})
	if err != nil {
		return
	}

	deviceConn := newDeviceConn(tcpConn)

	conn = &MainConnection{
		token:      listener.token,
		deviceConn: deviceConn,
	}

	return
}

func (conn *MainConnection) Close() (err error) {
	err = conn.deviceConn.Close()
	return
}

// RequestServices asks the device to return other TCP ports it is listening on and which services it provides on them.
func (conn *MainConnection) RequestServices() error {
	return conn.deviceConn.WriteMessage(&ServicesRequestMessage{
		tokenPrefixedMessage: tokenPrefixedMessage{
			Token: conn.token,
		},
	})
}

// AnnounceService tells the device about a service we provide.
func (conn *MainConnection) AnnounceService(name string, port uint16) error {
	return conn.deviceConn.WriteMessage(&ServiceAnnouncementMessage{
		tokenPrefixedMessage: tokenPrefixedMessage{
			Token: conn.token,
		},
		Service: name,
		Port:    port,
	})
}
