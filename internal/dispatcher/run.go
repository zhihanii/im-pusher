package dispatcher

import "github.com/zhihanii/im-pusher/internal/dispatcher/conf"

func Run(c *conf.Config) error {
	s, err := createServer(c)
	if err != nil {
		return err
	}
	return s.Prepare().Run()
}
