package biz

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhihanii/im-pusher/api/dispatcher"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/im-pusher/internal/logic/rpc"
	"github.com/zhihanii/im-pusher/pkg/gopool"
	"github.com/zhihanii/im-pusher/pkg/set"
	"github.com/zhihanii/im-pusher/pkg/timingwheel"
	"github.com/zhihanii/loadbalance"
	"github.com/zhihanii/retry"
	"github.com/zhihanii/zlog"
	"google.golang.org/protobuf/proto"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/zhihanii/im-pusher/internal/logic/conf"
)

type LogicHandler struct {
	c *conf.Config

	f *loadbalance.WatchBalancerFactory

	eps *EndpointContainer

	repo LogicRepo
	// online
	totalIPs   int64
	totalConns int64
	roomCount  map[string]int32
	regions    map[string]string // province -> region

	ackMap       sync.Map
	ackNotifyMap sync.Map
	groupAckMap  sync.Map
	timer        *timingwheel.TimingWheel
	retryer      retry.Retryer
}

func NewLogicHandler(c *conf.Config, f *loadbalance.WatchBalancerFactory, repo LogicRepo) (*LogicHandler, error) {
	lh := &LogicHandler{
		c:       c,
		f:       f,
		repo:    repo,
		retryer: retry.NewRetryer(retry.WithMaxRetryTimes(5), retry.WithMaxDuration(time.Second*8)),
	}
	lh.timer = timingwheel.NewTimingWheel(time.Second, 60)
	lh.timer.Start()
	return lh, nil
}

func (h *LogicHandler) NodeInstance(ctx context.Context) (string, error) {
	balancer, err := h.f.Get(ctx, "connect-service")
	if err != nil {
		return "", err
	}

	picker := balancer.GetPicker()
	ins := picker.Next(ctx, nil)
	if ins == nil {
		return "", fmt.Errorf("no connect-server instances are available")
	}
	zlog.Infof("next instance address:%s", ins.Address())

	return fmt.Sprintf("%s:8082", ins.Address()), nil
}

// Connect 验证用户数据
func (h *LogicHandler) Connect(ctx context.Context, server string, token []byte) (memberId uint64, key string,
	roomId string, accepts []int32, hb int64, err error) {
	var params struct {
		MemberId uint64 `json:"member_id"`
		Key      string `json:"key"`
		//RoomId   string  `json:"room_id"`
		//Accepts  []int32 `json:"accepts"`
	}
	//todo rpc调用鉴权服务验证token
	//解析token
	if err = json.Unmarshal(token, &params); err != nil {
		return
	}
	memberId = params.MemberId
	hb = int64(time.Second * 10)
	if key = params.Key; key == "" {
		key = uuid.New().String()
	}
	//记录映射关系
	if err = h.repo.AddMapping(ctx, memberId, key, server); err != nil {
		zlog.Errorf("add mapping:%v", err)
		return
	}
	return
}

func (h *LogicHandler) Disconnect(ctx context.Context, memberId uint64, key, server string) (has bool, err error) {
	if has, err = h.repo.DelMapping(ctx, memberId, key, server); err != nil {
		zlog.Errorf("del mapping:%v", err)
		return
	}
	return
}

func (h *LogicHandler) Heartbeat(ctx context.Context, memberId uint64, key, server string) (err error) {
	//尝试延长key的expire
	has, err := h.repo.ExpireMapping(ctx, memberId, key)
	if err != nil {
		zlog.Errorf("expire mapping:%v", err)
		return
	}
	if !has {
		//若key已被删除, 则重新添加key
		if err = h.repo.AddMapping(ctx, memberId, key, server); err != nil {
			zlog.Errorf("add mapping:%v", err)
			return
		}
	}
	return
}

func (h *LogicHandler) RenewOnline(ctx context.Context, server string, roomCount map[string]int32) (map[string]int32, error) {
	online := &Online{
		Server:    server,
		RoomCount: roomCount,
		Updated:   time.Now().Unix(),
	}
	if err := h.repo.AddServerOnline(ctx, server, online); err != nil {
		return nil, err
	}
	return h.roomCount, nil
}

func (h *LogicHandler) location(ctx context.Context, clientIP string) (province string, err error) {
	// province: config mapping
	return
}

