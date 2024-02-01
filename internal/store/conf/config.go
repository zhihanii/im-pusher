package conf

import "github.com/zhihanii/im-pusher/internal/store/options"

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
