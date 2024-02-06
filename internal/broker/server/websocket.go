package server

import (
	"context"
	"github.com/zhihanii/websocket"
	"github.com/zhihanii/zlog"
	"math/rand"
	"time"

	"github.com/zhihanii/im-pusher/api/protocol"
)

// 验证用户数据
func (s *Server) authWebsocket(ctx context.Context, conn *websocket.Conn,
	m *protocol.Message) (memberId uint64, key, roomId string, accepts []int32, hb time.Duration, err error) {
	for {
		//等待客户端发送用户验证数据
		if err = m.ReadWebsocket(conn); err != nil {
			return
		}
		if m.Operation == protocol.Auth {
			//zlog.Infoln("收到验证消息")
			break
		} else {
			//zlog.Infof("收到非验证消息:%d", m.Operation)
		}
	}
	//用户验证数据交由logic层解析, 同时会在logic层记录memberId与Connect Server的映射关系
	if memberId, key, roomId, accepts, hb, err = s.Connect(ctx, m); err != nil {
		zlog.Errorf("验证失败:%v", err)
		return
	}
	m.Operation = protocol.AuthReply
	m.Data = nil
	//响应验证
	if err = m.WriteWebsocket(conn); err != nil {
		zlog.Errorf("响应验证失败: %v", err)
		return
	}
	return
}

const (
	minServerHeartbeat = time.Second * 10
	maxServerHeartbeat = time.Second * 30
)

func (s *Server) randHeartbeat() time.Duration {
	return minServerHeartbeat + time.Duration(rand.Int63n(int64(maxServerHeartbeat-minServerHeartbeat)))
}
