// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: otp_authenticator.proto

package proto

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

const (
	OtpAuthenticator_GeneratePrivateKey_FullMethodName = "/otp_authenticator.OtpAuthenticator/GeneratePrivateKey"
	OtpAuthenticator_GenerateOtp_FullMethodName        = "/otp_authenticator.OtpAuthenticator/GenerateOtp"
)

// OtpAuthenticatorClient is the client API for OtpAuthenticator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OtpAuthenticatorClient interface {
	GeneratePrivateKey(ctx context.Context, in *GeneratePrivateKeyRequest, opts ...grpc.CallOption) (*GeneratePrivateKeyResponse, error)
	GenerateOtp(ctx context.Context, in *GenerateOtpRequest, opts ...grpc.CallOption) (*GenerateOtpRequestResponse, error)
}

type otpAuthenticatorClient struct {
	cc grpc.ClientConnInterface
}

func NewOtpAuthenticatorClient(cc grpc.ClientConnInterface) OtpAuthenticatorClient {
	return &otpAuthenticatorClient{cc}
}

func (c *otpAuthenticatorClient) GeneratePrivateKey(ctx context.Context, in *GeneratePrivateKeyRequest, opts ...grpc.CallOption) (*GeneratePrivateKeyResponse, error) {
	out := new(GeneratePrivateKeyResponse)
	err := c.cc.Invoke(ctx, OtpAuthenticator_GeneratePrivateKey_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *otpAuthenticatorClient) GenerateOtp(ctx context.Context, in *GenerateOtpRequest, opts ...grpc.CallOption) (*GenerateOtpRequestResponse, error) {
	out := new(GenerateOtpRequestResponse)
	err := c.cc.Invoke(ctx, OtpAuthenticator_GenerateOtp_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OtpAuthenticatorServer is the server API for OtpAuthenticator service.
// All implementations must embed UnimplementedOtpAuthenticatorServer
// for forward compatibility
type OtpAuthenticatorServer interface {
	GeneratePrivateKey(context.Context, *GeneratePrivateKeyRequest) (*GeneratePrivateKeyResponse, error)
	GenerateOtp(context.Context, *GenerateOtpRequest) (*GenerateOtpRequestResponse, error)
	mustEmbedUnimplementedOtpAuthenticatorServer()
}

// UnimplementedOtpAuthenticatorServer must be embedded to have forward compatible implementations.
type UnimplementedOtpAuthenticatorServer struct {
}

func (UnimplementedOtpAuthenticatorServer) GeneratePrivateKey(context.Context, *GeneratePrivateKeyRequest) (*GeneratePrivateKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GeneratePrivateKey not implemented")
}
func (UnimplementedOtpAuthenticatorServer) GenerateOtp(context.Context, *GenerateOtpRequest) (*GenerateOtpRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateOtp not implemented")
}
func (UnimplementedOtpAuthenticatorServer) mustEmbedUnimplementedOtpAuthenticatorServer() {}

// UnsafeOtpAuthenticatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OtpAuthenticatorServer will
// result in compilation errors.
type UnsafeOtpAuthenticatorServer interface {
	mustEmbedUnimplementedOtpAuthenticatorServer()
}

func RegisterOtpAuthenticatorServer(s grpc.ServiceRegistrar, srv OtpAuthenticatorServer) {
	s.RegisterService(&OtpAuthenticator_ServiceDesc, srv)
}

func _OtpAuthenticator_GeneratePrivateKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GeneratePrivateKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OtpAuthenticatorServer).GeneratePrivateKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OtpAuthenticator_GeneratePrivateKey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OtpAuthenticatorServer).GeneratePrivateKey(ctx, req.(*GeneratePrivateKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OtpAuthenticator_GenerateOtp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateOtpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OtpAuthenticatorServer).GenerateOtp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OtpAuthenticator_GenerateOtp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OtpAuthenticatorServer).GenerateOtp(ctx, req.(*GenerateOtpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OtpAuthenticator_ServiceDesc is the grpc.ServiceDesc for OtpAuthenticator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OtpAuthenticator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "otp_authenticator.OtpAuthenticator",
	HandlerType: (*OtpAuthenticatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GeneratePrivateKey",
			Handler:    _OtpAuthenticator_GeneratePrivateKey_Handler,
		},
		{
			MethodName: "GenerateOtp",
			Handler:    _OtpAuthenticator_GenerateOtp_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "otp_authenticator.proto",
}