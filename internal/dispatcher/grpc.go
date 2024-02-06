package dispatcher

//type grpcServer struct {
//	*grpc.Server
//	endpoint     *registry.Endpoint
//	etcdRegistry registry.Registry
//	//etcdClient *clientv3.Client
//}
//
//func (s *grpcServer) Run() {
//	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.endpoint.Address, s.endpoint.Port))
//	if err != nil {
//
//	}
//
//	go func() {
//		zlog.Infoln("开始运行grpc server")
//		if err := s.Serve(listener); err != nil {
//
//		}
//	}()
//
//	delayRegister := time.After(time.Second)
//	select {
//	case <-delayRegister:
//		if err := s.etcdRegistry.Register(s.endpoint); err != nil {
//			zlog.Errorf("register:%v", err)
//		}
//	}
//
//	//delayRegister := time.After(time.Second)
//	//<-delayRegister
//	//
//	//go func() {
//	//	var err1 error
//	//	ctx := context.Background()
//	//	em, err1 := endpoints.NewManager(s.etcdClient, s.endpoint.ServiceName)
//	//	if err1 != nil {
//	//		zlog.Errorf("endpoints.NewManager:%v", err)
//	//		return
//	//	}
//	//	lease, err1 := s.etcdClient.Grant(ctx, 15)
//	//	if err1 != nil {
//	//		zlog.Errorf("grant:%v", err)
//	//		return
//	//	}
//	//	err1 = em.AddEndpoint(ctx, fmt.Sprintf("%s/%s:%d", s.endpoint.ServiceName, s.endpoint.Address, s.endpoint.Port),
//	//		endpoints.Endpoint{
//	//			Addr: fmt.Sprintf("%s:%d", s.endpoint.Address, s.endpoint.Port),
//	//		}, clientv3.WithLease(lease.ID))
//	//	if err1 != nil {
//	//		zlog.Errorf("add endpoint:%v", err1)
//	//		return
//	//	}
//	//	ticker := time.NewTicker(5 * time.Second)
//	//	for {
//	//		select {
//	//		case <-ticker.C:
//	//			resp, err2 := s.etcdClient.KeepAliveOnce(ctx, lease.ID)
//	//			if err2 != nil {
//	//				zlog.Errorf("keep alive:%v", err2)
//	//			} else {
//	//				zlog.Infof("keep alive resp:%+v", resp)
//	//			}
//	//		case <-ctx.Done():
//	//			return
//	//		}
//	//	}
//	//}()
//}
//
//func (s *grpcServer) Close() {
//
//}
