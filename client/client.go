package client

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/im-pusher/pkg/timingwheel"
	"github.com/zhihanii/retry"
	"github.com/zhihanii/zlog"
	"google.golang.org/protobuf/proto"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

const (
	connected int32 = iota
	disconnected
)

type Response struct {
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type NodeInstance struct {
	Addr string
}

type Client struct {
	ctx          context.Context
	cancel       context.CancelFunc
	httpClient   *http.Client
	nodeInstance *NodeInstance
	wsConn       *websocket.Conn
	send         chan *protocol.Message
	success      chan *protocol.Message
	fail         chan *protocol.Message

	state        int32
	closeSignal  chan struct{}
	readerSignal chan struct{}
	writerSignal chan struct{}
	errCh        chan error

	token []byte
	mid   uint64

	ackQueue       *queue
	notifyQueue    *queue
	ackMap         sync.Map
	notifyMap      sync.Map
	mackCh         chan []byte
	mackBuffer     *ackSet
	flushFrequency time.Duration
	timer          *timingwheel.TimingWheel
	retryer        retry.Retryer
	handler        Handler
}

func New(token []byte, mid uint64, flushFrequency time.Duration, addr string) *Client {
	c := &Client{
		httpClient:     &http.Client{},
		send:           make(chan *protocol.Message, 10),
		success:        make(chan *protocol.Message, 10),
		fail:           make(chan *protocol.Message, 10),
		state:          disconnected,
		closeSignal:    make(chan struct{}),
		readerSignal:   make(chan struct{}),
		writerSignal:   make(chan struct{}),
		errCh:          make(chan error, 10),
		token:          token,
		mid:            mid,
		mackCh:         make(chan []byte, 10),
		mackBuffer:     newAckSet(),
		flushFrequency: flushFrequency,
		retryer:        retry.NewRetryer(retry.WithMaxRetryTimes(5), retry.WithMaxDuration(time.Second*8)),
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.timer = timingwheel.NewTimingWheel(time.Second, 60)
	c.timer.Start()
	if addr != "" {
		c.nodeInstance = &NodeInstance{
			Addr: addr,
		}
	}
	return c
}

func (c *Client) Run() error {
	err := c.connect()
	if err != nil {
		zlog.Errorf("connect:%v", err)
		return err
	}
	go c.reader()
	go c.writer()
	go c.mackSender()

	for {
		select {
		case <-c.ctx.Done():
			return nil
		case <-c.closeSignal:
			err := c.connect()
			if err != nil {
				zlog.Errorf("connect:%v", err)
				return err
			}
			c.readerSignal <- struct{}{}
			c.writerSignal <- struct{}{}
		}
	}
}

func (c *Client) Errs() <-chan error {
	return c.errCh
}

func (c *Client) IsConnected() bool {
	return atomic.LoadInt32(&c.state) == connected
}

func (c *Client) setState(state int32) {
	atomic.StoreInt32(&c.state, state)
}

func (c *Client) setDisconnectedCAS() bool {
	return atomic.CompareAndSwapInt32(&c.state, connected, disconnected)
}

func (c *Client) Shutdown() {
	c.wsConn.Close()
	c.cancel()
}

func (c *Client) getNodeInstance() error {
	res, err := c.httpClient.Get("http://127.0.0.1:8202/node-instance")
	if err != nil {
		zlog.Errorf("http get node-instance:%v", err)
		return err
	}
	r := Response{}
	reader := bufio.NewReader(res.Body)
	data, _, err := reader.ReadLine()
	if err != nil {
		zlog.Errorf("read:%v", err)
		return err
	}
	err = json.Unmarshal(data, &r)
	if err != nil {
		zlog.Errorf("unmarshal:%v", err)
		return err
	}
	if r.Msg != "success" {
		return fmt.Errorf("获取connect-server地址失败:%s", r.Msg)
	}
	//zlog.Infof("成功获取connect-server的addr:%s", r.Data)

	c.nodeInstance = &NodeInstance{
		Addr: r.Data,
	}
	return nil
}

func (c *Client) connect() error {
	var err error
	if c.nodeInstance == nil {
		err = c.getNodeInstance()
		if err != nil {
			return err
		}
	}
	u := url.URL{Scheme: "ws", Host: c.nodeInstance.Addr, Path: "/ws"}
	c.wsConn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		zlog.Errorf("dial: %v\n", err)
		return err
	}
	//zlog.Infoln("成功与connect-server建立连接")
	return c.auth()
}

func (c *Client) auth() error {
	authMsg := &protocol.Message{
		Operation: protocol.Auth,
		Data:      c.token,
	}
	err := authMsg.WriteWebsocket(c.wsConn)
	if err != nil {
		zlog.Errorf("write auth: %v\n", err)
		return err
	} else {
		//zlog.Infoln("发送auth消息成功")
	}
	return nil
}

func (c *Client) Send(m *protocol.Message, ackKey string) {
	go func(seq uint32) {
		err := c.retryer.Do(c.ctx, func() error {
			c.send <- m
			fail := make(chan struct{})
			success := make(chan struct{})
			waitAckDuration := 5 * time.Second
			task := c.timer.Add(waitAckDuration, func() {
				fail <- struct{}{}
			})
			e := &element{
				sequence: seq,
				task:     task,
				success:  success,
			}
			c.ackMap.Store(ackKey, e)

			select {
			case <-success:
				task.Cancel()
				c.ackMap.Delete(ackKey)
				return nil
			case <-fail:
				c.ackMap.Delete(ackKey)
				return errors.New("ack timeout")
			}
		})
		if err != nil {
			zlog.Errorf("retry:%v", err)
		}
	}(m.Sequence)
}

func (c *Client) SendWithNotify(m *protocol.Message, ackKey string) {
	go func(seq uint32) {
		err := c.retryer.Do(c.ctx, func() error {
			c.send <- m
			fail := make(chan struct{})
			success := make(chan struct{})
			waitAckDuration := 5 * time.Second
			task := c.timer.Add(waitAckDuration, func() {
				fail <- struct{}{}
			})
			e := &element{
				sequence: seq,
				task:     task,
				success:  success,
			}
			c.ackMap.Store(ackKey, e)

			select {
			case <-success:
				task.Cancel()
				c.ackMap.Delete(ackKey)
				return nil
			case <-fail:
				c.ackMap.Delete(ackKey)
				return errors.New("ack timeout")
			}
		})
		if err != nil {

		}

		//todo 等待notify
		//消息已成功发送给对方, 但因服务器crash长时间未收到notify怎么办?
		fail := make(chan struct{})
		success := make(chan struct{})
		waitNotifyDuration := 10 * time.Second
		task := c.timer.Add(waitNotifyDuration, func() {
			fail <- struct{}{}
		})
		e := &element{
			sequence: seq,
			task:     task,
			success:  success,
		}
		c.notifyMap.Store(ackKey, e)

		select {
		case <-success:
			task.Cancel()
			c.notifyMap.Delete(ackKey)
			c.success <- m
			//zlog.Infof("mid:%d, 成功收到notify", c.mid)
		case <-fail:
			c.notifyMap.Delete(ackKey)
			c.fail <- m
		}
	}(m.Sequence)
}

func (c *Client) sendAck(data []byte) {
	ack := &protocol.Message{
		Operation: protocol.Ack,
		Data:      data,
	}
	c.send <- ack
}

func (c *Client) sendMAck(data []byte) {
	c.mackCh <- data
}

func (c *Client) sendAckNotify(data []byte) {
	ackNotify := &protocol.Message{
		Operation: protocol.AckNotify,
		Data:      data,
	}
	c.send <- ackNotify
}

func (c *Client) reader() {
	defer func() {
		c.wsConn.Close()
	}()

	hb := 5 * time.Second
	task := c.timer.Add(hb, func() {
		//若长时间未收到服务端的heartbeat, 则关闭连接
		zlog.Infof("close connection, reason:heartbeat timeout")
		//主动发送CloseMessage
		err1 := c.wsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err1 != nil {
			zlog.Errorf("close: %v", err1)
		}
		_ = c.wsConn.SetReadDeadline(time.Now().Add(time.Millisecond * 200))
		_ = c.wsConn.Close()
	})

read:
	for {
		_, msg, err := c.wsConn.ReadMessage()
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			c.wsConn.Close()
			if c.setDisconnectedCAS() {
				c.closeSignal <- struct{}{}
			}
			goto wait
		}
		if err != nil {
			zlog.Errorf("read: %v\n", err)
			return
		}
		m := protocol.Message{}
		err = m.FromBytes(msg)
		if err != nil {
			zlog.Errorf("消息解析失败: %v\n", err)
			return
		}
		switch m.Operation {
		case protocol.AuthReply:
			//zlog.Infoln("收到auth响应")

		case protocol.HeartbeatReply:
			//zlog.Infoln("收到Heartbeat响应")
			c.timer.Set(task, hb)
			m.Data = nil

		case protocol.Chat:
			chatMessage := new(protocol.ChatMessage)
			err = proto.Unmarshal(m.Data, chatMessage)
			if err != nil {
				zlog.Errorf("mid:%d, unmarshal:%v", c.mid, err)
			} else {
				//zlog.Infof("mid:%d, 收到chat消息, conversation_id:%d, from:%d, to:%d, content:%s",
				//	c.mid, chatMessage.ConversationId, chatMessage.From, chatMessage.To, string(chatMessage.Content))
			}
			//zlog.Infof("mid:%d, 发送ack", c.mid)
			c.sendAck([]byte(fmt.Sprintf("%d:%d:%d", chatMessage.From, chatMessage.To, m.Sequence)))

		case protocol.GroupChat:
			groupChatMessage := new(protocol.GroupChatMessage)
			err = proto.Unmarshal(m.Data, groupChatMessage)
			if err != nil {
				zlog.Errorf("mid:%d, unmarshal:%v", c.mid, err)
			} else {
				//zlog.Infof("mid:%d, 收到group chat消息, group_id:%d, sender_id:%d, content:%s",
				//	c.mid, groupChatMessage.GroupId, groupChatMessage.SenderId, string(groupChatMessage.Content))
			}
			c.sendMAck([]byte(fmt.Sprintf("%d:%d:%d", groupChatMessage.GroupId, groupChatMessage.SenderId, m.Sequence)))

		case protocol.Ack:
			//zlog.Infof("mid:%d, 收到ack消息", c.mid)
			key := string(m.Data)
			v, ok := c.ackMap.Load(key)
			if ok {
				e := v.(*element)
				e.success <- struct{}{}
				c.ackMap.Delete(key)
			}
		case protocol.Notify:
			//zlog.Infof("mid:%d, 收到notify消息", c.mid)
			key := string(m.Data)
			v, ok := c.notifyMap.Load(key)
			if ok {
				e := v.(*element)
				e.success <- struct{}{}
				c.notifyMap.Delete(key)
			}
			//zlog.Infof("mid:%d, 发送ack_notify", c.mid)
			c.sendAckNotify(m.Data)
		default:
			zlog.Errorf("unknown operation:%d", m.Operation)
		}
	}

wait:
	select {
	case <-c.ctx.Done():
		return
	case <-c.readerSignal:
		goto read
	}
}

