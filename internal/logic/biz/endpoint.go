package biz

import (
	"context"
)

type Info struct {
	ReceiverId uint64
	Priority   int8
}

// Endpoint 若res为true, 则表示消息可以继续发送
type Endpoint func(ctx context.Context, info *Info) (res bool, err error)

type EndpointContainer struct {
	eps []Endpoint
}

func (c *EndpointContainer) Invoke(ctx context.Context, info *Info) (res bool, err error) {
	for _, ep := range c.eps {
		res, err = ep(ctx, info)
		if err != nil || !res {
			return
		}
	}
	return
}

//获取db对象来实现数据库操作

//// 检查用户是否开启推送
//func f1(ctx context.Context, memberId int64, metadata *MessageMetadata) (res bool, err error) {
//
//}
//
//// 检查用户是否收到重复的推送
//func f2(ctx context.Context, memberId int64, metadata *MessageMetadata) (res bool, err error) {
//
//}
//
//// 频率控制
//func f3(ctx context.Context, memberId int64, metadata *MessageMetadata) (res bool, err error) {
//
//}
//
//// 静默时间
//func f4(ctx context.Context, memberId int64, metadata *MessageMetadata) (res bool, err error) {
//
//}
//
//// 分级管理
//func f5(ctx context.Context, memberId int64, metadata *MessageMetadata) (res bool, err error) {
//
//}
