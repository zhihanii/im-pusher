package server

import (
	"context"
	pb "github.com/zhihanii/im-pusher/api/broker"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/im-pusher/internal/dispatcher/conf"
	"github.com/zhihanii/registry"
	"github.com/zhihanii/zlog"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync/atomic"
	"time"
)

type Broker struct {
	c *conf.Config

	serverId    string
	cli         pb.BrokerClient
	pushCh      []chan *pb.PushMsgReq
	roomCh      []chan *pb.BroadcastRoomReq
	broadcastCh chan *pb.BroadcastReq
	pushChNum   uint64
	roomChNum   uint64
	routineSize uint64

	ctx    context.Context
	cancel context.CancelFunc

	input  chan *protocol.TransMessage
	output chan *messageSet

	buffer     *messageSet
	timer      *time.Timer
	timerFired bool
}

func newBroker(config *conf.Config, serverId string, etcdClient *clientv3.Client) *Broker {
	b := &Broker{
		c:        config,
		serverId: serverId,
	}
	var err error
	etcdResolver, err := registry.NewGRPCEtcdResolverBuilder(etcdClient)
	if err != nil {
		zlog.Errorf("resolver new builder:%v", err)
	}
	conn, err := grpc.Dial("etcd:///broker-service", grpc.WithResolvers(etcdResolver),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	zlog.Infoln("成功与broker建立grpc连接")
	b.cli = pb.NewBrokerClient(conn)
	b.pushCh = make([]chan *pb.PushMsgReq, 1)
	b.pushChNum = 1
	b.ctx, b.cancel = context.WithCancel(context.Background())
	b.input = make(chan *protocol.TransMessage)
	b.output = make(chan *messageSet)
	b.buffer = newMessageSet(b.c)
	go b.run()
	go b.sender()
	return b
}

func (b *Broker) Push(msg *protocol.TransMessage) {
	b.input <- msg
}

func (b *Broker) run() {
	var output chan<- *messageSet
	var timerChan <-chan time.Time

	for {
		select {
		case msg, ok := <-b.input:
			if !ok {
				zlog.Infoln("input chan closed")
				return
			}

			if msg == nil {
				continue
			}

			//todo if needs retry

			if err := b.buffer.add(msg); err != nil {
				continue
			}

			if b.c.ServerOptions.Flush.Frequency > 0 && b.timer == nil {
				b.timer = time.NewTimer(b.c.ServerOptions.Flush.Frequency)
				timerChan = b.timer.C
			}
		case <-timerChan:
			b.timerFired = true
		case output <- b.buffer:
			b.rollOver()
			timerChan = nil

		}

		if b.timerFired || b.buffer.readyToFlush() {
			output = b.output
		} else {
			output = nil
		}
	}
}

func (b *Broker) sender() {
	for set := range b.output {
		if !set.empty() {
			req := set.buildRequest()
			_, err := b.cli.Push(b.ctx, req)
			if err != nil {
				zlog.Errorf("broker client push error:%v", err)
			}
		}
	}
}

func (b *Broker) rollOver() {
	if b.timer != nil {
		b.timer.Stop()
	}
	b.timer = nil
	b.timerFired = false
	b.buffer = newMessageSet(b.c)
}

//func (c *Connect) Push(req *pb.PushMsgReq) (err error) {
//	idx := atomic.AddUint64(&c.pushChNum, 1) % c.routineSize
//	c.pushCh[idx] <- req
//	return
//}

func (b *Broker) BroadcastRoom(req *pb.BroadcastRoomReq) (err error) {
	idx := atomic.AddUint64(&b.roomChNum, 1) % b.routineSize
	b.roomCh[idx] <- req
	return
}

func (b *Broker) Broadcast(req *pb.BroadcastReq) (err error) {
	b.broadcastCh <- req
	return
}

func (b *Broker) process(
	pushCh chan *pb.PushMsgReq,
	roomCh chan *pb.BroadcastRoomReq,
	broadcastCh chan *pb.BroadcastReq,
) {
	for {
		select {
		case req := <-pushCh:
			_, err := b.cli.PushMsg(context.Background(), &pb.PushMsgReq{
				ProtoOp: req.ProtoOp,
				Keys:    req.Keys,
				Msg:     req.Msg,
			})
			if err != nil {
				//log
			}

		case req := <-roomCh:
			_, err := b.cli.BroadcastRoom(context.Background(), &pb.BroadcastRoomReq{
				Msg:    req.Msg,
				RoomId: req.RoomId,
			})
			if err != nil {
				//log
			}

		case req := <-broadcastCh:
			_, err := b.cli.Broadcast(context.Background(), &pb.BroadcastReq{
				ProtoOp: req.ProtoOp,
				Msg:     req.Msg,
				Speed:   req.Speed,
			})
			if err != nil {
				//log
			}

		case <-b.ctx.Done():
			return
		}
	}
}
