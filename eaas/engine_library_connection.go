package eaas

import (
	"context"

	"github.com/icedream/go-stagelinq/eaas/proto/enginelibrary"
	"github.com/icedream/go-stagelinq/eaas/proto/networktrust"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// EngineLibraryConnection provides API functionality for interacting with
// remote Engine libraries.
type EngineLibraryConnection struct {
	grpc *grpc.ClientConn

	networktrust.NetworkTrustServiceClient
	enginelibrary.EngineLibraryServiceClient
}

// Close shuts down the gRPC connection to the EAAS API.
//
// Invoke this method once you're done using the connection to free up
// resources.
func (c *EngineLibraryConnection) Close() error {
	return c.grpc.Close()
}

func newClient(cc *grpc.ClientConn) *EngineLibraryConnection {
	return &EngineLibraryConnection{
		grpc:                       cc,
		NetworkTrustServiceClient:  networktrust.NewNetworkTrustServiceClient(cc),
		EngineLibraryServiceClient: enginelibrary.NewEngineLibraryServiceClient(cc),
	}
}

// Dial connects to the given EAAS gRPC API endpoint.
//
// Currently, the url value must be compatible with targets supported by the
// go-grpc library.
func Dial(url string) (*EngineLibraryConnection, error) {
	return DialContext(context.Background(), url)
}

// DialContext connects to the given EAAS gRPC API endpoint and sets up the
// backing gRPC client with the given context.
//
// Currently, the url value must be compatible with targets supported by the
// go-grpc library.
func DialContext(
	ctx context.Context,
	url string,
) (*EngineLibraryConnection, error) {
	// TODO - handle grpc://<IP>:<PORT> as that's the format Engine uses
	clientConn, err := grpc.DialContext(ctx, url,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return newClient(clientConn), nil
}
