package options

import (
	"encoding/json"
	"github.com/zhihanii/app/flag"
	coptions "github.com/zhihanii/im-pusher/internal/pkg/options"
	"github.com/zhihanii/zlog"
)

type Options struct {
	GrpcOptions            *coptions.GrpcOptions            `json:"grpc"  mapstructure:"grpc"`
	InsecureServingOptions *coptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	EtcdOptions            *coptions.EtcdOptions            `json:"etcd" mapstructure:"etcd"`
	LoggerOptions          *zlog.Options                    `json:"logger" mapstructure:"logger"`
}

func New() *Options {
	return &Options{
		GrpcOptions:            coptions.NewGrpcOptions(),
		InsecureServingOptions: coptions.NewInsecureServingOptions(),
		EtcdOptions:            coptions.NewEtcdOptions(),
		LoggerOptions:          zlog.NewOptions(),
	}
}

func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.GrpcOptions.Validate()...)
	errs = append(errs, o.InsecureServingOptions.Validate()...)
	errs = append(errs, o.EtcdOptions.Validate()...)
	errs = append(errs, o.LoggerOptions.Validate()...)

	return errs
}

func (o *Options) Flags() (fss flag.NamedFlagSets) {
	o.GrpcOptions.AddFlags(fss.FlagSet("grpc"))
	o.InsecureServingOptions.AddFlags(fss.FlagSet("insecure-serving"))
	o.EtcdOptions.AddFlags(fss.FlagSet("etcd"))
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
