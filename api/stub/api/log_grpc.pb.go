// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LogServiceClient is the client API for LogService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LogServiceClient interface {
	GetLogList(ctx context.Context, in *LogPaginationRequest, opts ...grpc.CallOption) (*GetLogListResult, error)
}

type logServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLogServiceClient(cc grpc.ClientConnInterface) LogServiceClient {
	return &logServiceClient{cc}
}

func (c *logServiceClient) GetLogList(ctx context.Context, in *LogPaginationRequest, opts ...grpc.CallOption) (*GetLogListResult, error) {
	out := new(GetLogListResult)
	err := c.cc.Invoke(ctx, "/api.LogService/GetLogList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogServiceServer is the server API for LogService service.
// All implementations should embed UnimplementedLogServiceServer
// for forward compatibility
type LogServiceServer interface {
	GetLogList(context.Context, *LogPaginationRequest) (*GetLogListResult, error)
}

// UnimplementedLogServiceServer should be embedded to have forward compatible implementations.
type UnimplementedLogServiceServer struct {
}

func (UnimplementedLogServiceServer) GetLogList(context.Context, *LogPaginationRequest) (*GetLogListResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLogList not implemented")
}

// UnsafeLogServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LogServiceServer will
// result in compilation errors.
type UnsafeLogServiceServer interface {
	mustEmbedUnimplementedLogServiceServer()
}

func RegisterLogServiceServer(s grpc.ServiceRegistrar, srv LogServiceServer) {
	s.RegisterService(&LogService_ServiceDesc, srv)
}

func _LogService_GetLogList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogPaginationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogServiceServer).GetLogList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.LogService/GetLogList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogServiceServer).GetLogList(ctx, req.(*LogPaginationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LogService_ServiceDesc is the grpc.ServiceDesc for LogService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LogService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.LogService",
	HandlerType: (*LogServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLogList",
			Handler:    _LogService_GetLogList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "log.proto",
}