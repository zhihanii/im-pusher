package options

import "github.com/spf13/pflag"

type KafkaOptions struct {
	Brokers []string `json:"brokers" mapstructure:"brokers"`
}

func NewKafkaOptions() *KafkaOptions {
	return &KafkaOptions{}
}

func (o *KafkaOptions) Validate() []error {
	var errs []error
	return errs
}

func (o *KafkaOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&o.Brokers, "kafka.brokers", o.Brokers, "The Kafka Broker's addresses.")
}
