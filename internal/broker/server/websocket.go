package server

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/hertz-contrib/websocket"
	"github.com/zhihanii/im-pusher/pkg/timingwheel"
	"github.com/zhihanii/zlog"
	"math/rand"
	"time"

	"github.com/zhihanii/im-pusher/api/protocol"
)

func (s *Server) initWebsocket() {
	go s.acceptWebsocket()
}

func (s *Server) acceptWebsocket() {
	//var wsUpgrader websocket.Upgrader
	//httpSrv := &http.Server{}
	//srvMux := http.NewServeMux()
	//srvMux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	//	conn, err := wsUpgrader.Upgrade(w, r, nil)
	//	if err != nil {
	//		return
	//	}
	//	go s.serveWebsocket(conn)
	//})
	//if err := httpSrv.ListenAndServe(); err != nil {
	//
	//}

	upgrader := websocket.HertzUpgrader{}
	h := server.New(server.WithHostPorts(fmt.Sprintf("%s:8082", s.c.InsecureServingOptions.BindAddress)),
		server.WithTransport(standard.NewTransporter))
	h.NoHijackConnPool = true
	h.GET("/ws", func(ctx context.Context, c *app.RequestContext) {
		zlog.Infoln("接收到websocket连接请求")
		//升级为Websocket协议, 并设置HertzHandler(会被封装为HijackHandler, 并在升级完成后执行)
		err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
			s.serveWebsocket(conn)
		})
		if err != nil {
			zlog.Errorf("upgrade:%v", err)
		}
	})
	zlog.Infoln("开始监听websocket连接")
	//err := h.Run()
	//if err != nil {
	//
	//}
	h.Spin()
}

func (s *Server) serveWebsocket(conn *websocket.Conn) {
	var (
		tr            = s.round.Timer()
		task          *timingwheel.Task
		m             *protocol.Message
		b             *Bucket
		ch            = NewChannel(s.c.Protocol.BufferSize, s.c.Protocol.MessageChannelSize)
		accepts       []int32
		hb            time.Duration
		lastHeartbeat time.Time //上一次调用logic层Heartbeat方法的时间
		roomId        string
		err           error
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if m, err = ch.Buffer.Set(); err == nil {
		//验证用户数据, 返回用户id、用户连接的唯一key、RoomId、代表用户订阅的Accepts
		if ch.MemberId, ch.Key, roomId, accepts, hb, err = s.authWebsocket(ctx, conn, m); err == nil {
			ch.Watch(accepts...)
			b = s.Bucket(ch.Key)
			//将用户Channel放入Room
			err = b.Put(roomId, ch)
		}
	}
	if err != nil {
		zlog.Infof("close conn")
		_ = conn.Close()
		return
	}

	//hb = 2 * time.Second
	//添加心跳监测任务
	task = tr.Add(hb, func() {
		//若长时间未收到客户端的heartbeat, 则关闭连接
		zlog.Infof("close connection, reason:heartbeat timeout, memberId:%d, key:%s", ch.MemberId, ch.Key)
		err1 := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err1 != nil {
			zlog.Errorf("close: %v", err)
		}
		_ = conn.SetReadDeadline(time.Now().Add(time.Millisecond * 200))
		_ = conn.Close()
	})

	//监听Channel中的消息, 并向用户推送消息
	go s.dispatchWebsocket(conn, ch)
	//两次调用Heartbeat方法的预期间隔
	serverHeartbeat := s.randHeartbeat()
	for {
		if m, err = ch.Buffer.Set(); err != nil {
			break
		}
		//从用户连接读取消息
		//收到的m都会被自动交由dispatchWebsocket处理
		if err = m.ReadWebsocket(conn); err != nil {
			//当客户端主动关闭连接时, 会返回err
			zlog.Infof("read: %v", err)
			break
		}
		upwardMessageCounter.Inc()
		if m.Operation == protocol.Heartbeat {
			//zlog.Infoln("收到Heartbeat消息")
			tr.Set(task, hb)
			m.Operation = protocol.HeartbeatReply
			m.Data = nil
			//两次调用Heartbeat方法的间隔时间较长
			//避免频繁调用Heartbeat方法
			if now := time.Now(); now.Sub(lastHeartbeat) > serverHeartbeat {
				if err1 := s.Heartbeat(ctx, ch.MemberId, ch.Key); err1 == nil {
					lastHeartbeat = now
				}
			}
		} else {
			//根据消息类型执行相应的操作
			var res OperateResult
			operateCounter.Inc()
			startTime := time.Now()
			if res = s.Operate(ctx, m, ch, b); res.Err != nil {
				break
			}
			endTime := time.Now()
			responseTime := endTime.Sub(startTime)
			zlog.Infof("operate response time:%d ms", responseTime.Milliseconds())
			//operateResponseTimeSummary.Observe(float64(responseTime))
			operateResponseTimeHistogram.Observe(float64(responseTime.Milliseconds()))
			//标记消息不需返回给发送者
			m.Operation = protocol.None
			m.Data = nil
			if res.Message != nil {
				zlog.Infof("msg operation:%d", res.Message.Operation)
				err = ch.Push(res.Message)
				if err != nil {
					zlog.Errorf("channel push:%v", err)
				}
			}
		}
		ch.Buffer.SetAdv()
		ch.Ready()
	}
	b.Del(ch)
	task.Cancel()
	conn.Close()
	ch.Close()

	//disconnect
}

// 监听Channel中的消息, 并向用户推送消息
func (s *Server) dispatchWebsocket(conn *websocket.Conn, ch *Channel) {
	var (
		err    error
		finish bool
		//online int32
	)
	for {
		var m = ch.Peek()
		switch m {
		case protocol.ProtoFinish:
			finish = true
			goto failed
		case protocol.ProtoReady: //Buffer中有来自客户端的消息写入
			for {
				//从Buffer中取出消息
				if m, err = ch.Buffer.Get(); err != nil {
					break
				}
				if m.Operation == protocol.HeartbeatReply {
					if ch.Room != nil {
						//online = ch.Room.OnlineNum()
					}
					//向客户端发送Heartbeat响应
					downwardMessageCounter.Inc()
					if err = m.WriteWebsocket(conn); err != nil {
						goto failed
					}
				} else if m.Operation != protocol.None {
					downwardMessageCounter.Inc()
					if err = m.WriteWebsocket(conn); err != nil {
						goto failed
					}
				}
				m.Data = nil
				ch.Buffer.GetAdv()
			}
		default:
			downwardMessageCounter.Inc()
			if err = m.WriteWebsocket(conn); err != nil {
				goto failed
			}
		}
	}
failed:
	//关闭连接
	conn.Close()
	for !finish {
		finish = ch.Peek() == protocol.ProtoFinish
	}
}

// 验证用户数据
func (s *Server) authWebsocket(ctx context.Context, conn *websocket.Conn,
	m *protocol.Message) (memberId uint64, key, roomId string, accepts []int32, hb time.Duration, err error) {
	for {
		//等待客户端发送用户验证数据
		if err = m.ReadWebsocket(conn); err != nil {
			return
		}
		if m.Operation == protocol.Auth {
			zlog.Infoln("收到验证消息")
			break
		} else {
			zlog.Infof("收到非验证消息:%d", m.Operation)
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
		zlog.Infoln("响应验证失败")
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
