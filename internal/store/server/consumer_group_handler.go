package server

import (
	"github.com/IBM/sarama"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/zlog"
	"google.golang.org/protobuf/proto"
)

var _ sarama.ConsumerGroupHandler = (*ConsumerGroupHandler)(nil)

type ConsumerGroupHandler struct {
	s     *Server
	ready chan bool
}

func newConsumerGroupHandler(s *Server) *ConsumerGroupHandler {
	h := &ConsumerGroupHandler{
		s:     s,
		ready: make(chan bool),
	}
	return h
}

func (h *ConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var err error
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				zlog.Infoln("message channel was closed")
				return nil
			}
			session.MarkMessage(msg, "")
			cm := new(protocol.ChatMessage)
			err = proto.Unmarshal(msg.Value, cm)
			if err != nil {
				zlog.Errorf("unmarshal error:%v", err)
				continue
			}
			err = h.s.repo.StoreChatMessage(nil, 100, cm)
			if err != nil {
				zlog.Errorf("store chat message:%v", err)
			}

		case <-session.Context().Done():
			return nil
		}
	}
}
