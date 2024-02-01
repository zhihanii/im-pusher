package conf

import (
	"github.com/zhihanii/im-pusher/internal/broker/options"
	"github.com/zhihanii/registry"
)

type Option func(c *Config)

type Config struct {
	*options.Options

	Registry     registry.Registry
	RegistryInfo registry.Endpoint

	Protocol Protocol
}

func New(o *options.Options, opts ...Option) (*Config, error) {
	c := &Config{Options: o}
	for _, opt := range opts {
		opt(c)
	}
	c.Protocol = Protocol{
		BufferSize:         10,
		MessageChannelSize: 10,
	}
	return c, nil
}

//func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
//	return &Config{opts}, nil
//}

type Protocol struct {
	BufferSize         int
	MessageChannelSize int
}