func (c *Client) writer() {
	defer func() {
		c.wsConn.Close()
	}()

	var err error
	hbMsg := &protocol.Message{
		Operation: protocol.Heartbeat,
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

write:
	for {
		select {
		case m, ok := <-c.send:
			if !ok {
				return
			}

			err = m.WriteWebsocket(c.wsConn)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				c.wsConn.Close()
				if c.setDisconnectedCAS() {
					c.closeSignal <- struct{}{}
				}
				goto wait
			}
			if err != nil {

			}
		case <-ticker.C:
			//zlog.Infoln("发送Heartbeat")
			err = hbMsg.WriteWebsocket(c.wsConn)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				c.wsConn.Close()
				if c.setDisconnectedCAS() {
					c.closeSignal <- struct{}{}
				}
				goto wait
			}
			if err != nil {
				zlog.Errorf("write heartbeat:%v", err)
				return
			}
		}
	}

wait:
	select {
	case <-c.ctx.Done():
		return
	case <-c.writerSignal:
		goto write
	}
}

func (c *Client) mackSender() {
	ticker := time.NewTicker(c.flushFrequency)

	for {
		select {
		case ackKey, ok := <-c.mackCh:
			if !ok {
				//zlog.Infoln("mack chan closed")
				return
			}
			if err := c.mackBuffer.add(ackKey); err != nil {
				continue
			}
		case <-ticker.C:
			if !c.mackBuffer.empty() {
				mack := &protocol.Message{
					Operation: protocol.MAck,
					Data:      c.mackBuffer.bytes(),
				}
				//zlog.Infof("member:%d 发送mack", c.mid)
				c.send <- mack
			}
			c.mackBuffer.reset()
		}
	}
}
