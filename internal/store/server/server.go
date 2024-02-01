package server

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"github.com/zhihanii/im-pusher/internal/store/conf"
	"github.com/zhihanii/im-pusher/internal/store/data"
	"github.com/zhihanii/zlog"
	"sync"
)

type Server struct {
	c        *conf.Config
	repo     data.StoreRepo
	consumer *ConsumerGroupHandler
}

func New(c *conf.Config, repo data.StoreRepo) *Server {
	s := &Server{
		c:    c,
		repo: repo,
	}
	s.consumer = newConsumerGroupHandler(s)
	return s
}

func (s *Server) Run() {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	client, err := sarama.NewConsumerGroup(s.c.KafkaOptions.Brokers, "group1", config)
	if err != nil {
		panic(err)
	}
	zlog.Infoln("build consumer group successfully")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(context.Background(), []string{"topic_chat_message"}, s.consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				zlog.Panicf("Error from consumer: %v", err)
			}
		}
	}()
	<-s.consumer.ready
	wg.Wait()
	if err = client.Close(); err != nil {
		zlog.Panicf("Error closing client: %v", err)
	}
}