func (h *LogicHandler) Receive(ctx context.Context, memberId uint64, m *protocol.Message) (ackMsg *protocol.Message, err error) {
	switch m.Operation {
	case protocol.Chat:
		return h.HandleChat(ctx, m)

	case protocol.GroupChat:
		return h.HandleGroupChat(ctx, memberId, m)

	case protocol.Ack:
		//zlog.Infof("ack from:%d", memberId)
		ackKey := string(m.Data)
		v, ok := h.ackMap.Load(ackKey)
		if ok {
			e := v.(*element)
			e.success <- struct{}{}
			h.ackMap.Delete(ackKey)
		}
		return

	case protocol.MAck:
		zlog.Infof("收到来自member:%d的mack", memberId)
		keys := bytes.Split(m.Data, []byte{';'})
		for _, k := range keys {
			v, ok := h.groupAckMap.Load(string(k))
			if ok {
				s := v.(*set.Set)
				s.Remove(memberId)
			}
			strs := strings.Split(string(k), ":")
			if len(strs) != 3 {
				zlog.Errorf("length is not equal to 3")
			} else {
				groupId, _ := strconv.Atoi(strs[0])
				messageId, _ := strconv.Atoi(strs[2])
				zlog.Infof("删除member:%d离线消息", memberId)
				err1 := h.repo.DeleteGroupChatOfflineMessage(ctx, uint64(groupId), memberId, uint64(messageId))
				if err1 != nil {
					zlog.Errorf("delete group chat offline message:%v", err1)
				}
			}
		}
		return

	case protocol.AckNotify:
		//zlog.Infof("ack_notify from:%d", memberId)
		ackKey := string(m.Data)
		v, ok := h.ackNotifyMap.Load(ackKey)
		if ok {
			e := v.(*element)
			e.success <- struct{}{}
			h.ackNotifyMap.Delete(ackKey)
		}
		return

	default:
		return nil, fmt.Errorf("unknown message operation:%d", m.Operation)
	}
}

func (h *LogicHandler) HandleChat(ctx context.Context, m *protocol.Message) (ackMsg *protocol.Message, err error) {
	chatMessage := new(protocol.ChatMessage)
	err = proto.Unmarshal(m.Data, chatMessage)
	if err != nil {
		zlog.Errorf("unmarshal:%v", err)
		return
	}

	//startTime := time.Now()
	//err = h.repo.StoreChatMessage(ctx, m.Sequence, chatMessage)
	//if err != nil {
	//	zlog.Errorf("store chat message:%v", err)
	//	return
	//}

	//err = h.repo.PushChatMessageToKafka(ctx, chatMessage)
	//if err != nil {
	//	zlog.Errorf("push chat message to kafka:%v", err)
	//	return
	//}

	//endTime := time.Now()
	//responseTime := endTime.Sub(startTime)
	//zlog.Infof("store chat message response time:%d ms", responseTime.Milliseconds())

	//异步推送消息
	ctx1 := context.Background()
	st1 := time.Now()
	submitErr := gopool.Submit(ctx1, func() {
		st2 := time.Now()
		h.repo.PushChatMessageToKafka(ctx, chatMessage)
		et2 := time.Now()
		rt2 := et2.Sub(st2)
		zlog.Infof("push to kafka响应时间:%d ms", rt2.Milliseconds())

		st := time.Now()
		//zlog.Infoln("正在异步推送消息")
		ackKey := key(chatMessage.From, chatMessage.To, m.Sequence)
		err1 := h.retryer.Do(ctx1, func() error {
			_, err2 := rpc.Push(ctx1, &dispatcher.PushReq{
				Receivers: []uint64{chatMessage.To},
				Message:   m,
			})
			if err2 != nil {
				zlog.Errorf("dispatcher client push:%v", err2)
				return err2
			}
			fail := make(chan struct{})
			success := make(chan struct{})
			waitAckDuration := 5 * time.Second
			task := h.timer.Add(waitAckDuration, func() {
				fail <- struct{}{}
			})
			e := &element{
				task:    task,
				success: success,
			}
			h.ackMap.Store(ackKey, e)
			select {
			case <-success:
				task.Cancel()
				h.ackMap.Delete(ackKey)
				return nil
			case <-fail:
				h.ackMap.Delete(ackKey)
				return errors.New("ack timeout")
			}
		})
		if err1 != nil {
			zlog.Errorf("retry do:%v", err1)
			return
		}

		err1 = h.retryer.Do(ctx1, func() error {
			notifyMessage := &protocol.Message{
				Operation: protocol.Notify,
				Data:      []byte(ackKey),
			}
			_, err2 := rpc.Push(ctx1, &dispatcher.PushReq{
				Receivers: []uint64{chatMessage.From},
				Message:   notifyMessage,
			})
			if err2 != nil {
				zlog.Errorf("dispatcher client push:%v", err2)
				return err2
			}
			fail := make(chan struct{})
			success := make(chan struct{})
			waitAckDuration := 5 * time.Second
			task := h.timer.Add(waitAckDuration, func() {
				fail <- struct{}{}
			})
			e := &element{
				task:    task,
				success: success,
			}
			h.ackNotifyMap.Store(ackKey, e)
			select {
			case <-success:
				task.Cancel()
				h.ackNotifyMap.Delete(ackKey)
				return nil
			case <-fail:
				h.ackNotifyMap.Delete(ackKey)
				return errors.New("ack_notify timeout")
			}
		})
		if err1 != nil {
			zlog.Errorf("retry do:%v", err1)
			return
		}
		et := time.Now()
		rt := et.Sub(st)
		zlog.Infof("异步消息推送所用时间:%d ms", rt.Milliseconds())
	})
	if submitErr != nil {
		if errors.Is(submitErr, gopool.ErrPoolIsOverload) {
			zlog.Infof("submit:go pool is overload")
		}
	}
	et1 := time.Now()
	rt1 := et1.Sub(st1)
	zlog.Infof("submit response time:%d ms", rt1.Milliseconds())

	ackMsg = &protocol.Message{
		Operation: protocol.Ack,
		Data:      []byte(fmt.Sprintf("%d:%d:%d", chatMessage.From, chatMessage.To, m.Sequence)),
	}

	return
}

