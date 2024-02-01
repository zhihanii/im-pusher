package broker

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zhihanii/zlog"
	"net/http"
)

type httpServer struct {
	ctx  context.Context
	addr string
	//InsecureServingInfo *InsecureServingInfo
	//router              *gin.Engine
	mux            *http.ServeMux
	insecureServer *http.Server
}

func (s *httpServer) init() {
	s.mux = http.NewServeMux()
	s.mux.Handle("/metrics", promhttp.Handler())
}

func (s *httpServer) Run() error {
	s.insecureServer = &http.Server{
		Addr:    s.addr,
		Handler: s.mux,
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

//type InsecureServingInfo struct {
//	Addr string
//	Port int
//}
