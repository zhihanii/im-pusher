package logic

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	pb "github.com/zhihanii/im-pusher/api/logic"
	"github.com/zhihanii/im-pusher/internal/logic/biz"
	"github.com/zhihanii/im-pusher/internal/logic/conf"
	"github.com/zhihanii/im-pusher/internal/logic/data"
	"github.com/zhihanii/im-pusher/internal/logic/data/kafka"
	"github.com/zhihanii/im-pusher/internal/logic/rpc"
	"github.com/zhihanii/im-pusher/internal/logic/service"
	"github.com/zhihanii/loadbalance"
	"github.com/zhihanii/registry"
	"github.com/zhihanii/zlog"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type appServer struct {
	c          *conf.Config
	grpcServer *grpcServer
	httpServer *httpServer
	etcdClient *clientv3.Client
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
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Flush.Frequency = 200 * time.Millisecond
	config.Producer.Flush.Messages = 200

	cfg := clientv3.Config{
		Endpoints: c.EtcdOptions.Endpoints,
	}
	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	builder, err := registry.NewEtcdWatchResolverBuilder(etcdClient)
	if err != nil {
		return nil, err
	}
	f := loadbalance.NewWatchBalancerFactory(func() loadbalance.LoadBalancer {
		return loadbalance.NewWeightedRoundRobinBalancer()
	}, builder)
	err = f.BuildBalancers(context.Background(), []string{"broker-service"})
	if err != nil {
		return nil, err
	}

	kafkaCli, err := kafka.NewClient(c.KafkaOptions.Brokers, config)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(mysql.Open(c.MySQLOptions.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	redisCli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", c.RedisOptions.Host, c.RedisOptions.Port),
	})
	d := data.NewData(kafkaCli, redisCli, db)
	repo := data.NewLogicRepo(c, d, time.Hour)
	lh, err := biz.NewLogicHandler(c, f, repo)
	if err != nil {
		return nil, err
	}
	svc, err := service.New(lh)
	if err != nil {
		return nil, err
	}
	pb.RegisterLogicServer(grpcSrv, svc)

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
		lh:   lh,
	}

	return &appServer{
		c:          c,
		grpcServer: gs,
		httpServer: hs,
		etcdClient: etcdClient,
	}, nil
}

func (s *appServer) Prepare() preparedServer {
	s.httpServer.initRouter()
	return preparedServer{s}
}
