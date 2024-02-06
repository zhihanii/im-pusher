package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	"github.com/zhenjl/cityhash"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/im-pusher/internal/logic/biz"
	"github.com/zhihanii/im-pusher/internal/logic/conf"
	"github.com/zhihanii/im-pusher/internal/logic/data/db"
	"github.com/zhihanii/im-pusher/internal/pkg/database"
	"github.com/zhihanii/zlog"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"strconv"
	"time"
)

const (
	_prefixMidServer    = "mid:%d" // mid -> key:server
	_prefixKeyServer    = "key:%s" // key -> server
	_prefixServerOnline = "ol:%s"  // server -> online
	_prefixMidOffline   = "offline:%d"
)

func keyMidServer(mid uint64) string {
	return fmt.Sprintf(_prefixMidServer, mid)
}

func keyKeyServer(key string) string {
	return fmt.Sprintf(_prefixKeyServer, key)
}

func keyServerOnline(key string) string {
	return fmt.Sprintf(_prefixServerOnline, key)
}

func keyMidOffline(mid uint64) string {
	return fmt.Sprintf(_prefixMidOffline, mid)
}

type logicRepo struct {
	c           *conf.Config
	data        *Data
	redisExpire time.Duration

	msgCh chan *protocol.DispatcherMessage
	//workers []*worker
}

func NewLogicRepo(c *conf.Config, data *Data, redisExpire time.Duration) biz.LogicRepo {
	r := &logicRepo{
		c:           c,
		data:        data,
		redisExpire: redisExpire,
	}
	r.msgCh = make(chan *protocol.DispatcherMessage, 10000)
	//for i := 0; i < 1; i++ {
	//	r.workers = append(r.workers, newWorker(c, r, r.msgCh))
	//}
	go r.run()
	return r
}

func (r *logicRepo) run() {
	var err error
	for msg := range r.msgCh {
		err = r.pushMessage(nil, msg)
		if err != nil {
			zlog.Errorf("push message: %v", err)
		}
	}
}

// 可以根据redis中是否存在映射关系来判断用户是否在线

func (r *logicRepo) AddMapping(ctx context.Context, memberId uint64, key, server string) error {
	_, err := r.data.redisCli.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, keyMidServer(memberId), key, server)
		pipe.Expire(ctx, keyMidServer(memberId), r.redisExpire)
		pipe.Set(ctx, keyKeyServer(key), server, r.redisExpire)
		return nil
	})
	return err
}

func (r *logicRepo) ExpireMapping(ctx context.Context, memberId uint64, key string) (has bool, err error) {
	cmds, err := r.data.redisCli.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		if memberId > 0 {
			pipe.Expire(ctx, keyMidServer(memberId), r.redisExpire)
		}
		pipe.Expire(ctx, keyKeyServer(key), r.redisExpire)
		return nil
	})
	for _, cmd := range cmds {
		has = cmd.(*redis.BoolCmd).Val()
	}
	return
}

func (r *logicRepo) DelMapping(ctx context.Context, memberId uint64, key, server string) (has bool, err error) {
	_, err = r.data.redisCli.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HDel(ctx, keyMidServer(memberId), key)
		pipe.Del(ctx, keyKeyServer(key))
		return nil
	})
	return
}

//func (r *logicRepo) GetMapping(ctx context.Context, memberId uint64) (res map[string]string, err error) {
//	strs, err := getMappingScript.Run(ctx, r.data.redisCli,
//		[]string{keyMidServer(memberId)}).StringSlice()
//	if err != nil {
//
//	}
//	n := len(strs)
//	if n > 0 && n%2 == 1 {
//		//log
//	}
//	for i := 0; i < n-1; i += 2 {
//		res[strs[i]] = strs[i+1]
//	}
//	return
//}

type OfflineMsg struct {
	Operation int32
	Data      []byte
}

