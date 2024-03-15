package eaas

import (
	"context"

	"github.com/icedream/go-stagelinq/eaas/proto/enginelibrary"
	"github.com/icedream/go-stagelinq/eaas/proto/networktrust"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EngineLibraryConnection struct {
	grpc *grpc.ClientConn

	networktrust.NetworkTrustServiceClient
	enginelibrary.EngineLibraryServiceClient
}

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

func Dial(url string) (*EngineLibraryConnection, error) {
	return DialContext(context.Background(), url)
}

func DialContext(ctx context.Context, url string) (*EngineLibraryConnection, error) {
	clientConn, err := grpc.DialContext(ctx, url,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return newClient(clientConn), nil
}
