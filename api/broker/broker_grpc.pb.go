// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.0--rc2
// source: api/broker/broker.proto

package broker

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
	Broker_Push_FullMethodName          = "/api.broker.Broker/Push"
	Broker_PushMsg_FullMethodName       = "/api.broker.Broker/PushMsg"
	Broker_Broadcast_FullMethodName     = "/api.broker.Broker/Broadcast"
	Broker_BroadcastRoom_FullMethodName = "/api.broker.Broker/BroadcastRoom"
	Broker_Rooms_FullMethodName         = "/api.broker.Broker/Rooms"
)

// BrokerClient is the client API for Broker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BrokerClient interface {
	Push(ctx context.Context, in *PushReq, opts ...grpc.CallOption) (*PushReply, error)
	PushMsg(ctx context.Context, in *PushMsgReq, opts ...grpc.CallOption) (*PushMsgReply, error)
	Broadcast(ctx context.Context, in *BroadcastReq, opts ...grpc.CallOption) (*BroadcastReply, error)
	BroadcastRoom(ctx context.Context, in *BroadcastRoomReq, opts ...grpc.CallOption) (*BroadcastRoomReply, error)
	Rooms(ctx context.Context, in *RoomsReq, opts ...grpc.CallOption) (*RoomsReply, error)
}

type brokerClient struct {
	cc grpc.ClientConnInterface
}

func NewBrokerClient(cc grpc.ClientConnInterface) BrokerClient {
	return &brokerClient{cc}
}

func (c *brokerClient) Push(ctx context.Context, in *PushReq, opts ...grpc.CallOption) (*PushReply, error) {
	out := new(PushReply)
	err := c.cc.Invoke(ctx, Broker_Push_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) PushMsg(ctx context.Context, in *PushMsgReq, opts ...grpc.CallOption) (*PushMsgReply, error) {
	out := new(PushMsgReply)
	err := c.cc.Invoke(ctx, Broker_PushMsg_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) Broadcast(ctx context.Context, in *BroadcastReq, opts ...grpc.CallOption) (*BroadcastReply, error) {
	out := new(BroadcastReply)
	err := c.cc.Invoke(ctx, Broker_Broadcast_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) BroadcastRoom(ctx context.Context, in *BroadcastRoomReq, opts ...grpc.CallOption) (*BroadcastRoomReply, error) {
	out := new(BroadcastRoomReply)
	err := c.cc.Invoke(ctx, Broker_BroadcastRoom_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) Rooms(ctx context.Context, in *RoomsReq, opts ...grpc.CallOption) (*RoomsReply, error) {
	out := new(RoomsReply)
	err := c.cc.Invoke(ctx, Broker_Rooms_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BrokerServer is the server API for Broker service.
// All implementations must embed UnimplementedBrokerServer
// for forward compatibility
type BrokerServer interface {
	Push(context.Context, *PushReq) (*PushReply, error)
	PushMsg(context.Context, *PushMsgReq) (*PushMsgReply, error)
	Broadcast(context.Context, *BroadcastReq) (*BroadcastReply, error)
	BroadcastRoom(context.Context, *BroadcastRoomReq) (*BroadcastRoomReply, error)
	Rooms(context.Context, *RoomsReq) (*RoomsReply, error)
	mustEmbedUnimplementedBrokerServer()
}

// UnimplementedBrokerServer must be embedded to have forward compatible implementations.
type UnimplementedBrokerServer struct {
}

func (UnimplementedBrokerServer) Push(context.Context, *PushReq) (*PushReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Push not implemented")
}
func (UnimplementedBrokerServer) PushMsg(context.Context, *PushMsgReq) (*PushMsgReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushMsg not implemented")
}
func (UnimplementedBrokerServer) Broadcast(context.Context, *BroadcastReq) (*BroadcastReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Broadcast not implemented")
}
func (UnimplementedBrokerServer) BroadcastRoom(context.Context, *BroadcastRoomReq) (*BroadcastRoomReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BroadcastRoom not implemented")
}
func (UnimplementedBrokerServer) Rooms(context.Context, *RoomsReq) (*RoomsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rooms not implemented")
}
func (UnimplementedBrokerServer) mustEmbedUnimplementedBrokerServer() {}

// UnsafeBrokerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BrokerServer will
// result in compilation errors.
type UnsafeBrokerServer interface {
	mustEmbedUnimplementedBrokerServer()
}

func RegisterBrokerServer(s grpc.ServiceRegistrar, srv BrokerServer) {
	s.RegisterService(&Broker_ServiceDesc, srv)
}

func _Broker_Push_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).Push(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broker_Push_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).Push(ctx, req.(*PushReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_PushMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushMsgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).PushMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broker_PushMsg_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).PushMsg(ctx, req.(*PushMsgReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_Broadcast_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BroadcastReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).Broadcast(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broker_Broadcast_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).Broadcast(ctx, req.(*BroadcastReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_BroadcastRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BroadcastRoomReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).BroadcastRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broker_BroadcastRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).BroadcastRoom(ctx, req.(*BroadcastRoomReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_Rooms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoomsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).Rooms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broker_Rooms_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).Rooms(ctx, req.(*RoomsReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Broker_ServiceDesc is the grpc.ServiceDesc for Broker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Broker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.broker.Broker",
	HandlerType: (*BrokerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Push",
			Handler:    _Broker_Push_Handler,
		},
		{
			MethodName: "PushMsg",
			Handler:    _Broker_PushMsg_Handler,
		},
		{
			MethodName: "Broadcast",
			Handler:    _Broker_Broadcast_Handler,
		},
		{
			MethodName: "BroadcastRoom",
			Handler:    _Broker_BroadcastRoom_Handler,
		},
		{
			MethodName: "Rooms",
			Handler:    _Broker_Rooms_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/broker/broker.proto",
}