func (r *logicRepo) PushOfflineMsg(ctx context.Context, op int32, mids []uint64, msg []byte) (err error) {
	offlineMsg := &OfflineMsg{
		Operation: op,
		Data:      msg,
	}
	_, err = r.data.redisCli.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, mid := range mids {
			pipe.ZAdd(ctx, keyMidOffline(mid), redis.Z{
				Score:  float64(time.Now().Unix()),
				Member: offlineMsg,
			})
		}
		return nil
	})
	return
}

func (r *logicRepo) ServersByKeys(ctx context.Context, keys []string) (res []string, err error) {
	args := make([]string, len(keys))
	for i, key := range keys {
		args[i] = keyKeyServer(key)
	}
	values, err := r.data.redisCli.MGet(ctx, args...).Result()
	if err != nil {
		return nil, err
	}
	for _, value := range values {
		switch v := value.(type) {
		case string:
			res = append(res, v)
		case []byte:
			res = append(res, string(v))
		default:
			return nil, fmt.Errorf("unexpected element type, got type %T", v)
		}
	}
	return
}

// KeysByMid 根据memberId获取key, 并可判断用户是否在线
func (r *logicRepo) KeysByMid(ctx context.Context, mid uint64) (res map[string]string, isOnline bool, err error) {
	res, err = r.data.redisCli.HGetAll(ctx, keyMidServer(mid)).Result()
	if err != nil {
		return nil, false, err
	}
	if len(res) > 0 {
		isOnline = true
	} else {
		isOnline = false
	}
	return
}

func (r *logicRepo) KeysByMids(ctx context.Context, mids []uint64) (ress map[string]string, onlineMids, offlineMids []uint64, err error) {
	cmds, err := r.data.redisCli.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, mid := range mids {
			pipe.HGetAll(ctx, keyMidServer(mid))
		}
		return nil
	})
	if err != nil {
		return
	}
	for i, cmd := range cmds {
		res := cmd.(*redis.MapStringStringCmd).Val()
		if len(res) > 0 {
			onlineMids = append(onlineMids, mids[i])
		} else {
			offlineMids = append(offlineMids, mids[i])
		}
		for k, v := range res {
			ress[k] = v
		}
	}
	return
}

func (r *logicRepo) AddServerOnline(ctx context.Context, server string, online *biz.Online) (err error) {
	roomsMap := map[uint32]map[string]int32{}
	for room, count := range online.RoomCount {
		hashKey := cityhash.CityHash32([]byte(room), uint32(len(room))) % 64
		rMap := roomsMap[hashKey]
		if rMap == nil {
			rMap = make(map[string]int32)
			roomsMap[hashKey] = rMap
		}
		rMap[room] = count
	}
	key := keyServerOnline(server)
	for hashKey, value := range roomsMap {
		err = r.addServerOnline(ctx, key, strconv.FormatInt(int64(hashKey), 10), &biz.Online{
			Server:    online.Server,
			RoomCount: value,
			Updated:   online.Updated,
		})
		if err != nil {
			return
		}
	}
	return
}

func (r *logicRepo) addServerOnline(ctx context.Context, key, hashKey string, online *biz.Online) (err error) {
	b, _ := json.Marshal(online)
	_, err = r.data.redisCli.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, key, hashKey, b)
		pipe.Expire(ctx, key, r.redisExpire)
		return nil
	})
	return
}

func (r *logicRepo) ServerOnline(ctx context.Context, server string) (online *biz.Online, err error) {
	online = &biz.Online{
		RoomCount: map[string]int32{},
	}
	key := keyServerOnline(server)
	for i := 0; i < 64; i++ {
		ol, err := r.serverOnline(ctx, key, strconv.FormatInt(int64(i), 10))
		if err == nil && ol != nil {
			online.Server = ol.Server
			if ol.Updated > online.Updated {
				online.Updated = ol.Updated
			}
			for room, count := range ol.RoomCount {
				online.RoomCount[room] = count
			}
		}
	}
	return
}

func (r *logicRepo) serverOnline(ctx context.Context, key, hashKey string) (online *biz.Online, err error) {
	b, err := r.data.redisCli.HGet(ctx, key, hashKey).Bytes()
	if err != nil {
		if err != redis.Nil {

		}
		return
	}
	online = new(biz.Online)
	if err = json.Unmarshal(b, online); err != nil {
		return
	}
	return
}

