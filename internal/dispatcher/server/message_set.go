package server

import (
	"github.com/zhihanii/im-pusher/api/broker"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/im-pusher/internal/dispatcher/conf"
)

type messageSet struct {
	c           *conf.Config
	msgs        []*protocol.TransMessage
	bufferBytes int
	bufferCount int
}

func newMessageSet(c *conf.Config) *messageSet {
	return &messageSet{
		c: c,
	}
}

func (s *messageSet) add(msg *protocol.TransMessage) error {
	//var err error
	s.msgs = append(s.msgs, msg)
	s.bufferCount++
	return nil
}

func (s *messageSet) buildRequest() *broker.PushReq {
	return &broker.PushReq{
		Messages: s.msgs,
	}
}

func (s *messageSet) readyToFlush() bool {
	switch {
	case s.empty():
		return false
	case s.c.ServerOptions.Flush.Frequency == 0 && s.c.ServerOptions.Flush.Bytes == 0 && s.c.ServerOptions.Flush.Messages == 0:
		return true
	// If we've passed the message trigger-point
	case s.c.ServerOptions.Flush.Messages > 0 && s.bufferCount >= s.c.ServerOptions.Flush.Messages:
		return true
	// If we've passed the byte trigger-point
	case s.c.ServerOptions.Flush.Bytes > 0 && s.bufferBytes >= s.c.ServerOptions.Flush.Bytes:
		return true
	default:
		return false
	}
}

func (s *messageSet) empty() bool {
	return s.bufferCount == 0
}
