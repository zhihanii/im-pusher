package biz

import (
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/im-pusher/internal/dispatcher/conf"
	"github.com/zhihanii/zlog"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// todo 使用WatchResolver监听connect-service, 当有节点变化时, 需要增加或删除Connect对象

type MessageSender struct {
	c *conf.Config
	//ctx context.Context
	msgCh      chan *protocol.TransMessage
	etcdClient *clientv3.Client
	brokers    map[string]*Broker
}

func NewMessageSender(c *conf.Config, etcdClient *clientv3.Client) *MessageSender {
	s := &MessageSender{
		c:          c,
		msgCh:      make(chan *protocol.TransMessage, 10),
		etcdClient: etcdClient,
		brokers:    make(map[string]*Broker),
	}
	return s
}

func (s *MessageSender) Init() {
	s.brokers["broker1"] = newBroker(s.c, "broker1", s.etcdClient)
	s.start()
}

func (s *MessageSender) send(tm *protocol.TransMessage) {
	s.msgCh <- tm
}

func (s *MessageSender) start() {
	go s.run()
}

func (s *MessageSender) run() {
	for {
		select {
		case tm, ok := <-s.msgCh:
			if !ok {
				zlog.Infoln("message channel was closed")
				return
			}

			broker, ok := s.brokers[tm.Server]
			if !ok {
				//返回logic层以重发消息
				zlog.Errorf("invalid server:%s", tm.Server)
			} else {
				broker.Push(tm)
			}
		}
	}
}
