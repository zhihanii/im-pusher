package rpc

import (
	"context"
	"github.com/zhihanii/im-pusher/api/dispatcher"
	"github.com/zhihanii/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var dispatcherClient dispatcher.DispatcherClient

func initDispatcherClient(etcdClient *clientv3.Client) {
	var err error
	//etcdResolver, err := resolver.NewBuilder(etcdClient)
	etcdResolver, err := registry.NewGRPCEtcdResolverBuilder(etcdClient)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial("etcd:///dispatcher-service", grpc.WithResolvers(etcdResolver),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	dispatcherClient = dispatcher.NewDispatcherClient(conn)
}

func Push(ctx context.Context, req *dispatcher.PushReq) (offlineMembers []uint64, err error) {
	var reply *dispatcher.PushReply
	reply, err = dispatcherClient.Push(ctx, req)
	if err != nil {
		return
	}
	return reply.OfflineMembers, nil
}
