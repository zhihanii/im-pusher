package server

import (
	"github.com/zhihanii/im-pusher/pkg/timingwheel"
	"github.com/zhihanii/websocket"
	"time"
)

type connContext struct {
	channel           *Channel
	bucket            *Bucket
	wsConn            *websocket.Conn
	tr                *timingwheel.TimingWheel
	task              *timingwheel.Task
	heartbeat         time.Duration
	heartbeatDuration time.Duration
	lastHeartbeat     time.Time
}

func (c *connContext) Channel() *Channel {
	return c.channel
}

func (c *connContext) Bucket() *Bucket {
	return c.bucket
}

func (c *connContext) WebsocketConn() *websocket.Conn {
	return c.wsConn
}

func (c *connContext) Timer() *timingwheel.TimingWheel {
	return c.tr
}

func (c *connContext) TimerTask() *timingwheel.Task {
	return c.task
}

func (c *connContext) Heartbeat() time.Duration {
	return c.heartbeat
}

func (c *connContext) HeartbeatDuration() time.Duration {
	return c.heartbeatDuration
}

func (c *connContext) LastHeartbeat() time.Time {
	return c.lastHeartbeat
}
