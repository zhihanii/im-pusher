package biz

import (
	"context"
	"errors"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/im-pusher/internal/dispatcher/conf"
	"github.com/zhihanii/zlog"
)

var (
	MemberOfflineErr = errors.New("member offline")
)

type DispatcherHandler struct {
	c *conf.Config

	s *MessageSender

	repo DispatcherRepo
}

func NewDispatchHandler(c *conf.Config, s *MessageSender, repo DispatcherRepo) (*DispatcherHandler, error) {
	dh := &DispatcherHandler{
		c:    c,
		s:    s,
		repo: repo,
	}
	return dh, nil
}

func (h *DispatcherHandler) HandlePush(ctx context.Context, receiverId uint64, m *protocol.Message) (err error) {
	zlog.Infof("handle push, receiver id:%d", receiverId)
	res, err := h.repo.GetMapping(ctx, receiverId)
	if err != nil {
		return err
	}
	if len(res) == 0 {
		return MemberOfflineErr
	}
	//zlog.Infoln("get mapping res:")
	//for k, v := range res {
	//	zlog.Infof("key:%s value:%s", k, v)
	//}
	serverKeys := make(map[string][]string)
	for key, server := range res {
		if key != "" && server != "" {
			serverKeys[server] = append(serverKeys[server], key)
		}
	}
	for server, keys := range serverKeys {
		if err = h.push(ctx, protocol.Channel, 0, server, keys, m.Operation, m.Sequence, m.Data); err != nil {
			return
		}
	}
	return
}

func (h *DispatcherHandler) push(ctx context.Context, t uint32, priority int, server string, keys []string,
	operation uint32, sequence uint32, data []byte) (err error) {
	msg := &protocol.TransMessage{
		Type:      t,
		Priority:  uint32(priority),
		Server:    server,
		Keys:      keys,
		Operation: operation,
		Sequence:  sequence,
		Data:      data,
	}
	//return h.repo.PushMessage(ctx, keys, msg)
	h.s.send(msg)
	return nil
}
