package main

import (
	"context"
	"time"

	"github.com/icedream/go-stagelinq/eaas/proto/networktrust"
)

var _ networktrust.NetworkTrustServiceServer = &NetworkTrustServer{}

type NetworkTrustServer struct {
	networktrust.UnimplementedNetworkTrustServiceServer
}

// CreateTrust implements networktrust.NetworkTrustServiceServer.
func (n *NetworkTrustServer) CreateTrust(ctx context.Context, _ *networktrust.CreateTrustRequest) (*networktrust.CreateTrustResponse, error) {
	// safety sleep to not confuse Engine
	time.After(time.Second)

	// Just allow all for now
	return &networktrust.CreateTrustResponse{
		Response: &networktrust.CreateTrustResponse_Granted{},
	}, nil
}
