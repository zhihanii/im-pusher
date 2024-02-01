package broker

import "github.com/zhihanii/im-pusher/internal/broker/conf"

func Run(cfg *conf.Config) error {
	s, err := createServer(cfg)
	if err != nil {
		return err
	}
	return s.Prepare().Run()
}
