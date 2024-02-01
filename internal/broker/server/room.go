package server

import (
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/im-pusher/internal/broker/errors"
	"sync"
)

type Room struct {
	Id string
	sync.RWMutex
	next      *Channel
	drop      bool //房间的在线人数是否为0
	Online    int32
	AllOnline int32
}

func NewRoom(id string) *Room {
	return &Room{
		Id:     id,
		drop:   false,
		next:   nil,
		Online: 0,
	}
}

// Put 将Channel加入Room
func (r *Room) Put(ch *Channel) (err error) {
	r.Lock()
	if !r.drop {
		if r.next != nil {
			r.next.Prev = ch //头插法
		}
		ch.Next = r.next
		ch.Prev = nil
		r.next = ch
		r.Online++
	} else {
		err = errors.ErrRoomDroped
	}
	r.Unlock()
	return
}

// Del 将Channel移出Room, 若Room的在线人数为0, 则返回true
func (r *Room) Del(ch *Channel) bool {
	r.Lock()
	if ch.Next != nil {
		ch.Next.Prev = ch.Prev
	}
	if ch.Prev != nil {
		ch.Prev.Next = ch.Next
	} else {
		r.next = ch.Next
	}
	ch.Next = nil
	ch.Prev = nil
	r.Online--
	r.drop = r.Online == 0
	r.Unlock()
	return r.drop
}

// Push 放入消息
func (r *Room) Push(m *protocol.Message) {
	r.RLock()
	//给Room中的各个Channel发送消息
	for ch := r.next; ch != nil; ch = ch.Next {
		_ = ch.Push(m)
	}
	r.RUnlock()
}

func (r *Room) Close() {
	r.RLock()
	//关闭所有Channel
	for ch := r.next; ch != nil; ch = ch.Next {
		ch.Close()
	}
	r.RUnlock()
}

func (r *Room) OnlineNum() int32 {
	if r.AllOnline > 0 {
		return r.AllOnline
	}
	return r.Online
}
