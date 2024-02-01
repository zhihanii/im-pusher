package kafka

import (
	"github.com/IBM/sarama"
	"github.com/zhihanii/zlog"
)

type Client struct {
	asyncProducer sarama.AsyncProducer
}

func NewClient(addrs []string, config *sarama.Config) (*Client, error) {
	var err error
	c := &Client{}
	c.asyncProducer, err = sarama.NewAsyncProducer(addrs, config)
	if err != nil {
		return nil, err
	}
	zlog.Info("async producer build successfully")
	go c.run()
	return c, nil
}

func (c *Client) SendMessageAsync(m *sarama.ProducerMessage) {
	c.asyncProducer.Input() <- m
}

func (c *Client) run() {
	for {
		select {
		case err := <-c.asyncProducer.Errors():
			zlog.Errorf("async producer:%v", err)
		case msg := <-c.asyncProducer.Successes():
			zlog.Infof("message send successfully, topic:%s, partition:%d, offset:%d",
				msg.Topic, msg.Partition, msg.Offset)
		}
	}
}
