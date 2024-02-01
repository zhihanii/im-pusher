package timingwheel

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
	"unsafe"
)

type ExactTimingWheelConfig struct {
	tick time.Duration
	size int64
}

type ExactTimingWheel struct {
	levelId int64
	isLast  bool //是否是最下层的时间轮

	tick int64
	size int64

	interval    int64
	currentTime int64
	buckets     []*bucket
	queue       *DelayQueue

	nextWheel unsafe.Pointer // type: *TimingWheel

	exitC chan struct{}
	wg    waitGroupWrapper
}

func NewExactTimingWheel(level int, configs []*ExactTimingWheelConfig) (*ExactTimingWheel, error) {
	if len(configs) == 0 {
		return nil, errors.New("len of configs must be greater than zero")
	}
	if level != len(configs) {
		return nil, errors.New("len of configs must be equal to level")
	}
	var (
		root, p, tmp *ExactTimingWheel
		tickMs       int64
	)
	startMs := timeToMs(time.Now().UTC())
	tickMs = int64(configs[0].tick / time.Millisecond)
	root = newExactTimingWheel(tickMs, configs[0].size, startMs, NewDelayQueue(int(configs[0].size)))
	p = root
	for i := 1; i < level; i++ {
		tickMs = int64(configs[i].tick / time.Millisecond)
		tmp = newExactTimingWheel(tickMs, configs[i].size, startMs, NewDelayQueue(int(configs[i].size)))
		p.nextWheel = unsafe.Pointer(tmp)
		p = tmp
	}
	return root, nil
}

func newExactTimingWheel(tickMs int64, size int64, startMs int64, queue *DelayQueue) *ExactTimingWheel {
	buckets := make([]*bucket, size)
	for i := range buckets {
		buckets[i] = newBucket()
	}
	return &ExactTimingWheel{
		tick:        tickMs,
		size:        size,
		currentTime: truncate(startMs, tickMs),
		interval:    tickMs * size,
		buckets:     buckets,
		queue:       queue,
		exitC:       make(chan struct{}),
	}
}

func (t *ExactTimingWheel) Start() {
	t.wg.Wrap(func() {
		t.queue.Poll(t.exitC, func() int64 {
			return timeToMs(time.Now().UTC())
		})
	})

	t.wg.Wrap(func() {
		for {
			select {
			//取出到达执行时间的bucket
			case e := <-t.queue.C:
				b := e.(*bucket)
				//t.advanceClock(b.Expiration())

				if t.isLast {
					//执行到达执行时间的任务 否则重新加入
					b.Flush(t.addOrRun)
				} else {
					nextWheel := atomic.LoadPointer(&t.nextWheel)
					if nextWheel != nil {
						//将任务放入下一层时间轮
						b.FlushWithError((*ExactTimingWheel)(nextWheel).add)
					}
				}

			case <-t.exitC:
				return
			}
		}
	})
}

func (t *ExactTimingWheel) Stop() {
	close(t.exitC)
	t.wg.Wait()
}

// 添加任务
func (t *ExactTimingWheel) add(task *Task) error {
	//当前时间
	currentTime := atomic.LoadInt64(&t.currentTime)
	if task.expiration < currentTime+t.interval {
		// 在轮盘的刻度范围内
		id := task.expiration / t.tick
		b := t.buckets[id%t.size] //找到相应的bucket
		b.Add(task)               //将任务放入bucket

		//更新expiration
		if b.SetExpiration(id * t.tick) {
			//往bucket中放入任务可能使bucket的expiration提前
			//所以需要在expiration改变时调用Offer
			t.queue.Offer(b, b.Expiration())
		}

		return nil
	} else {
		//不在轮盘的刻度范围内
		return fmt.Errorf("overflow levelId: %d", t.levelId)
	}
}

func (t *ExactTimingWheel) addOrRun(task *Task) {
	currentTime := atomic.LoadInt64(&t.currentTime)
	if task.expiration < currentTime+t.tick {
		go task.f()
	} else {
		_ = t.add(task)
	}
}
