package eaas

import (
	"context"
	"time"
)

// DiscovererConfiguration contains configurable values for setting up an EAAS
// discoverer.
type DiscovererConfiguration struct {
	// Context can be set to allow cancellation of network operations from
	// somewhere else in the code.
	Context context.Context

	// DiscoveryTimeout is the duration for which Listener.Discover will wait
	// for EAAS devices to announce themselves. If this is not set, no timeout
	// will occur.
	DiscoveryTimeout time.Duration
}
