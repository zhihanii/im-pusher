package options

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/spf13/pflag"
	"os"
)

type EtcdOptions struct {
	Endpoints            []string `json:"endpoints"               mapstructure:"endpoints"`
	Timeout              int      `json:"timeout"                 mapstructure:"timeout"`
	RequestTimeout       int      `json:"request-timeout"         mapstructure:"request-timeout"`
	LeaseExpire          int      `json:"lease-expire"            mapstructure:"lease-expire"`
	Username             string   `json:"username"                mapstructure:"username"`
	Password             string   `json:"password"                mapstructure:"password"`
	UseTLS               bool     `json:"use-tls"                 mapstructure:"use-tls"`
	CaCert               string   `json:"ca-cert"                 mapstructure:"ca-cert"`
	Cert                 string   `json:"cert"                    mapstructure:"cert"`
	Key                  string   `json:"key"                     mapstructure:"key"`
	HealthBeatPathPrefix string   `json:"health-beat-path-prefix" mapstructure:"health-beat-path-prefix"`
	HealthBeatIFaceName  string   `json:"health-beat-iface-name"  mapstructure:"health-beat-iface-name"`
	Namespace            string   `json:"namespace"               mapstructure:"namespace"`
}

func NewEtcdOptions() *EtcdOptions {
	return &EtcdOptions{
		Timeout:        5,
		RequestTimeout: 2,
		LeaseExpire:    5,
	}
}

func (o *EtcdOptions) Validate() []error {
	var errs []error

	if len(o.Endpoints) == 0 {
		errs = append(errs, fmt.Errorf("etcd endpoints can not be empty"))
	}

	if o.RequestTimeout <= 0 {
		errs = append(errs, fmt.Errorf("--etcd.request-timeout cannot be negative"))
	}

	return errs
}

func (o *EtcdOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&o.Endpoints, "etcd.endpoints", o.Endpoints, "Endpoints of etcd cluster.")
	fs.StringVar(&o.Username, "etcd.username", o.Username, "Username of etcd cluster.")
	fs.StringVar(&o.Password, "etcd.password", o.Password, "Password of etcd cluster.")
	fs.IntVar(&o.Timeout, "etcd.timeout", o.Timeout, "Etcd dial timeout in seconds.")
	fs.IntVar(&o.RequestTimeout, "etcd.request-timeout", o.RequestTimeout, "Etcd request timeout in seconds.")
	fs.IntVar(&o.LeaseExpire, "etcd.lease-expire", o.LeaseExpire, "Etcd expire timeout in seconds.")
	fs.BoolVar(&o.UseTLS, "etcd.use-tls", o.UseTLS, "Use tls transport to connect etcd cluster.")
	fs.StringVar(&o.CaCert, "etcd.ca-cert", o.CaCert, "Path to cacert for connecting to etcd cluster.")
	fs.StringVar(&o.Cert, "etcd.cert", o.Cert, "Path to cert file for connecting to etcd cluster.")
	fs.StringVar(&o.Key, "etcd.key", o.Key, "Path to key file for connecting to etcd cluster.")
	fs.StringVar(
		&o.HealthBeatPathPrefix,
		"etcd.health-beat-path-prefix",
		o.HealthBeatPathPrefix,
		"health beat path prefix.",
	)
	fs.StringVar(
		&o.HealthBeatIFaceName,
		"etcd.health-beat-iface-name",
		o.HealthBeatIFaceName,
		"health beat registry iface name, such as eth0.",
	)
	fs.StringVar(&o.Namespace, "etcd.namespace", o.Namespace, "Etcd storage namespace.")
}

func (o *EtcdOptions) GetEtcdTLSConfig() (*tls.Config, error) {
	var (
		cert       tls.Certificate
		certLoaded bool
		capool     *x509.CertPool
	)
	if o.Cert != "" && o.Key != "" {
		var err error
		cert, err = tls.LoadX509KeyPair(o.Cert, o.Key)
		if err != nil {
			return nil, err
		}
		certLoaded = true
		o.UseTLS = true
	}
	if o.CaCert != "" {
		data, err := os.ReadFile(o.CaCert)
		if err != nil {
			return nil, err
		}
		capool = x509.NewCertPool()
		for {
			var block *pem.Block
			block, _ = pem.Decode(data)
			if block == nil {
				break
			}
			cacert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			capool.AddCert(cacert)
		}
		o.UseTLS = true
	}

	if o.UseTLS {
		cfg := &tls.Config{
			RootCAs:            capool,
			InsecureSkipVerify: false,
		}
		if certLoaded {
			cfg.Certificates = []tls.Certificate{cert}
		}

		return cfg, nil
	}

	return &tls.Config{}, nil
}
