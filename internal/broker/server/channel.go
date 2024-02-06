package server

import (
	"bufio"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/websocket"
	"sync"
)

// Channel 一个Channel对应一个websocket连接
type Channel struct {
	Room   *Room
	Buffer *CircularBuffer        //用于存放来自客户端的消息
	msgCh  chan *protocol.Message //需要推送给客户端的消息
	Writer bufio.Writer
	Reader bufio.Reader
	Next   *Channel
	Prev   *Channel

	MemberId uint64
	Key      string
	IP       string
	watchOps map[int32]struct{}
	sync.RWMutex

	groupMux      sync.RWMutex
	watchGroupIDs map[uint64]struct{}

	conn *websocket.Conn
}

func NewChannel(bufferSize, messageChannelSize int, conn *websocket.Conn) *Channel {
	c := &Channel{
		Room:   NewRoom("some"),
		Buffer: NewCircularBuffer(16),
		conn:   conn,
	}
	c.Buffer.Init(bufferSize)
	c.msgCh = make(chan *protocol.Message, messageChannelSize)
	c.watchOps = make(map[int32]struct{})
	return c
}

func (c *Channel) Watch(accepts ...int32) {
	c.Lock()
	for _, op := range accepts {
		c.watchOps[op] = struct{}{}
	}
	c.Unlock()
}

func (c *Channel) UnWatch(accepts ...int32) {
	c.Lock()
	for _, op := range accepts {
		delete(c.watchOps, op)
	}
	c.Unlock()
}

// NeedPush Channel是否Watch了op
func (c *Channel) NeedPush(op int32) bool {
	c.RLock()
	if _, ok := c.watchOps[op]; ok {
		c.RUnlock()
		return true
	}
	c.RUnlock()
	return false
}

func (c *Channel) NeedPushGroup(groupID uint64) bool {
	c.groupMux.RLock()
	if _, ok := c.watchGroupIDs[groupID]; ok {
		c.groupMux.RUnlock()
		return true
	}
	c.groupMux.RUnlock()
	return false
}

// Push 放入消息
func (c *Channel) Push(m *protocol.Message) (err error) {
	//select {
	//case c.msgCh <- m:
	//default:
	//	err = errors.ErrSignalFullMsgDropped
	//}

	err = c.write(m)
	return
}

func (c *Channel) write(m *protocol.Message) (err error) {
	//downwardMessageCounter.Inc()
	return m.WriteWebsocket(c.conn)
}

// Peek 取出消息
func (c *Channel) Peek() *protocol.Message {
	return <-c.msgCh
}

// Ready 通知Buffer中已有消息写入
func (c *Channel) Ready() {
	c.msgCh <- protocol.ProtoReady
}

func (c *Channel) Close() {
	c.msgCh <- protocol.ProtoFinish
}
