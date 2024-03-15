package eaas

import (
	"context"
)

// AnnouncerConfiguration contains configurable values for setting up a EAAS
// announcer.
type AnnouncerConfiguration struct {
	// Context can be set to allow cancellation of network operations from somewhere else in the code.
	Context context.Context

	// Name is the name under which we announce ourselves to the network. Denon
	// software tends to use hostname here.
	Name string

	// SoftwareVersion is your application's version. It is used for StagelinQ announcements to the network.
	SoftwareVersion string

	// Token is used as part of announcements and main data communication. It is currently recommended to leave this empty.
	Token Token
}
