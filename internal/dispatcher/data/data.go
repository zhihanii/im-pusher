package data

import (
	"github.com/redis/go-redis/v9"
)

type Data struct {
	redisCli redis.Cmdable
}

func NewData(redisCli redis.Cmdable) *Data {
	return &Data{
		redisCli: redisCli,
	}
}
