package service

import (
	"context"
	"errors"
	pb "github.com/zhihanii/im-pusher/api/dispatcher"
	"github.com/zhihanii/im-pusher/internal/dispatcher/biz"
	"github.com/zhihanii/zlog"
)

type DispatchService struct {
	pb.UnimplementedDispatcherServer

	dh *biz.DispatcherHandler
}

func New(dh *biz.DispatcherHandler) (*DispatchService, error) {
	svc := &DispatchService{
		dh: dh,
	}
	return svc, nil
}

func (s *DispatchService) Push(ctx context.Context, req *pb.PushReq) (*pb.PushReply, error) {
	var err error
	var offlineMembers []uint64
	for _, receiver := range req.Receivers {
		err = s.dh.HandlePush(ctx, receiver, req.Message)
		if err != nil {
			zlog.Errorf("handle push:%v", err)
			if errors.Is(err, biz.MemberOfflineErr) {
				offlineMembers = append(offlineMembers, receiver)
			}
		}
	}
	return &pb.PushReply{
		OfflineMembers: offlineMembers,
	}, nil
}
