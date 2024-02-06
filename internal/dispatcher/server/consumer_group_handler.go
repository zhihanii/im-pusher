package server

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/taskpool"
	"github.com/zhihanii/zlog"
	"google.golang.org/protobuf/proto"
)

var _ sarama.ConsumerGroupHandler = (*ConsumerGroupHandler)(nil)

type ConsumerGroupHandler struct {
	ctx   context.Context
	s     *Server
	ready chan bool
}

func newConsumerGroupHandler(s *Server) *ConsumerGroupHandler {
	h := &ConsumerGroupHandler{
		ctx:   context.Background(),
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
			dm := new(protocol.DispatcherMessage)
			err = proto.Unmarshal(msg.Value, dm)
			if err != nil {
				zlog.Errorf("unmarshal error:%v", err)
				continue
			}
			//zlog.Infof("handle push dm")
			//chatReceiveMessage := new(protocol.ChatReceiveMessage)
			//err = proto.Unmarshal(dm.Data, chatReceiveMessage)
			//if err != nil {
			//	zlog.Errorf("unmarshal error:%v", err)
			//	continue
			//}
			//zlog.Infof("receivers:%v, operation:%d, sequence:%d", dm.Receivers, dm.Operation, dm.Sequence)
			//zlog.Infof("chat receive message: %v", chatReceiveMessage)
			//h.s.repo.GetMapping(context.Background(), 1)
			ctx1 := context.Background()
			taskpool.Submit(ctx1, func() {
				h.s.HandlePush(h.ctx, dm.Receivers[0], dm)
			})

		case <-session.Context().Done():
			return nil
		}
	}
}
