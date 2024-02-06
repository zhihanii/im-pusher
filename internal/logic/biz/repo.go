package biz

import (
	"context"
	"fmt"
	"github.com/zhihanii/im-pusher/api/protocol"
)

func EncodeRoomKey(t, room string) string {
	return fmt.Sprintf("%s://%s", t, room)
}

type LogicRepo interface {
	AddMapping(ctx context.Context, memberId uint64, key, server string) error
	ExpireMapping(ctx context.Context, memberId uint64, key string) (has bool, err error)
	DelMapping(ctx context.Context, memberId uint64, key, server string) (bool, error)
	//GetMapping(ctx context.Context, memberId uint64) (map[string]string, error)
	PushOfflineMsg(ctx context.Context, op int32, mids []uint64, msg []byte) error
	ServersByKeys(ctx context.Context, keys []string) ([]string, error)
	KeysByMids(ctx context.Context, mids []uint64) (map[string]string, []uint64, []uint64, error)
	AddServerOnline(ctx context.Context, server string, online *Online) error
	ServerOnline(ctx context.Context, server string) (*Online, error)
	DelServerOnline(ctx context.Context, server string) error
	StoreChatMessage(ctx context.Context, seq uint32, chatMessage *protocol.ChatSendMessage) error
	StoreGroupChatMessage(ctx context.Context, seq uint32, groupChatMessage *protocol.GroupChatMessage) ([]uint64, error)
	DeleteGroupChatOfflineMessage(ctx context.Context, groupId, userId, messageId uint64) error
	//PushChatMessageToKafka(ctx context.Context, chatMessage *protocol.ChatSendMessage)
	PushMessage(ctx context.Context, dm *protocol.DispatcherMessage)
}
