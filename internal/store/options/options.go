package options

import (
	"encoding/json"
	"github.com/zhihanii/app/flag"
	coptions "github.com/zhihanii/im-pusher/internal/pkg/options"
	"github.com/zhihanii/zlog"
	"time"
)

type Options struct {
	KafkaOptions  *coptions.KafkaOptions `json:"kafka" mapstructure:"kafka"`
	MySQLOptions  *coptions.MySQLOptions `json:"mysql" mapstructure:"mysql"`
	LoggerOptions *zlog.Options          `json:"logger" mapstructure:"logger"`
	ServerOptions *ServerOptions
}

func New() *Options {
	o := &Options{
		KafkaOptions:  coptions.NewKafkaOptions(),
		MySQLOptions:  coptions.NewMySQLOptions(),
		LoggerOptions: zlog.NewOptions(),
		ServerOptions: &ServerOptions{},
	}
	o.ServerOptions.Flush.Frequency = 200 * time.Millisecond
	o.ServerOptions.Flush.Bytes = 0
	o.ServerOptions.Flush.Messages = 200
	return o
}

func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.KafkaOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.LoggerOptions.Validate()...)

	return errs
}

func (o *Options) Flags() (fss flag.NamedFlagSets) {
	o.KafkaOptions.AddFlags(fss.FlagSet("kafka"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.LoggerOptions.AddFlags(fss.FlagSet("logger"))
	return fss
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)
	return string(data)
}

func (o *Options) Complete() error {
	return nil
}

type ServerOptions struct {
	Flush struct {
		Frequency time.Duration
		Bytes     int
		Messages  int
	}
}
