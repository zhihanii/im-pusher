package biz

import (
	"context"
	"github.com/zhihanii/im-pusher/api/protocol"
)

type DispatcherRepo interface {
	GetMapping(ctx context.Context, memberId uint64) (res map[string]string, err error)
	PushMessage(ctx context.Context, keys []string, tm *protocol.TransMessage) error
}
