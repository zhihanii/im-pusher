package server

import (
	"github.com/zhenjl/cityhash"
	"github.com/zhihanii/im-pusher/internal/broker/conf"
)

type Server struct {
	c          *conf.Config
	round      *Round
	buckets    []*Bucket
	bucketSize uint32

	serverId string
}

func New(c *conf.Config) *Server {
	s := &Server{
		c:        c,
		round:    NewRound(c),
		serverId: "broker1",
	}
	s.buckets = make([]*Bucket, 10)
	s.bucketSize = 10
	for i := 0; i < 10; i++ {
		s.buckets[i] = NewBucket(c)
	}
	s.initWebsocket()
	return s
}

func (s *Server) Buckets() []*Bucket {
	return s.buckets
}

func (s *Server) Bucket(subKey string) *Bucket {
	idx := cityhash.CityHash32([]byte(subKey), uint32(len(subKey))) % s.bucketSize
	return s.buckets[idx]
}
