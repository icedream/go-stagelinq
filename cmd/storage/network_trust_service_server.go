package main

import (
	"context"
	"log"

	"github.com/icedream/go-stagelinq/eaas/proto/networktrust"
)

var _ networktrust.NetworkTrustServiceServer = &NetworkTrustServiceServer{}

// NetworkTrustServiceServer is an example network trust service server
// implementation. All it does is approve everything that asks.
type NetworkTrustServiceServer struct {
	networktrust.UnimplementedNetworkTrustServiceServer
}

// CreateTrust implements networktrust.NetworkTrustServiceServer.
func (n *NetworkTrustServiceServer) CreateTrust(ctx context.Context, req *networktrust.CreateTrustRequest) (*networktrust.CreateTrustResponse, error) {
	log.Printf("CreateTrust: %+v", req)

	// Just allow all for now
	return &networktrust.CreateTrustResponse{
		Response: &networktrust.CreateTrustResponse_Granted{},
	}, nil
}
