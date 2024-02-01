package broker

import (
	"github.com/zhihanii/app"
	"github.com/zhihanii/im-pusher/internal/broker/conf"
	"github.com/zhihanii/im-pusher/internal/broker/options"
)

func NewApp(name string) *app.App {
	opts := options.New()
	a := app.New(name,
		"",
		app.WithOptions(opts),
		app.WithDescription(""),
		app.WithDefaultArgs(),
		app.WithRunFunc(run(opts)),
	)
	return a
}

func run(opts *options.Options) app.RunFunc {
	return func(name string) error {
		cfg, err := conf.New(opts)
		if err != nil {
			return err
		}
		return Run(cfg)
	}
}
