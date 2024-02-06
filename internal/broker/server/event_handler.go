package server

import (
	"context"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/im-pusher/pkg/timingwheel"
	"github.com/zhihanii/websocket"
	"github.com/zhihanii/zlog"
	"github.com/zhihanii/znet"
	"time"
)

var _ znet.EventHandler = &Server{}

func (s *Server) OnConnect(ctx context.Context, conn znet.Conn) (err error) {
	var req *websocket.Request
	if req, err = websocket.ReadRequest(conn); err != nil {
		zlog.Errorf("websocket read request: %v", err)
		conn.Close()
		return err
	}

	var wsConn *websocket.Conn
	if wsConn, err = websocket.Upgrade(req, conn); err != nil {
		zlog.Errorf("websocket upgrade: %v", err)
		conn.Close()
		return err
	}

	var (
		m        = new(protocol.Message)
		memberId uint64
		key      string
		roomId   string
		accepts  []int32
		hb       time.Duration
		ch       *Channel
		bucket   *Bucket
		tr       *timingwheel.TimingWheel
		task     *timingwheel.Task
	)

	//验证用户数据, 返回用户id、用户连接的唯一key、RoomId、代表用户订阅的Accepts
	memberId, key, roomId, accepts, hb, err = s.authWebsocket(context.Background(), wsConn, m)
	if err != nil {
		zlog.Errorf("auth failed: %v", err)
		_ = conn.Close()
		return
	}
	ch = NewChannel(s.c.Protocol.BufferSize, s.c.Protocol.MessageChannelSize, wsConn)
	ch.MemberId = memberId
	ch.Key = key
	ch.Watch(accepts...)
	bucket = s.Bucket(ch.Key)
	bucket.Put(roomId, ch)
	tr = s.round.Timer()
	var hbDuration = s.randHeartbeat()

	//添加心跳监测任务
	task = tr.Add(hb, func() {
		//若长时间未收到客户端的heartbeat, 则关闭连接
		zlog.Infof("close connection, reason:heartbeat timeout, memberId:%d, key:%s", ch.MemberId, ch.Key)
		err1 := wsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err1 != nil {
			zlog.Errorf("close: %v", err1)
		}
		_ = conn.SetReadDeadline(time.Now().Add(time.Millisecond * 200))
		_ = conn.Close()
	})

	cc := &connContext{
		channel:           ch,
		bucket:            bucket,
		wsConn:            wsConn,
		tr:                tr,
		task:              task,
		heartbeat:         hb,
		heartbeatDuration: hbDuration,
		lastHeartbeat:     time.Time{},
	}

	conn.StoreValue(cc)

	return nil
}

func (s *Server) OnRead(ctx context.Context, conn znet.Conn) (err error) {
	var v = conn.LoadValue()
	if v == nil {
		return s.OnConnect(ctx, conn)
	}
	var cc *connContext
	cc, _ = v.(*connContext)
	var ch = cc.Channel()
	var bucket = cc.bucket
	var wsConn = cc.WebsocketConn()
	var tr = cc.Timer()
	var task = cc.TimerTask()
	var heartbeat = cc.Heartbeat()
	var hbDuration = cc.HeartbeatDuration()
	var lastHeartbeat = cc.LastHeartbeat()
	//todo 对象复用
	var m = new(protocol.Message)
	//从用户连接读取消息
	//收到的m都会被自动交由dispatchWebsocket处理
	if err = m.ReadWebsocket(wsConn); err != nil {
		//当客户端主动关闭连接时, 会返回err
		zlog.Errorf("read: %v", err)
		return
	}
	//upwardMessageCounter.Inc()
	if m.Operation == protocol.Heartbeat {
		//zlog.Infoln("收到Heartbeat消息")
		tr.Set(task, heartbeat)
		m.Operation = protocol.HeartbeatReply
		m.Data = nil
		s.Push(&protocol.TransMessage{
			Keys:      []string{ch.Key},
			Operation: protocol.HeartbeatReply,
		})
		//两次调用Heartbeat方法的间隔时间较长
		//避免频繁调用Heartbeat方法
		if now := time.Now(); now.Sub(lastHeartbeat) > hbDuration {
			if err1 := s.Heartbeat(ctx, ch.MemberId, ch.Key); err1 == nil {
				lastHeartbeat = now
			}
		}
	} else {
		//根据消息类型执行相应的操作
		var res OperateResult
		//operateCounter.Inc()
		//startTime := time.Now()
		if res = s.Operate(ctx, m, ch, bucket); res.Err != nil {
			//zlog.Errorf("operate result error: %v", res.Err)
		}
		//endTime := time.Now()
		//responseTime := endTime.Sub(startTime)
		//zlog.Infof("operate response time:%d ms", responseTime.Milliseconds())
		//operateResponseTimeSummary.Observe(float64(responseTime))
		//operateResponseTimeHistogram.Observe(float64(responseTime.Milliseconds()))
		//标记消息不需返回给发送者
		//m.Operation = protocol.None
		m.Data = nil
		if res.Message != nil {
			//zlog.Infof("msg operation:%d", res.Message.Operation)
			s.Push(&protocol.TransMessage{
				Keys:      []string{ch.Key},
				Operation: res.Message.Operation,
				Sequence:  res.Message.Sequence,
				Data:      res.Message.Data,
			})
		}
	}

	return nil
}
