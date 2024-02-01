package options

import (
	"fmt"
	"github.com/spf13/pflag"
)

type GrpcOptions struct {
	ServiceName string `json:"service-name" mapstructure:"service-name"`
	Network     string `json:"network" mapstructure:"network"`
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port" mapstructure:"bind-port"`
	Weight      int    `json:"weight" mapstructure:"weight"`
	MaxMsgSize  int    `json:"max-msg-size" mapstructure:"max-msg-size"`
}

func NewGrpcOptions() *GrpcOptions {
	return &GrpcOptions{
		ServiceName: "",
		Network:     "tcp",
		BindAddress: "0.0.0.0",
		BindPort:    8081,
		Weight:      10,
		MaxMsgSize:  4 * 1024 * 1024,
	}
}

func (o *GrpcOptions) Validate() []error {
	var errs []error

	if o.ServiceName == "" {
		errs = append(errs, fmt.Errorf("service-name can not be empty"))
	}

	if o.BindPort < 0 || o.BindPort > 65535 {
		errs = append(errs,
			fmt.Errorf(
				"--insecure-port %v must be between 0 and 65535, inclusive. 0 for turning off insecure (HTTP) port",
				o.BindPort,
			),
		)
	}

	return errs
}

func (o *GrpcOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ServiceName, "grpc.service-name", o.ServiceName, "The service name is for etcd registry.")

	fs.StringVar(&o.Network, "grpc.network", o.Network, "grpc network")

	fs.StringVar(&o.BindAddress, "grpc.bind-address", o.BindAddress, ""+
		"The IP address on which to serve the --grpc.port(set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")

	fs.IntVar(&o.BindPort, "grpc.bind-port", o.BindPort, ""+
		"The port on which to serve unsecured, unauthenticated grpc access. It is assumed "+
		"that firewall rules are set up such that this port is not reachable from outside of "+
		"the deployed machine and that port 443 on the iam public address is proxied to this "+
		"port. This is performed by nginx in the default setup. Set to zero to disable.")

	fs.IntVar(&o.Weight, "grpc.weight", o.Weight, "load balance weight")

	fs.IntVar(&o.MaxMsgSize, "grpc.max-msg-size", o.MaxMsgSize, "gRPC max message size.")
}
