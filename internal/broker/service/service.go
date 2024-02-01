package service

import (
	"context"
	pb "github.com/zhihanii/im-pusher/api/broker"
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/im-pusher/internal/broker/errors"
	"github.com/zhihanii/im-pusher/internal/broker/server"
	"time"
)

type ConnectService struct {
	pb.UnimplementedBrokerServer

	srv *server.Server
}

func New(srv *server.Server) (*ConnectService, error) {
	return &ConnectService{
		srv: srv,
	}, nil
}

func (s *ConnectService) Push(ctx context.Context, req *pb.PushReq) (reply *pb.PushReply, err error) {
	for _, tm := range req.Messages {
		for _, key := range tm.Keys {
			bucket := s.srv.Bucket(key)
			if bucket == nil {
				continue
			}
			if ch := bucket.Channel(key); ch != nil {
				//if !ch.NeedPush() {
				//	continue
				//}
				//todo 构造protocol.Message
				msg := &protocol.Message{
					Operation: tm.Operation,
					Sequence:  tm.Sequence,
					Data:      tm.Data,
				}

				//err = proto.Unmarshal(tm.Data, msg)
				//if err != nil {
				//	return &pb.PushReply{}, err
				//}

				if err = ch.Push(msg); err != nil {
					return
				}
			}
		}
	}
	return &pb.PushReply{}, nil
}

func (s *ConnectService) PushMsg(ctx context.Context, req *pb.PushMsgReq) (reply *pb.PushMsgReply, err error) {
	if len(req.Keys) == 0 || req.Msg == nil {
		return nil, errors.ErrPushMsgArg
	}
	for _, key := range req.Keys {
		bucket := s.srv.Bucket(key)
		if bucket == nil {
			continue
		}
		if ch := bucket.Channel(key); ch != nil {
			if !ch.NeedPush(req.ProtoOp) {
				continue
			}
			if err = ch.Push(req.Msg); err != nil {
				return
			}
		}
	}
	return &pb.PushMsgReply{}, nil
}

func (s *ConnectService) Broadcast(ctx context.Context, req *pb.BroadcastReq) (*pb.BroadcastReply, error) {
	if req.Msg == nil {
		return nil, errors.ErrBroadCastArg
	}
	// TODO use broadcast queue
	go func() {
		for _, bucket := range s.srv.Buckets() {
			bucket.Broadcast(req.GetMsg(), req.ProtoOp)
			if req.Speed > 0 {
				t := bucket.ChannelCount() / int(req.Speed)
				time.Sleep(time.Duration(t) * time.Second)
			}
		}
	}()
	return &pb.BroadcastReply{}, nil
}

func (s *ConnectService) BroadcastRoom(ctx context.Context, req *pb.BroadcastRoomReq) (*pb.BroadcastRoomReply, error) {
	if req.Msg == nil || req.RoomId == "" {
		return nil, errors.ErrBroadCastRoomArg
	}
	for _, bucket := range s.srv.Buckets() {
		bucket.BroadcastRoom(req)
	}
	return &pb.BroadcastRoomReply{}, nil
}

func (s *ConnectService) Rooms(ctx context.Context, req *pb.RoomsReq) (*pb.RoomsReply, error) {
	var roomIds = make(map[string]bool)
	for _, bucket := range s.srv.Buckets() {
		for roomID := range bucket.Rooms() {
			roomIds[roomID] = true
		}
	}
	return &pb.RoomsReply{Rooms: roomIds}, nil
}
