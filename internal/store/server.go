package store

import (
	"github.com/zhihanii/im-pusher/internal/store/conf"
	"github.com/zhihanii/im-pusher/internal/store/data"
	"github.com/zhihanii/im-pusher/internal/store/server"
	"github.com/zhihanii/zlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
)

type appServer struct {
	c           *conf.Config
	storeServer *server.Server
}

type preparedServer struct {
	*appServer
}

func (s preparedServer) Run() error {
	s.storeServer.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	}

	return nil
}

func createServer(c *conf.Config) (*appServer, error) {
	zlog.Init(c.LoggerOptions)
	db, err := gorm.Open(mysql.Open(c.MySQLOptions.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	d := data.NewData(db)
	repo := data.NewStoreRepo(c, d)
	srv := server.New(c, repo)

	return &appServer{
		c:           c,
		storeServer: srv,
	}, nil
}

func (s *appServer) Prepare() preparedServer {
	return preparedServer{s}
}
