package data

import (
	"context"
	"fmt"
	"github.com/zhihanii/im-pusher/internal/dispatcher/conf"
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

type DispatcherRepo interface {
	GetMapping(ctx context.Context, memberId uint64) (res map[string]string, err error)
}

type dispatcherRepo struct {
	c    *conf.Config
	data *Data
}

func NewDispatcherRepo(c *conf.Config, data *Data) DispatcherRepo {
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
	//r.data.redisCli.Get(ctx, "123")
	return
}
