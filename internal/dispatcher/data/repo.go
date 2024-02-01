package data

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/im-pusher/internal/dispatcher/biz"
	"github.com/zhihanii/im-pusher/internal/dispatcher/conf"
	"google.golang.org/protobuf/proto"
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

type dispatcherRepo struct {
	c    *conf.Config
	data *Data
}

func NewDispatcherRepo(c *conf.Config, data *Data) biz.DispatcherRepo {
	return &dispatcherRepo{
		c:    c,
		data: data,
	}
}

func (r *dispatcherRepo) GetMapping(ctx context.Context, memberId uint64) (res map[string]string, err error) {
	res = make(map[string]string)
	strs, err := getMappingScript.Run(ctx, r.data.redisCli,
		[]string{keyMidServer(memberId)}).StringSlice()
	if err != nil {
		return nil, err
	}
	n := len(strs)
	if n > 0 && n%2 == 1 {
		//log
	}
	for i := 0; i < n-1; i += 2 {
		res[strs[i]] = strs[i+1]
	}
	return
}

func (r *dispatcherRepo) PushMessage(ctx context.Context, keys []string, msg *protocol.TransMessage) (err error) {
	b, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	m := &sarama.ProducerMessage{
		Topic: "my_topic1",
		Key:   sarama.StringEncoder(keys[0]),
		Value: sarama.ByteEncoder(b),
	}
	r.data.kafkaCli.SendMessageAsync(m)
	return
}
