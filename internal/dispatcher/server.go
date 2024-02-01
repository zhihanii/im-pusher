package dispatcher

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	pb "github.com/zhihanii/im-pusher/api/dispatcher"
	"github.com/zhihanii/im-pusher/internal/dispatcher/biz"
	"github.com/zhihanii/im-pusher/internal/dispatcher/conf"
	"github.com/zhihanii/im-pusher/internal/dispatcher/data"
	"github.com/zhihanii/im-pusher/internal/dispatcher/service"
	"github.com/zhihanii/registry"
	"github.com/zhihanii/zlog"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
)

type appServer struct {
	c          *conf.Config
	grpcServer *grpcServer
	delayFunc  func()
}

type preparedServer struct {
	*appServer
}

func (s preparedServer) Run() error {
	s.grpcServer.Run()

	s.delayFunc()

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
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	cfg := clientv3.Config{
		Endpoints: c.EtcdOptions.Endpoints,
	}
	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	s := biz.NewMessageSender(c, etcdClient)
	redisCli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", c.RedisOptions.Host, c.RedisOptions.Port),
	})
	d := data.NewData(nil, redisCli)
	repo := data.NewDispatcherRepo(c, d)
	dh, err := biz.NewDispatchHandler(c, s, repo)
	if err != nil {
		return nil, err
	}
	svc, err := service.New(dh)
	if err != nil {
		return nil, err
	}
	pb.RegisterDispatcherServer(grpcSrv, svc)

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

	delayFunc := func() {
		s.Init()
	}

	return &appServer{
		c:          c,
		grpcServer: gs,
		delayFunc:  delayFunc,
	}, nil
}

func (s *appServer) Prepare() preparedServer {
	return preparedServer{s}
}