func (h *LogicHandler) HandleGroupChat(ctx context.Context, memberId uint64, m *protocol.Message) (ackMsg *protocol.Message, err error) {
	groupChatMessage := new(protocol.GroupChatMessage)
	err = proto.Unmarshal(m.Data, groupChatMessage)
	if err != nil {
		zlog.Errorf("unmarshal:%v", err)
		return
	}

	users, err := h.repo.StoreGroupChatMessage(ctx, m.Sequence, groupChatMessage)
	if err != nil {
		zlog.Errorf("store group chat message:%v", err)
		return
	}

	ackKey := fmt.Sprintf("%d:%d:%d", groupChatMessage.GroupId, groupChatMessage.SenderId, m.Sequence)

	if len(users) == 0 {
		ackMsg = &protocol.Message{
			Operation: protocol.Ack,
			Data:      []byte(ackKey),
		}
		return
	}

	//异步推送消息
	ctx1 := context.Background()
	s := set.New(users...)
	h.groupAckMap.Store(ackKey, s)
	leftUsers := users
	submitErr := gopool.Submit(ctx1, func() {
		err1 := h.retryer.Do(ctx1, func() error {
			offlineMembers, err2 := rpc.Push(ctx1, &dispatcher.PushReq{
				Receivers: leftUsers,
				Message:   m,
			})
			if err2 != nil {
				zlog.Errorf("dispatch client push:%v", err2)
				return err2
			}

			if len(offlineMembers) > 0 {
				s.Remove(offlineMembers...)
			}

			waitAckDuration := 30 * time.Second
			select {
			case <-time.After(waitAckDuration):
				leftUsers = s.Values()
				if len(leftUsers) > 0 {
					return errors.New("ack_notify timeout")
				} else {
					zlog.Infof("在线用户的ack已完成, ackKey:%s", ackKey)
					h.groupAckMap.Delete(ackKey)
					return nil
				}
			}
		})
		if err1 != nil {
			zlog.Errorf("retry do:%v", err1)
			return
		}
	})
	if submitErr != nil {
		if errors.Is(submitErr, gopool.ErrPoolIsOverload) {
			zlog.Infof("submit:go pool is overload")
		}
	}

	ackMsg = &protocol.Message{
		Operation: protocol.Ack,
		Data:      []byte(ackKey),
	}
	return
}

func key(mid1, mid2 uint64, seq uint32) string {
	return fmt.Sprintf("%d:%d:%d", mid1, mid2, seq)
}

type element struct {
	task    *timingwheel.Task
	success chan struct{}
}
