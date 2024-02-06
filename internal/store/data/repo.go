package data

import (
	"context"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/im-pusher/internal/pkg/database"
	"github.com/zhihanii/im-pusher/internal/store/conf"
	"github.com/zhihanii/im-pusher/internal/store/data/db"
	"github.com/zhihanii/zlog"
	"time"
)

type StoreRepo interface {
	StoreChatMessage(ctx context.Context, seq uint32, chatMessage *protocol.ChatSendMessage) error
}

type storeRepo struct {
	c    *conf.Config
	data *Data

	msgCh  chan *db.ChatMessage
	worker *worker
	id     uint64
}

func NewStoreRepo(c *conf.Config, data *Data) StoreRepo {
	r := &storeRepo{
		c:    c,
		data: data,
	}
	r.initId()
	r.msgCh = make(chan *db.ChatMessage, 10000)
	r.worker = newWorker(c, r, r.msgCh)
	return r
}

func (r *storeRepo) initId() {
	//在表数据量很大时, 响应时间会较长(319W响应时间=1.67s)
	//todo 若kafka的消息未消费完, 则会造成主键冲突的问题
	err := r.data.db.Raw("select count(*) from t_chat_message;").Scan(&r.id).Error
	if err != nil {
		zlog.Errorf("init id:%v", err)
		return
	}
	zlog.Infof("init id:%d", r.id)
}

// StoreChatMessage 非并发安全
func (r *storeRepo) StoreChatMessage(ctx context.Context, seq uint32, chatMessage *protocol.ChatSendMessage) error {
	now := time.Now()
	cm := &db.ChatMessage{
		Base: database.Base{
			Id:          r.id,
			IsDeleted:   0,
			GMTCreate:   now,
			GMTModified: now,
		},
		ConversationId: chatMessage.ConversationId,
		FromUserId:     chatMessage.From,
		ToUserId:       chatMessage.To,
		Sequence:       seq,
		Content:        string(chatMessage.Content),
	}
	r.id++

	//if err := r.data.db.Table("t_chat_message").Create(cm).Error; err != nil {
	//	return err
	//}
	r.msgCh <- cm

	return nil
}

func (r *storeRepo) storeBatchChatMessage(ctx context.Context, msgs []*db.ChatMessage) error {
	startTime := time.Now()
	err := r.data.db.Table("t_chat_message").Create(&msgs).Error
	if err != nil {
		return err
	}
	endTime := time.Now()
	responseTime := endTime.Sub(startTime)
	zlog.Infof("store batch chat message to db response time:%d ms", responseTime.Milliseconds())
	return nil
}
