package server

import (
	pb "github.com/zhihanii/im-pusher/api/broker"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/im-pusher/internal/broker/conf"
	"sync"
	"sync/atomic"
)

type Bucket struct {
	c *conf.Config

	sync.RWMutex
	chs         map[string]*Channel
	rooms       map[string]*Room
	routines    []chan *pb.BroadcastRoomReq
	routinesNum uint64

	ipCnts map[string]int32
}

func NewBucket(c *conf.Config) *Bucket {
	b := &Bucket{
		c:        c,
		chs:      make(map[string]*Channel),
		rooms:    make(map[string]*Room),
		routines: make([]chan *pb.BroadcastRoomReq, 10),
		ipCnts:   make(map[string]int32),
	}
	return b
}

// ChannelCount Channel数量
func (b *Bucket) ChannelCount() int {
	return len(b.chs)
}

// RoomCount Room数量
func (b *Bucket) RoomCount() int {
	return len(b.rooms)
}

// RoomsCount 统计每个Room的在线人数
func (b *Bucket) RoomsCount() (res map[string]int32) {
	var (
		roomId string
		room   *Room
	)
	b.RLock()
	res = make(map[string]int32)
	for roomId, room = range b.rooms {
		if room.Online > 0 {
			res[roomId] = room.Online
		}
	}
	b.RUnlock()
	return
}

// ChangeRoom 改变Channel的Room
func (b *Bucket) ChangeRoom(newRoomId string, ch *Channel) (err error) {
	var (
		newRoom *Room
		ok      bool
		oldRoom = ch.Room
	)
	if newRoomId == "" {
		//从oldRoom删除ch, 若oldRoom在线人数为0, 则删除oldRoom
		if oldRoom != nil && oldRoom.Del(ch) {
			b.DelRoom(oldRoom)
		}
		ch.Room = nil
		return
	}
	b.Lock()
	if newRoom, ok = b.rooms[newRoomId]; !ok {
		newRoom = NewRoom(newRoomId)
		b.rooms[newRoomId] = newRoom
	}
	b.Unlock()
	//退出原来的Room
	if oldRoom != nil && oldRoom.Del(ch) {
		b.DelRoom(oldRoom)
	}
	//加入新的Room
	if err = newRoom.Put(ch); err != nil {
		return
	}
	ch.Room = newRoom
	return
}

// Put 将Channel加入指定的Room
func (b *Bucket) Put(rid string, ch *Channel) (err error) {
	var (
		room *Room
		ok   bool
	)
	b.Lock()
	if dch := b.chs[ch.Key]; dch != nil {
		dch.Close() //关闭原来存在的Channel
	}
	b.chs[ch.Key] = ch //记录Channel
	if rid != "" {
		if room, ok = b.rooms[rid]; !ok {
			room = NewRoom(rid)
			b.rooms[rid] = room
		}
		ch.Room = room //加入Room
	}
	b.ipCnts[ch.IP]++ //IP数+1
	b.Unlock()
	if room != nil {
		err = room.Put(ch) //加入Room
	}
	return
}

// Del 删除Channel
func (b *Bucket) Del(dch *Channel) {
	room := dch.Room
	b.Lock()
	if ch, ok := b.chs[dch.Key]; ok {
		if ch == dch {
			delete(b.chs, ch.Key)
		}
		if b.ipCnts[ch.IP] > 1 { //同一个IP存在多个会话
			b.ipCnts[ch.IP]--
		} else {
			delete(b.ipCnts, ch.IP)
		}
	}
	b.Unlock()
	//从Room移出Channel
	if room != nil && room.Del(dch) {
		b.DelRoom(room)
	}
}

// Channel 获取Channel
func (b *Bucket) Channel(key string) (ch *Channel) {
	b.RLock()
	ch = b.chs[key]
	b.RUnlock()
	return
}

// Broadcast 向所有Watch了op的Channel广播消息
func (b *Bucket) Broadcast(m *protocol.Message, op int32) {
	var ch *Channel
	b.RLock()
	for _, ch = range b.chs {
		if !ch.NeedPush(op) {
			continue
		}
		_ = ch.Push(m)
	}
	b.RUnlock()
}

// Room 获取Room
func (b *Bucket) Room(rid string) (room *Room) {
	b.RLock()
	room = b.rooms[rid]
	b.RUnlock()
	return
}

// DelRoom 删除Room
func (b *Bucket) DelRoom(room *Room) {
	b.Lock()
	delete(b.rooms, room.Id)
	b.Unlock()
	room.Close()
}

func (b *Bucket) BroadcastRoom(req *pb.BroadcastRoomReq) {
	//num := atomic.AddUint64(&b.routinesNum, 1) % b.c.RoutineAmount
	num := atomic.AddUint64(&b.routinesNum, 1) % 10
	b.routines[num] <- req
}

// Rooms 返回所有在线人数>0的RoomId
func (b *Bucket) Rooms() (res map[string]struct{}) {
	var (
		roomId string
		room   *Room
	)
	res = make(map[string]struct{})
	b.RLock()
	for roomId, room = range b.rooms {
		if room.Online > 0 {
			res[roomId] = struct{}{}
		}
	}
	b.RUnlock()
	return
}

// IPCount 返回所有IP
func (b *Bucket) IPCount() (res map[string]struct{}) {
	var ip string
	b.RLock()
	res = make(map[string]struct{}, len(b.ipCnts))
	for ip = range b.ipCnts {
		res[ip] = struct{}{}
	}
	b.RUnlock()
	return
}

// UpdateRoomsCount 更新Room的在线人数
func (b *Bucket) UpdateRoomsCount(roomCntMap map[string]int32) {
	var (
		roomId string
		room   *Room
	)
	b.RLock()
	for roomId, room = range b.rooms {
		room.AllOnline = roomCntMap[roomId]
	}
	b.RUnlock()
}

func (b *Bucket) roomProc(c chan *pb.BroadcastRoomReq) {
	for {
		req := <-c
		if room := b.Room(req.RoomId); room != nil {
			room.Push(req.Msg)
		}
	}
}
