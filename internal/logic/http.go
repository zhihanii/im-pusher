package logic

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zhihanii/im-pusher/internal/logic/biz"
	"github.com/zhihanii/zlog"
	"net/http"
)

type httpServer struct {
	ctx            context.Context
	addr           string
	router         *gin.Engine
	insecureServer *http.Server
	lh             *biz.LogicHandler
}

func (s *httpServer) Run() error {
	s.insecureServer = &http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}
	go func() {
		zlog.Info("http server start")
		if err := s.insecureServer.ListenAndServe(); err != nil {
			zlog.Errorf("http server error:%v", err)
			return
		}
	}()
	return nil
}
