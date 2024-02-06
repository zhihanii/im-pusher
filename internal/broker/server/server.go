package server

import (
	"fmt"
	"github.com/zhenjl/cityhash"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/im-pusher/internal/broker/conf"
	"github.com/zhihanii/zlog"
	"github.com/zhihanii/znet"
	"net"
)

type Server struct {
	c          *conf.Config
	round      *Round
	buckets    []*Bucket
	bucketSize uint32

	serverId string

	msgCh chan *protocol.TransMessage

	evl znet.EventLoop
}

func New(c *conf.Config) *Server {
	s := &Server{
		c:        c,
		round:    NewRound(c),
		serverId: "broker1",
		msgCh:    make(chan *protocol.TransMessage, 10000),
	}
	s.buckets = make([]*Bucket, 10)
	s.bucketSize = 10
	for i := 0; i < 10; i++ {
		s.buckets[i] = NewBucket(c)
	}
	//s.initWebsocket()
	go s.dispatch()
	return s
}

func (s *Server) Start() {
	var listener net.Listener
	var err error
	listener, err = net.Listen("tcp", fmt.Sprintf("%s:8082", s.c.InsecureServingOptions.BindAddress))
	if err != nil {
		zlog.Errorf("listen: %v", err)
		return
	}
	s.evl = znet.NewEventLoop(s)
	go func() {
		if err := s.evl.Serve(listener); err != nil {
			zlog.Errorf("event loop serve: %v", err)
		}
	}()
}

func (s *Server) Buckets() []*Bucket {
	return s.buckets
}

func (s *Server) Bucket(subKey string) *Bucket {
	idx := cityhash.CityHash32([]byte(subKey), uint32(len(subKey))) % s.bucketSize
	return s.buckets[idx]
}

func (s *Server) Push(tm *protocol.TransMessage) {
	s.msgCh <- tm
}

func (s *Server) dispatch() {
	var err error
	for tm := range s.msgCh {
		for _, key := range tm.Keys {
			bucket := s.Bucket(key)
			if bucket == nil {
				continue
			}
			if ch := bucket.Channel(key); ch != nil {
				//if !ch.NeedPush() {
				//	continue
				//}
				//todo 构造protocol.Message
				msg := &protocol.Message{
					Operation: tm.Operation,
					Sequence:  tm.Sequence,
					Data:      tm.Data,
				}

				//err = proto.Unmarshal(tm.Data, msg)
				//if err != nil {
				//	return &pb.PushReply{}, err
				//}

				if err = ch.Push(msg); err != nil {
					zlog.Errorf("channel push message: %v", err)
				}
			}
		}
	}
}
