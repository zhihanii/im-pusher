package data

import (
	"context"
	"github.com/zhihanii/im-pusher/internal/store/conf"
	"github.com/zhihanii/im-pusher/internal/store/data/db"
	"github.com/zhihanii/zlog"
	"time"
)

type worker struct {
	c          *conf.Config
	repo       *storeRepo
	ctx        context.Context
	cancel     context.CancelFunc
	input      chan *db.ChatMessage
	output     chan *messageSet
	buffer     *messageSet
	timer      *time.Timer
	timerFired bool

	//msgCount int64
}

func newWorker(c *conf.Config, repo *storeRepo, input chan *db.ChatMessage) *worker {
	w := &worker{
		c:     c,
		repo:  repo,
		input: input,
	}
	w.ctx, w.cancel = context.WithCancel(context.Background())
	w.output = make(chan *messageSet)
	w.buffer = newMessageSet(w.c)
	go w.run()
	go w.sender()
	return w
}

func (w *worker) run() {
	var output chan<- *messageSet
	var timerChan <-chan time.Time
	//ticker := time.NewTicker(2 * time.Second)

	for {
		select {
		case msg, ok := <-w.input:
			if !ok {
				//zlog.Infoln("input chan closed")
				return
			}

			if msg == nil {
				continue
			}

			//w.msgCount++

			//todo if needs retry

			if err := w.buffer.add(msg); err != nil {
				continue
			}

			if w.c.ServerOptions.Flush.Frequency > 0 && w.timer == nil {
				w.timer = time.NewTimer(w.c.ServerOptions.Flush.Frequency)
				timerChan = w.timer.C
			}
		case <-timerChan:
			w.timerFired = true
		case output <- w.buffer:
			w.rollOver()
			timerChan = nil
			//case <-ticker.C:
			//	zlog.Infof("message count:%v", w.msgCount)
		}

		if w.timerFired || w.buffer.readyToFlush() {
			output = w.output
		} else {
			output = nil
		}
	}
}

func (w *worker) sender() {
	for set := range w.output {
		if !set.empty() {
			err := w.repo.storeBatchChatMessage(w.ctx, set.msgs)
			if err != nil {
				zlog.Errorf("store batch chat message:%v", err)
			}
		}
	}
}

func (w *worker) rollOver() {
	if w.timer != nil {
		w.timer.Stop()
	}
	w.timer = nil
	w.timerFired = false
	w.buffer = newMessageSet(w.c)
}
