package server

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/im-pusher/internal/dispatcher/conf"
	"github.com/zhihanii/im-pusher/internal/dispatcher/data"
	"github.com/zhihanii/zlog"
	"sync"
)

type Server struct {
	c        *conf.Config
	consumer *ConsumerGroupHandler
	//dh       *DispatcherHandler
	sender *MessageSender
	repo   data.DispatcherRepo
}

func New(c *conf.Config, sender *MessageSender, repo data.DispatcherRepo) *Server {
	s := &Server{
		c:      c,
		sender: sender,
		repo:   repo,
	}
	s.consumer = newConsumerGroupHandler(s)
	return s
}

func (s *Server) Run() {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	client, err := sarama.NewConsumerGroup(s.c.KafkaOptions.Brokers, "group2", config)
	if err != nil {
		panic(err)
	}
	zlog.Infoln("build consumer group successfully")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(context.Background(), []string{"topic_dispatcher_message"}, s.consumer); err != nil {
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

var (
	MemberOfflineErr = errors.New("member offline")
)

func (s *Server) HandlePush(ctx context.Context, receiverId uint64, dm *protocol.DispatcherMessage) (err error) {
	//receiverId := dm.Receivers[0]
	//zlog.Infof("handle push, receiver id:%d", receiverId)
	res, err := s.repo.GetMapping(ctx, receiverId)
	if err != nil {
		zlog.Errorf("get mapping: %v", err)
		return err
	}
	if len(res) == 0 {
		zlog.Infof("member[%d] offline", receiverId)
		return MemberOfflineErr
	}
	serverKeys := make(map[string][]string)
	for key, server := range res {
		if key != "" && server != "" {
			serverKeys[server] = append(serverKeys[server], key)
		}
	}
	for server, keys := range serverKeys {
		_ = s.push(ctx, protocol.Channel, 0, server, keys, dm.Operation, dm.Sequence, dm.Data)
	}
	return
}

func (s *Server) push(ctx context.Context, t uint32, priority int, server string, keys []string,
	operation uint32, sequence uint32, data []byte) (err error) {
	tm := &protocol.TransMessage{
		Type:      t,
		Priority:  uint32(priority),
		Server:    server,
		Keys:      keys,
		Operation: operation,
		Sequence:  sequence,
		Data:      data,
	}
	s.sender.send(tm)
	return nil
}
