package server

//var (
//	MemberOfflineErr = errors.New("member offline")
//)
//
//type DispatcherHandler struct {
//	c *conf.Config
//
//	s *MessageSender
//
//	repo data.DispatcherRepo
//}
//
//func NewDispatcherHandler(c *conf.Config, s *MessageSender, repo data.DispatcherRepo) (*DispatcherHandler, error) {
//	dh := &DispatcherHandler{
//		c:    c,
//		s:    s,
//		repo: repo,
//	}
//	return dh, nil
//}
//
//func (h *DispatcherHandler) HandlePush(ctx context.Context, receiverId uint64, dm *protocol.DispatcherMessage) (err error) {
//	//receiverId := dm.Receivers[0]
//	zlog.Infof("handle push, receiver id:%d", receiverId)
//	res, err := h.repo.GetMapping(ctx, receiverId)
//	if err != nil {
//		zlog.Errorf("get mapping: %v", err)
//		return err
//	}
//	if len(res) == 0 {
//		zlog.Infof("member[%d] offline", receiverId)
//		return MemberOfflineErr
//	}
//	//serverKeys := make(map[string][]string)
//	//for key, server := range res {
//	//	if key != "" && server != "" {
//	//		serverKeys[server] = append(serverKeys[server], key)
//	//	}
//	//}
//	//for server, keys := range serverKeys {
//	//	_ = h.push(ctx, protocol.Channel, 0, server, keys, dm.Operation, dm.Sequence, dm.Data)
//	//}
//	return
//}
//
//func (h *DispatcherHandler) push(ctx context.Context, t uint32, priority int, server string, keys []string,
//	operation uint32, sequence uint32, data []byte) (err error) {
//	tm := &protocol.TransMessage{
//		Type:      t,
//		Priority:  uint32(priority),
//		Server:    server,
//		Keys:      keys,
//		Operation: operation,
//		Sequence:  sequence,
//		Data:      data,
//	}
//	h.s.send(tm)
//	return nil
//}
