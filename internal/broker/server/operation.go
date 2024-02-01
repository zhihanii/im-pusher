package server

import (
	"context"
	"github.com/zhihanii/im-pusher/api/logic"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/im-pusher/internal/broker/rpc"
	"github.com/zhihanii/zlog"
	"strconv"
	"strings"
	"time"
)

// Connect 用户验证数据交由logic层解析
func (s *Server) Connect(ctx context.Context, m *protocol.Message) (memberId uint64, key, roomId string, accepts []int32, hb time.Duration, err error) {
	return rpc.Connect(ctx, &logic.ConnectReq{
		Server: s.serverId,
		Token:  m.Data,
	})
}

func (s *Server) Disconnect(ctx context.Context, mid uint64, key string) (err error) {
	return rpc.Disconnect(ctx, &logic.DisconnectReq{
		Server:   s.serverId,
		MemberId: mid,
		Key:      key,
	})
}

// Heartbeat 可以保持用户与Connect Server在Redis中的映射关系
func (s *Server) Heartbeat(ctx context.Context, mid uint64, key string) (err error) {
	return rpc.Heartbeat(ctx, &logic.HeartbeatReq{
		Server:   s.serverId,
		MemberId: mid,
		Key:      key,
	})
}

type OperateResult struct {
	Err     error
	Message *protocol.Message
}

func (s *Server) Operate(ctx context.Context, m *protocol.Message, ch *Channel, b *Bucket) (res OperateResult) {
	switch m.Operation {
	case protocol.ChangeRoom:
		if err := b.ChangeRoom(string(m.Data), ch); err != nil {

		}
		m.Operation = protocol.ChangeRoomReply
	case protocol.Sub:
		if ops, err := splitInt32(string(m.Data), ","); err == nil {
			ch.Watch(ops...)
		}
		m.Operation = protocol.SubReply
	case protocol.Unsub:
		if ops, err := splitInt32(string(m.Data), ","); err == nil {
			ch.UnWatch(ops...)
		}
		m.Operation = protocol.UnsubReply
	default:
		replyMessage, err := rpc.Receive(ctx, &logic.ReceiveReq{
			MemberId: ch.MemberId,
			Message:  m,
		})
		if err != nil {
			zlog.Errorf("logic client receive error:%v", err)
			res.Err = err
			return
		}
		res.Message = replyMessage
	}
	return
}

func splitInt32(s, p string) ([]int32, error) {
	if s == "" {
		return nil, nil
	}
	strs := strings.Split(s, p)
	res := make([]int32, len(strs))
	for _, str := range strs {
		i, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return nil, err
		}
		res = append(res, int32(i))
	}
	return res, nil
}
