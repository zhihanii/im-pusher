package server

import (
	"github.com/zhihanii/im-pusher/internal/broker/conf"
	"github.com/zhihanii/im-pusher/pkg/timingwheel"
	"sync"
	"time"
)

// RoundOptions round options.
type RoundOptions struct {
	Timer        int
	TimerSize    int
	Reader       int
	ReadBuf      int
	ReadBufSize  int
	Writer       int
	WriteBuf     int
	WriteBufSize int
}

// Round used for connection round-robin get a reader/writer/timer for split big lock.
type Round struct {
	//readers []bytes.Pool
	//writers []bytes.Pool
	idx         int
	timersMutex sync.Mutex
	timers      []*timingwheel.TimingWheel
	options     RoundOptions
}

// NewRound new a round struct.
func NewRound(c *conf.Config) (r *Round) {
	var i int
	r = &Round{
		options: RoundOptions{
			//Reader:       c.TCP.Reader,
			//ReadBuf:      c.TCP.ReadBuf,
			//ReadBufSize:  c.TCP.ReadBufSize,
			//Writer:       c.TCP.Writer,
			//WriteBuf:     c.TCP.WriteBuf,
			//WriteBufSize: c.TCP.WriteBufSize,
			//Timer:        c.Protocol.Timer,
			Timer: 5,
			//TimerSize:    c.Protocol.TimerSize,
			TimerSize: 60,
		},
		idx:         0,
		timersMutex: sync.Mutex{},
		//timers:      make([]timingwheel.TimingWheel, r.options.Timer),
	}
	// reader
	//r.readers = make([]bytes.Pool, r.options.Reader)
	//for i = 0; i < r.options.Reader; i++ {
	//	r.readers[i].Init(r.options.ReadBuf, r.options.ReadBufSize)
	//}
	//// writer
	//r.writers = make([]bytes.Pool, r.options.Writer)
	//for i = 0; i < r.options.Writer; i++ {
	//	r.writers[i].Init(r.options.WriteBuf, r.options.WriteBufSize)
	//}
	// timer
	r.timers = make([]*timingwheel.TimingWheel, r.options.Timer)
	for i = 0; i < r.options.Timer; i++ {
		//r.timers[i].Init(r.options.TimerSize)
		r.timers[i] = timingwheel.NewTimingWheel(time.Second, int64(r.options.TimerSize))
		r.timers[i].Start()
	}
	return
}

// Timer get a timer.
func (r *Round) Timer() (t *timingwheel.TimingWheel) {
	r.timersMutex.Lock()
	t = r.timers[r.idx]
	r.idx = (r.idx + 1) % r.options.Timer
	r.timersMutex.Unlock()
	return
}

// Reader get a reader memory buffer.
//func (r *Round) Reader(rn int) *bytes.Pool {
//	return &(r.readers[rn%r.options.Reader])
//}
//
//// Writer get a writer memory buffer pool.
//func (r *Round) Writer(rn int) *bytes.Pool {
//	return &(r.writers[rn%r.options.Writer])
//}
