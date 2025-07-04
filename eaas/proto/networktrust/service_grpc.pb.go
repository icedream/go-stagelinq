// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: proto/networktrust/service.proto

package networktrust

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	NetworkTrustService_CreateTrust_FullMethodName = "/networktrust.v1.NetworkTrustService/CreateTrust"
)

// NetworkTrustServiceClient is the client API for NetworkTrustService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NetworkTrustServiceClient interface {
	CreateTrust(ctx context.Context, in *CreateTrustRequest, opts ...grpc.CallOption) (*CreateTrustResponse, error)
}

type networkTrustServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNetworkTrustServiceClient(cc grpc.ClientConnInterface) NetworkTrustServiceClient {
	return &networkTrustServiceClient{cc}
}

func (c *networkTrustServiceClient) CreateTrust(ctx context.Context, in *CreateTrustRequest, opts ...grpc.CallOption) (*CreateTrustResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateTrustResponse)
	err := c.cc.Invoke(ctx, NetworkTrustService_CreateTrust_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NetworkTrustServiceServer is the server API for NetworkTrustService service.
// All implementations must embed UnimplementedNetworkTrustServiceServer
// for forward compatibility.
type NetworkTrustServiceServer interface {
	CreateTrust(context.Context, *CreateTrustRequest) (*CreateTrustResponse, error)
	mustEmbedUnimplementedNetworkTrustServiceServer()
}

// UnimplementedNetworkTrustServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNetworkTrustServiceServer struct{}

func (UnimplementedNetworkTrustServiceServer) CreateTrust(context.Context, *CreateTrustRequest) (*CreateTrustResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTrust not implemented")
}
func (UnimplementedNetworkTrustServiceServer) mustEmbedUnimplementedNetworkTrustServiceServer() {}
func (UnimplementedNetworkTrustServiceServer) testEmbeddedByValue()                             {}

// UnsafeNetworkTrustServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NetworkTrustServiceServer will
// result in compilation errors.
type UnsafeNetworkTrustServiceServer interface {
	mustEmbedUnimplementedNetworkTrustServiceServer()
}

func RegisterNetworkTrustServiceServer(s grpc.ServiceRegistrar, srv NetworkTrustServiceServer) {
	// If the following call pancis, it indicates UnimplementedNetworkTrustServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&NetworkTrustService_ServiceDesc, srv)
}

func _NetworkTrustService_CreateTrust_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTrustRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkTrustServiceServer).CreateTrust(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NetworkTrustService_CreateTrust_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkTrustServiceServer).CreateTrust(ctx, req.(*CreateTrustRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NetworkTrustService_ServiceDesc is the grpc.ServiceDesc for NetworkTrustService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NetworkTrustService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "networktrust.v1.NetworkTrustService",
	HandlerType: (*NetworkTrustServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateTrust",
			Handler:    _NetworkTrustService_CreateTrust_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/networktrust/service.proto",
}
