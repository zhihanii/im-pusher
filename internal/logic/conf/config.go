package conf

import (
	"github.com/zhihanii/im-pusher/internal/logic/options"
	"time"
)

type Option func(c *Config)

type Config struct {
	*options.Options
}

func New(o *options.Options, opts ...Option) (*Config, error) {
	c := &Config{Options: o}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

type Node struct {
	DefaultDomain string
	HostDomain    string
	TCPPort       int
	WSPort        int
	WSSPort       int
	HeartbeatMax  int
	Heartbeat     time.Duration
	RegionWeight  float64
}

type Backoff struct {
	MaxDelay  int32
	BaseDelay int32
	Factor    float32
	Jitter    float32
}
