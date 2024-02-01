package options

import (
	"encoding/json"
	"github.com/zhihanii/app/flag"
	coptions "github.com/zhihanii/im-pusher/internal/pkg/options"
	"github.com/zhihanii/zlog"
	"time"
)

type Options struct {
	GrpcOptions            *coptions.GrpcOptions            `json:"grpc"  mapstructure:"grpc"`
	InsecureServingOptions *coptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	EtcdOptions            *coptions.EtcdOptions            `json:"etcd"  mapstructure:"etcd"`
	KafkaOptions           *coptions.KafkaOptions           `json:"kafka" mapstructure:"kafka"`
	RedisOptions           *coptions.RedisOptions           `json:"redis" mapstructure:"redis"`
	LoggerOptions          *zlog.Options                    `json:"logger" mapstructure:"logger"`
	ServerOptions          *ServerOptions
}

func New() *Options {
	o := &Options{
		GrpcOptions:            coptions.NewGrpcOptions(),
		InsecureServingOptions: coptions.NewInsecureServingOptions(),
		EtcdOptions:            coptions.NewEtcdOptions(),
		KafkaOptions:           coptions.NewKafkaOptions(),
		RedisOptions:           coptions.NewRedisOptions(),
		LoggerOptions:          zlog.NewOptions(),
		ServerOptions:          &ServerOptions{},
	}
	o.ServerOptions.Flush.Frequency = time.Second
	o.ServerOptions.Flush.Bytes = 0
	o.ServerOptions.Flush.Messages = 0
	return o
}

func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.GrpcOptions.Validate()...)
	errs = append(errs, o.InsecureServingOptions.Validate()...)
	errs = append(errs, o.EtcdOptions.Validate()...)
	errs = append(errs, o.KafkaOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.LoggerOptions.Validate()...)

	return errs
}

func (o *Options) Flags() (fss flag.NamedFlagSets) {
	o.GrpcOptions.AddFlags(fss.FlagSet("grpc"))
	o.InsecureServingOptions.AddFlags(fss.FlagSet("insecure-serving"))
	o.EtcdOptions.AddFlags(fss.FlagSet("etcd"))
	o.KafkaOptions.AddFlags(fss.FlagSet("kafka"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
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
