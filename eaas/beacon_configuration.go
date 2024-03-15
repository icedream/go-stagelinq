package eaas

import "context"

// The default EAAS gRPC API port.
const DefaultEAASGRPCPort uint16 = 50010

// The default EAAS HTTP server port.
const DefaultEAASHTTPPort uint16 = DefaultEAASGRPCPort + 10

// BeaconConfiguration contains configurable values for setting up a EAAS
// announcer.
type BeaconConfiguration struct {
	// Context can be set to allow cancellation of network operations from
	// somewhere else in the code.
	Context context.Context

	// Name is the name under which we announce ourselves to the network. Denon
	// software tends to use hostname here.
	Name string

	// SoftwareVersion is your application's version. It is used for StagelinQ
	// announcements to the network.
	SoftwareVersion string

	// Token is used as part of announcements and main data communication. It is
	// currently recommended to leave this empty.
	Token Token

	// The host to report back to clients with where the EAAS gRPC API is
	// listening on.
	//
	// If left empty, defaults to the IP the beacon is bound to. It is
	// recommended to leave this empty.
	GRPCHost string

	// The port to report back to clients which the EAAS gRPC API is listening on.
	//
	// If left zero, defaults to the default EAAS gRPC API port (50010).
	GRPCPort uint16
}