func (r *logicRepo) DelServerOnline(ctx context.Context, server string) (err error) {
	if _, err = r.data.redisCli.Del(ctx, keyServerOnline(server)).Result(); err != nil {

	}
	return
}

func (r *logicRepo) StoreChatMessage(ctx context.Context, seq uint32, chatMessage *protocol.ChatSendMessage) error {
	now := time.Now()
	cm := &db.ChatMessage{
		Base: database.Base{
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

	if err := r.data.db.Table("t_chat_message").Create(cm).Error; err != nil {
		return err
	}
	//r.msgCh <- cm

	return nil
}

func (r *logicRepo) storeBatchChatMessage(ctx context.Context, msgs []*db.ChatMessage) error {
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

func (r *logicRepo) StoreGroupChatMessage(ctx context.Context, seq uint32, groupChatMessage *protocol.GroupChatMessage) ([]uint64, error) {
	var err error
	now := time.Now()
	gcm := &db.GroupChatMessage{
		Base: database.Base{
			IsDeleted:   0,
			GMTCreate:   now,
			GMTModified: now,
		},
		GroupId:   groupChatMessage.GroupId,
		SenderId:  groupChatMessage.SenderId,
		MessageId: uint64(seq),
		Content:   string(groupChatMessage.Content),
	}

	var users []uint64
	err = r.data.db.Table("t_group_user").
		Select("user_id").
		Where("group_id = ?", groupChatMessage.GroupId).
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	n := len(users)
	oms := make([]*db.GroupChatOfflineMessage, 0, n-1)
	k := -1
	for i := 0; i < n; i++ {
		if users[i] != groupChatMessage.SenderId {
			oms = append(oms, &db.GroupChatOfflineMessage{
				GroupId:   groupChatMessage.GroupId,
				UserId:    users[i],
				MessageId: uint64(seq),
				GMTCreate: now,
			})
		} else {
			k = i
		}
	}
	if k >= 0 {
		users = append(users[:k], users[k+1:]...)
	}

	err = r.data.db.Transaction(func(tx *gorm.DB) error {
		var err1 error
		if err1 = r.data.db.Table("t_group_chat_message").Create(gcm).Error; err1 != nil {
			return err1
		}
		if err1 = r.data.db.Table("t_group_chat_offline_message").Create(oms).Error; err1 != nil {
			return err1
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *logicRepo) DeleteGroupChatOfflineMessage(ctx context.Context, groupId, userId, messageId uint64) error {
	return r.data.db.Table("t_group_chat_offline_message").
		Where("group_id = ? and user_id = ? and msg_id = ?", groupId, userId, messageId).
		Delete(&db.GroupChatOfflineMessage{}).Error
}

//func (r *logicRepo) PushChatMessageToKafka(ctx context.Context, chatMessage *protocol.ChatSendMessage) {
//	r.msgCh <- chatMessage
//}

//func (r *logicRepo) pushChatMessageToKafka(ctx context.Context, chatMessage *protocol.ChatSendMessage) (err error) {
//	b, err := proto.Marshal(chatMessage)
//	if err != nil {
//		return
//	}
//	m := &sarama.ProducerMessage{
//		Topic: "topic_chat_message",
//		//Key:   sarama.StringEncoder(keys[0]),
//		Value: sarama.ByteEncoder(b),
//	}
//	r.data.kafkaCli.SendMessageAsync(m)
//	return
//}

func (r *logicRepo) PushMessage(ctx context.Context, dm *protocol.DispatcherMessage) {
	r.msgCh <- dm
}

func (r *logicRepo) pushMessage(ctx context.Context, dm *protocol.DispatcherMessage) (err error) {
	b, err := proto.Marshal(dm)
	if err != nil {
		return
	}
	m := &sarama.ProducerMessage{
		Topic: "topic_dispatcher_message",
		Value: sarama.ByteEncoder(b),
	}
	r.data.kafkaCli.SendMessageAsync(m)
	return
}
