package stagelinq

import (
	"context"
	"time"
)

// ListenerConfiguration contains configurable values for setting up a StagelinQ listener.
type ListenerConfiguration struct {
	// Context can be set to allow cancellation of network operations from somewhere else in the code.
	Context context.Context

	// DiscoveryTimeout is the duration for which Listener.Discover will wait for StagelinQ devices to announce themselves.
	// If this is not set, no timeout will occur.
	DiscoveryTimeout time.Duration

	// Name is the name under which we announce ourselves to the network.
	// For example, Resolute uses the computer user name here, and Denon devices use their identifying abbreviation (the Prime 4 uses "prime4").
	Name string

	// SoftwareName is your application's name. It is used for StagelinQ announcements to the network.
	SoftwareName string

	// SoftwareVersion is your application's version. It is used for StagelinQ announcements to the network.
	SoftwareVersion string

	// Token is used as part of announcements and main data communication. It is currently recommended to leave this empty.
	Token Token
}
