package broker

import (
	"context"
	"fmt"
	pb "github.com/zhihanii/im-pusher/api/broker"
	"github.com/zhihanii/im-pusher/internal/broker/conf"
	"github.com/zhihanii/im-pusher/internal/broker/rpc"
	"github.com/zhihanii/im-pusher/internal/broker/server"
	"github.com/zhihanii/im-pusher/internal/broker/service"
	"github.com/zhihanii/registry"
	"github.com/zhihanii/zlog"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
)

type appServer struct {
	c             *conf.Config
	grpcServer    *grpcServer
	httpServer    *httpServer
	connectServer *server.Server
	etcdClient    *clientv3.Client
}

type preparedServer struct {
	*appServer
}

func (s preparedServer) Run() error {
	s.grpcServer.Run()
	_ = s.httpServer.Run()

	rpc.Init(s.etcdClient)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	}

	return nil
}

func createServer(c *conf.Config) (*appServer, error) {
	zlog.Init(c.LoggerOptions)
	opts := []grpc.ServerOption{}
	grpcSrv := grpc.NewServer(opts...)
	srv := server.New(c)
	svc, err := service.New(srv)
	if err != nil {
		return nil, err
	}
	pb.RegisterBrokerServer(grpcSrv, svc)

	cfg := clientv3.Config{
		Endpoints: c.EtcdOptions.Endpoints,
	}
	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	gs := &grpcServer{
		Server: grpcSrv,
		endpoint: &registry.Endpoint{
			ServiceName: c.GrpcOptions.ServiceName,
			Network:     c.GrpcOptions.Network,
			Address:     c.GrpcOptions.BindAddress,
			Port:        c.GrpcOptions.BindPort,
			Weight:      c.GrpcOptions.Weight,
		},
		//etcdClient: etcdClient,
	}

	gs.etcdRegistry, err = registry.NewEtcdRegistryWithClient(etcdClient, int64(c.EtcdOptions.LeaseExpire))
	if err != nil {
		return nil, err
	}

	hs := &httpServer{
		ctx:  context.Background(),
		addr: fmt.Sprintf("%s:%d", c.InsecureServingOptions.BindAddress, c.InsecureServingOptions.BindPort),
	}

	return &appServer{
		c:          c,
		grpcServer: gs,
		httpServer: hs,
		etcdClient: etcdClient,
	}, nil
}

func (s *appServer) Prepare() preparedServer {
	s.httpServer.init()
	return preparedServer{s}
}
