package data

import (
	"github.com/redis/go-redis/v9"
	"github.com/zhihanii/im-pusher/internal/logic/data/kafka"
	"gorm.io/gorm"
)

type Data struct {
	kafkaCli *kafka.Client
	redisCli redis.Cmdable
	db       *gorm.DB
}

func NewData(kafkaCli *kafka.Client, redisCli redis.Cmdable, db *gorm.DB) *Data {
	return &Data{
		kafkaCli: kafkaCli,
		redisCli: redisCli,
		db:       db,
	}
}
