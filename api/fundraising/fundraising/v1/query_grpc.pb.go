// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: fundraising/fundraising/v1/query.proto

package fundraisingv1

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
	Query_Params_FullMethodName            = "/fundraising.fundraising.v1.Query/Params"
	Query_ListAuction_FullMethodName       = "/fundraising.fundraising.v1.Query/ListAuction"
	Query_GetAuction_FullMethodName        = "/fundraising.fundraising.v1.Query/GetAuction"
	Query_ListAllowedBidder_FullMethodName = "/fundraising.fundraising.v1.Query/ListAllowedBidder"
	Query_GetAllowedBidder_FullMethodName  = "/fundraising.fundraising.v1.Query/GetAllowedBidder"
	Query_ListBid_FullMethodName           = "/fundraising.fundraising.v1.Query/ListBid"
	Query_GetBid_FullMethodName            = "/fundraising.fundraising.v1.Query/GetBid"
	Query_ListVestingQueue_FullMethodName  = "/fundraising.fundraising.v1.Query/ListVestingQueue"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// Queries a list of Auction items.
	ListAuction(ctx context.Context, in *QueryAllAuctionRequest, opts ...grpc.CallOption) (*QueryAllAuctionResponse, error)
	GetAuction(ctx context.Context, in *QueryGetAuctionRequest, opts ...grpc.CallOption) (*QueryGetAuctionResponse, error)
	// Queries a list of AllowedBidder items.
	ListAllowedBidder(ctx context.Context, in *QueryAllAllowedBidderRequest, opts ...grpc.CallOption) (*QueryAllAllowedBidderResponse, error)
	GetAllowedBidder(ctx context.Context, in *QueryGetAllowedBidderRequest, opts ...grpc.CallOption) (*QueryGetAllowedBidderResponse, error)
	// Queries a list of Bid items.
	ListBid(ctx context.Context, in *QueryAllBidRequest, opts ...grpc.CallOption) (*QueryAllBidResponse, error)
	GetBid(ctx context.Context, in *QueryGetBidRequest, opts ...grpc.CallOption) (*QueryGetBidResponse, error)
	// Queries a list of VestingQueue items.
	ListVestingQueue(ctx context.Context, in *QueryAllVestingQueueRequest, opts ...grpc.CallOption) (*QueryAllVestingQueueResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, Query_Params_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ListAuction(ctx context.Context, in *QueryAllAuctionRequest, opts ...grpc.CallOption) (*QueryAllAuctionResponse, error) {
	out := new(QueryAllAuctionResponse)
	err := c.cc.Invoke(ctx, Query_ListAuction_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetAuction(ctx context.Context, in *QueryGetAuctionRequest, opts ...grpc.CallOption) (*QueryGetAuctionResponse, error) {
	out := new(QueryGetAuctionResponse)
	err := c.cc.Invoke(ctx, Query_GetAuction_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ListAllowedBidder(ctx context.Context, in *QueryAllAllowedBidderRequest, opts ...grpc.CallOption) (*QueryAllAllowedBidderResponse, error) {
	out := new(QueryAllAllowedBidderResponse)
	err := c.cc.Invoke(ctx, Query_ListAllowedBidder_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetAllowedBidder(ctx context.Context, in *QueryGetAllowedBidderRequest, opts ...grpc.CallOption) (*QueryGetAllowedBidderResponse, error) {
	out := new(QueryGetAllowedBidderResponse)
	err := c.cc.Invoke(ctx, Query_GetAllowedBidder_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ListBid(ctx context.Context, in *QueryAllBidRequest, opts ...grpc.CallOption) (*QueryAllBidResponse, error) {
	out := new(QueryAllBidResponse)
	err := c.cc.Invoke(ctx, Query_ListBid_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetBid(ctx context.Context, in *QueryGetBidRequest, opts ...grpc.CallOption) (*QueryGetBidResponse, error) {
	out := new(QueryGetBidResponse)
	err := c.cc.Invoke(ctx, Query_GetBid_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ListVestingQueue(ctx context.Context, in *QueryAllVestingQueueRequest, opts ...grpc.CallOption) (*QueryAllVestingQueueResponse, error) {
	out := new(QueryAllVestingQueueResponse)
	err := c.cc.Invoke(ctx, Query_ListVestingQueue_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// Queries a list of Auction items.
	ListAuction(context.Context, *QueryAllAuctionRequest) (*QueryAllAuctionResponse, error)
	GetAuction(context.Context, *QueryGetAuctionRequest) (*QueryGetAuctionResponse, error)
	// Queries a list of AllowedBidder items.
	ListAllowedBidder(context.Context, *QueryAllAllowedBidderRequest) (*QueryAllAllowedBidderResponse, error)
	GetAllowedBidder(context.Context, *QueryGetAllowedBidderRequest) (*QueryGetAllowedBidderResponse, error)
	// Queries a list of Bid items.
	ListBid(context.Context, *QueryAllBidRequest) (*QueryAllBidResponse, error)
	GetBid(context.Context, *QueryGetBidRequest) (*QueryGetBidResponse, error)
	// Queries a list of VestingQueue items.
	ListVestingQueue(context.Context, *QueryAllVestingQueueRequest) (*QueryAllVestingQueueResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) ListAuction(context.Context, *QueryAllAuctionRequest) (*QueryAllAuctionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAuction not implemented")
}
func (UnimplementedQueryServer) GetAuction(context.Context, *QueryGetAuctionRequest) (*QueryGetAuctionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuction not implemented")
}
func (UnimplementedQueryServer) ListAllowedBidder(context.Context, *QueryAllAllowedBidderRequest) (*QueryAllAllowedBidderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAllowedBidder not implemented")
}
func (UnimplementedQueryServer) GetAllowedBidder(context.Context, *QueryGetAllowedBidderRequest) (*QueryGetAllowedBidderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllowedBidder not implemented")
}
func (UnimplementedQueryServer) ListBid(context.Context, *QueryAllBidRequest) (*QueryAllBidResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBid not implemented")
}
func (UnimplementedQueryServer) GetBid(context.Context, *QueryGetBidRequest) (*QueryGetBidResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBid not implemented")
}
func (UnimplementedQueryServer) ListVestingQueue(context.Context, *QueryAllVestingQueueRequest) (*QueryAllVestingQueueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListVestingQueue not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Params_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ListAuction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllAuctionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ListAuction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ListAuction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ListAuction(ctx, req.(*QueryAllAuctionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetAuction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetAuctionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetAuction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetAuction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetAuction(ctx, req.(*QueryGetAuctionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ListAllowedBidder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllAllowedBidderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ListAllowedBidder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ListAllowedBidder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ListAllowedBidder(ctx, req.(*QueryAllAllowedBidderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetAllowedBidder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetAllowedBidderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetAllowedBidder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetAllowedBidder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetAllowedBidder(ctx, req.(*QueryGetAllowedBidderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ListBid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllBidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ListBid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ListBid_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ListBid(ctx, req.(*QueryAllBidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetBid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetBidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetBid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetBid_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetBid(ctx, req.(*QueryGetBidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ListVestingQueue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllVestingQueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ListVestingQueue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ListVestingQueue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ListVestingQueue(ctx, req.(*QueryAllVestingQueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fundraising.fundraising.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "ListAuction",
			Handler:    _Query_ListAuction_Handler,
		},
		{
			MethodName: "GetAuction",
			Handler:    _Query_GetAuction_Handler,
		},
		{
			MethodName: "ListAllowedBidder",
			Handler:    _Query_ListAllowedBidder_Handler,
		},
		{
			MethodName: "GetAllowedBidder",
			Handler:    _Query_GetAllowedBidder_Handler,
		},
		{
			MethodName: "ListBid",
			Handler:    _Query_ListBid_Handler,
		},
		{
			MethodName: "GetBid",
			Handler:    _Query_GetBid_Handler,
		},
		{
			MethodName: "ListVestingQueue",
			Handler:    _Query_ListVestingQueue_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fundraising/fundraising/v1/query.proto",
}
