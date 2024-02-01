package db

import (
	"github.com/zhihanii/im-pusher/internal/pkg/database"
	"time"
)

type ChatConversation struct {
	database.Base
	User1 uint64 `json:"user_1" gorm:"column:user_1"`
	User2 uint64 `json:"user_2" gorm:"column:user_2"`
}

type ChatMessage struct {
	database.Base
	ConversationId uint64 `json:"conversation_id" gorm:"column:conversation_id"`
	FromUserId     uint64 `json:"from_user_id" gorm:"column:from_user_id"`
	ToUserId       uint64 `json:"to_user_id" gorm:"column:to_user_id"`
	Sequence       uint32 `json:"sequence" gorm:"column:sequence"`
	Content        string `json:"content" gorm:"column:content"`
}

type GroupChatMessage struct {
	database.Base
	GroupId   uint64 `json:"group_id" gorm:"column:group_id"`
	SenderId  uint64 `json:"sender_id" gorm:"column:sender_id"`
	MessageId uint64 `json:"msg_id" gorm:"column:msg_id"`
	Content   string `json:"content" gorm:"column:content"`
}

type GroupChatOfflineMessage struct {
	Id        uint64    `json:"id" gorm:"column:id"`
	GroupId   uint64    `json:"group_id" gorm:"column:group_id"`
	UserId    uint64    `json:"user_id" gorm:"column:user_id"`
	MessageId uint64    `json:"msg_id" gorm:"column:msg_id"`
	GMTCreate time.Time `json:"gmt_create" gorm:"column:gmt_create"`
}
