package rpc

import (
	"context"
	"github.com/zhihanii/im-pusher/api/logic"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var logicClient logic.LogicClient

func initLogicClient(etcdClient *clientv3.Client) {
	var err error
	//etcdResolver, err := resolver.NewBuilder(etcdClient)
	etcdResolver, err := registry.NewGRPCEtcdResolverBuilder(etcdClient)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial("etcd:///logic-service", grpc.WithResolvers(etcdResolver),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	logicClient = logic.NewLogicClient(conn)
}

func Connect(ctx context.Context, req *logic.ConnectReq) (memberId uint64, key, roomId string, accepts []int32, hb time.Duration, err error) {
	var reply *logic.ConnectReply
	reply, err = logicClient.Connect(ctx, req)
	if err != nil {
		return
	}
	return reply.MemberId, reply.Key, reply.RoomId, reply.Accepts, time.Duration(reply.Heartbeat), nil
}

func Disconnect(ctx context.Context, req *logic.DisconnectReq) (err error) {
	var reply *logic.DisconnectReply
	reply, err = logicClient.Disconnect(ctx, req)
	if err != nil {
		return
	}
	if reply.Has {

	}
	return nil
}

func Heartbeat(ctx context.Context, req *logic.HeartbeatReq) (err error) {
	//var reply *logic.HeartbeatReply
	_, err = logicClient.Heartbeat(ctx, req)
	if err != nil {
		return
	}
	return nil
}

func Receive(ctx context.Context, req *logic.ReceiveReq) (m *protocol.Message, err error) {
	var reply *logic.ReceiveReply
	reply, err = logicClient.Receive(ctx, req)
	if err != nil {
		return
	}
	return reply.Message, nil
}
