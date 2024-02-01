package data

import (
	"github.com/redis/go-redis/v9"
	"github.com/zhihanii/im-pusher/internal/dispatcher/data/kafka"
)

type Data struct {
	kafkaCli *kafka.Client
	redisCli redis.Cmdable
}

func NewData(kafkaCli *kafka.Client, redisCli redis.Cmdable) *Data {
	return &Data{
		kafkaCli: kafkaCli,
		redisCli: redisCli,
	}
}
