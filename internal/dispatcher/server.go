package dispatcher

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zhihanii/im-pusher/internal/dispatcher/conf"
	"github.com/zhihanii/im-pusher/internal/dispatcher/data"
	"github.com/zhihanii/im-pusher/internal/dispatcher/server"
	"github.com/zhihanii/zlog"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"os/signal"
	"syscall"
)

type appServer struct {
	c                *conf.Config
	delayFunc        func()
	dispatcherServer *server.Server
}

type preparedServer struct {
	*appServer
}

func (s preparedServer) Run() error {
	s.delayFunc()

	s.dispatcherServer.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	}

	return nil
}

func createServer(c *conf.Config) (*appServer, error) {
	zlog.Init(c.LoggerOptions)
	cfg := clientv3.Config{
		Endpoints: c.EtcdOptions.Endpoints,
	}
	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	s := server.NewMessageSender(c, etcdClient)
	addr := fmt.Sprintf("%s:%d", c.RedisOptions.Host, c.RedisOptions.Port)
	zlog.Infof("redis addr: %s", addr)
	redisCli := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	d := data.NewData(redisCli)
	repo := data.NewDispatcherRepo(c, d)
	srv := server.New(c, s, repo)
	if err != nil {
		return nil, err
	}

	delayFunc := func() {
		s.Init()
	}

	//s.Init()

	return &appServer{
		c:                c,
		delayFunc:        delayFunc,
		dispatcherServer: srv,
	}, nil
}

func (s *appServer) Prepare() preparedServer {
	return preparedServer{s}
}
