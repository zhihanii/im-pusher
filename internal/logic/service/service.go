package service

import (
	"context"
	pb "github.com/zhihanii/im-pusher/api/logic"
	"github.com/zhihanii/im-pusher/internal/logic/biz"
)

type LogicService struct {
	pb.UnimplementedLogicServer

	lh *biz.LogicHandler
}

func New(lh *biz.LogicHandler) (*LogicService, error) {
	//lh, err := biz.NewLogicHandler(c)
	//if err != nil {
	//	return nil, err
	//}
	svc := &LogicService{
		lh: lh,
	}
	return svc, nil
}

func (s *LogicService) Connect(ctx context.Context, req *pb.ConnectReq) (*pb.ConnectReply, error) {
	memberId, key, roomId, accepts, hb, err := s.lh.Connect(ctx, req.Server, req.Token)
	if err != nil {
		return &pb.ConnectReply{}, err
	}
	return &pb.ConnectReply{
		MemberId:  memberId,
		Key:       key,
		RoomId:    roomId,
		Accepts:   accepts,
		Heartbeat: hb,
	}, nil
}

func (s *LogicService) Disconnect(ctx context.Context, req *pb.DisconnectReq) (*pb.DisconnectReply, error) {
	has, err := s.lh.Disconnect(ctx, req.MemberId, req.Key, req.Server)
	if err != nil {
		return &pb.DisconnectReply{}, err
	}
	return &pb.DisconnectReply{
		Has: has,
	}, nil
}

func (s *LogicService) Heartbeat(ctx context.Context, req *pb.HeartbeatReq) (*pb.HeartbeatReply, error) {
	if err := s.lh.Heartbeat(ctx, req.MemberId, req.Key, req.Server); err != nil {
		return &pb.HeartbeatReply{}, err
	}
	return &pb.HeartbeatReply{}, nil
}

func (s *LogicService) RenewOnline(ctx context.Context, req *pb.OnlineReq) (*pb.OnlineReply, error) {
	allRoomCount, err := s.lh.RenewOnline(ctx, req.Server, req.RoomCount)
	if err != nil {
		return &pb.OnlineReply{}, err
	}
	return &pb.OnlineReply{AllRoomCount: allRoomCount}, nil
}

func (s *LogicService) Nodes(ctx context.Context, req *pb.NodesReq) (*pb.NodesReply, error) {
	//return s.lh.NodesWeighted(ctx, req.Platform, req.ClientIp), nil
	return &pb.NodesReply{}, nil
}

func (s *LogicService) PushKeys(ctx context.Context, req *pb.PushKeysReq) (*pb.PushKeysReply, error) {
	//err := s.lh.PushKeys(ctx, req.Op, req.Keys, req.Msg)
	//if err != nil {
	//	return &pb.PushKeysReply{}, err
	//}
	return &pb.PushKeysReply{}, nil
}

func (s *LogicService) PushMids(ctx context.Context, req *pb.PushMidsReq) (*pb.PushMidsReply, error) {
	//err := s.lh.PushMids(ctx, req.Op, req.Mids, req.Msg)
	//if err != nil {
	//	return &pb.PushMidsReply{}, err
	//}
	return &pb.PushMidsReply{}, nil
}

func (s *LogicService) PushRoom(ctx context.Context, req *pb.PushRoomReq) (*pb.PushRoomReply, error) {
	//err := s.lh.PushRoom(ctx, req.Op, req.Type, req.Room, req.Msg)
	//if err != nil {
	//	return &pb.PushRoomReply{}, err
	//}
	return &pb.PushRoomReply{}, nil
}

func (s *LogicService) PushAll(ctx context.Context, req *pb.PushAllReq) (*pb.PushAllReply, error) {
	//err := s.lh.PushAll(ctx, req.Op, req.Speed, req.Msg)
	//if err != nil {
	//	return &pb.PushAllReply{}, err
	//}
	return &pb.PushAllReply{}, nil
}

func (s *LogicService) Receive(ctx context.Context, req *pb.ReceiveReq) (*pb.ReceiveReply, error) {
	ackMsg, err := s.lh.Receive(ctx, req.MemberId, req.Message)
	if err != nil {
		return &pb.ReceiveReply{}, err
	}
	return &pb.ReceiveReply{
		Message: ackMsg,
	}, nil
}

func (s *LogicService) Push(ctx context.Context, req *pb.PushReq) (*pb.PushReply, error) {
	return &pb.PushReply{}, nil
}
