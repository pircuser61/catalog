// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: api/api.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CatalogClient is the client API for Catalog service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CatalogClient interface {
	GoodCreate(ctx context.Context, in *GoodCreateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GoodGet(ctx context.Context, in *GoodGetRequest, opts ...grpc.CallOption) (*GoodGetResponse, error)
	GoodList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GoodListResponse, error)
	GoodUpdate(ctx context.Context, in *GoodUpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GoodDelete(ctx context.Context, in *GoodDeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type catalogClient struct {
	cc grpc.ClientConnInterface
}

func NewCatalogClient(cc grpc.ClientConnInterface) CatalogClient {
	return &catalogClient{cc}
}

func (c *catalogClient) GoodCreate(ctx context.Context, in *GoodCreateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/catalog.api.catalog/GoodCreate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogClient) GoodGet(ctx context.Context, in *GoodGetRequest, opts ...grpc.CallOption) (*GoodGetResponse, error) {
	out := new(GoodGetResponse)
	err := c.cc.Invoke(ctx, "/catalog.api.catalog/GoodGet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogClient) GoodList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GoodListResponse, error) {
	out := new(GoodListResponse)
	err := c.cc.Invoke(ctx, "/catalog.api.catalog/GoodList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogClient) GoodUpdate(ctx context.Context, in *GoodUpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/catalog.api.catalog/GoodUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogClient) GoodDelete(ctx context.Context, in *GoodDeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/catalog.api.catalog/GoodDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CatalogServer is the server API for Catalog service.
// All implementations must embed UnimplementedCatalogServer
// for forward compatibility
type CatalogServer interface {
	GoodCreate(context.Context, *GoodCreateRequest) (*emptypb.Empty, error)
	GoodGet(context.Context, *GoodGetRequest) (*GoodGetResponse, error)
	GoodList(context.Context, *emptypb.Empty) (*GoodListResponse, error)
	GoodUpdate(context.Context, *GoodUpdateRequest) (*emptypb.Empty, error)
	GoodDelete(context.Context, *GoodDeleteRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedCatalogServer()
}

// UnimplementedCatalogServer must be embedded to have forward compatible implementations.
type UnimplementedCatalogServer struct {
}

func (UnimplementedCatalogServer) GoodCreate(context.Context, *GoodCreateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GoodCreate not implemented")
}
func (UnimplementedCatalogServer) GoodGet(context.Context, *GoodGetRequest) (*GoodGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GoodGet not implemented")
}
func (UnimplementedCatalogServer) GoodList(context.Context, *emptypb.Empty) (*GoodListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GoodList not implemented")
}
func (UnimplementedCatalogServer) GoodUpdate(context.Context, *GoodUpdateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GoodUpdate not implemented")
}
func (UnimplementedCatalogServer) GoodDelete(context.Context, *GoodDeleteRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GoodDelete not implemented")
}
func (UnimplementedCatalogServer) mustEmbedUnimplementedCatalogServer() {}

// UnsafeCatalogServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CatalogServer will
// result in compilation errors.
type UnsafeCatalogServer interface {
	mustEmbedUnimplementedCatalogServer()
}

func RegisterCatalogServer(s grpc.ServiceRegistrar, srv CatalogServer) {
	s.RegisterService(&Catalog_ServiceDesc, srv)
}

func _Catalog_GoodCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GoodCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogServer).GoodCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/catalog.api.catalog/GoodCreate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogServer).GoodCreate(ctx, req.(*GoodCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Catalog_GoodGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GoodGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogServer).GoodGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/catalog.api.catalog/GoodGet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogServer).GoodGet(ctx, req.(*GoodGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Catalog_GoodList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogServer).GoodList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/catalog.api.catalog/GoodList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogServer).GoodList(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Catalog_GoodUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GoodUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogServer).GoodUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/catalog.api.catalog/GoodUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogServer).GoodUpdate(ctx, req.(*GoodUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Catalog_GoodDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GoodDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogServer).GoodDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/catalog.api.catalog/GoodDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogServer).GoodDelete(ctx, req.(*GoodDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Catalog_ServiceDesc is the grpc.ServiceDesc for Catalog service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Catalog_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "catalog.api.catalog",
	HandlerType: (*CatalogServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GoodCreate",
			Handler:    _Catalog_GoodCreate_Handler,
		},
		{
			MethodName: "GoodGet",
			Handler:    _Catalog_GoodGet_Handler,
		},
		{
			MethodName: "GoodList",
			Handler:    _Catalog_GoodList_Handler,
		},
		{
			MethodName: "GoodUpdate",
			Handler:    _Catalog_GoodUpdate_Handler,
		},
		{
			MethodName: "GoodDelete",
			Handler:    _Catalog_GoodDelete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/api.proto",
}
