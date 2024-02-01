package timingwheel

import (
	"errors"
	"sync/atomic"
	"time"
	"unsafe"
)

// TimingWheel 溢出型分层时间轮
// 多层时间轮使用同一个DelayQueue
type TimingWheel struct {
	tick int64 //每个时间格的基本时间跨度, 单位是Millisecond
	size int64 //时间格的个数

	interval    int64       //时间轮的总体跨度, 单位是Millisecond
	currentTime int64       //当前时间, 是tick的整数倍, 指向的时间格表示该时间格到期, 单位是Millisecond
	buckets     []*bucket   //每个bucket装着同一时间格的任务
	queue       *DelayQueue //延迟队列, 管理bucket, 返回到达执行时间的bucket

	//溢出时间轮
	overflowWheel unsafe.Pointer // type: *TimingWheel

	exitC chan struct{}
	wg    waitGroupWrapper
}

// NewTimingWheel size必须大于1, 因为下层轮盘的tick必须大于上层轮盘的tick, 若size为1, 则无法满足该条件, 此时若Task的延迟时间大于tick,
// 会导致不断地新增overflowWheel, 最终造成stack overflow, 实际上即使Task的延迟时间<=tick, 也可能造成stack overflow
// 由于currentTime会根据tickMs进行truncate, 所以一个延迟时间小于tick的Task可能不会立即执行
func NewTimingWheel(tick time.Duration, size int64) *TimingWheel {
	tickMs := int64(tick / time.Millisecond)
	if tickMs <= 0 {
		panic(errors.New("tick must be greater than or equal to 1ms"))
	} else if size <= 1 {
		panic(errors.New("size must be greater than 1"))
	}

	startMs := timeToMs(time.Now().UTC())

	return newTimingWheel(
		tickMs,
		size,
		startMs,
		NewDelayQueue(int(size)),
	)
}

func newTimingWheel(tickMs int64, size int64, startMs int64, queue *DelayQueue) *TimingWheel {
	buckets := make([]*bucket, size)
	for i := range buckets {
		buckets[i] = newBucket()
	}
	return &TimingWheel{
		tick:        tickMs,
		size:        size,
		currentTime: truncate(startMs, tickMs),
		interval:    tickMs * size,
		buckets:     buckets,
		queue:       queue,
		exitC:       make(chan struct{}),
	}
}

func (t *TimingWheel) Init(size int) {

}

// 添加任务
func (t *TimingWheel) add(task *Task) bool {
	//当前时间
	currentTime := atomic.LoadInt64(&t.currentTime)
	if task.expiration < currentTime+t.tick {
		// 到达过期时间 可以执行任务
		return false
	} else if task.expiration < currentTime+t.interval {
		// 未到达过期时间且在轮盘的刻度范围内
		id := task.expiration / t.tick
		b := t.buckets[id%t.size] //找到相应的bucket
		b.Add(task)               //将任务放入bucket

		//更新expiration
		if b.SetExpiration(id * t.tick) {
			//将bucket放入DelayQueue
			t.queue.Offer(b, b.Expiration())
		}

		return true
	} else {
		//不在轮盘的刻度范围内
		//获取下一层的时间轮
		overflowWheel := atomic.LoadPointer(&t.overflowWheel)
		if overflowWheel == nil {
			atomic.CompareAndSwapPointer(&t.overflowWheel,
				nil,
				unsafe.Pointer(newTimingWheel(
					t.interval,
					t.size,
					currentTime,
					t.queue,
				)),
			)
			overflowWheel = atomic.LoadPointer(&t.overflowWheel)
		}
		//将任务添加到下一层的时间轮
		return (*TimingWheel)(overflowWheel).add(task)
	}
}

func (t *TimingWheel) addOrRun(task *Task) {
	if !t.add(task) {
		go task.f()
	}
}

// 更新当前时间
func (t *TimingWheel) advanceClock(expiration int64) {
	currentTime := atomic.LoadInt64(&t.currentTime)
	if expiration >= currentTime+t.tick {
		currentTime = truncate(expiration, t.tick)
		atomic.StoreInt64(&t.currentTime, currentTime)

		// Try to advance the clock of the overflow wheel if present
		overflowWheel := atomic.LoadPointer(&t.overflowWheel)
		if overflowWheel != nil {
			(*TimingWheel)(overflowWheel).advanceClock(currentTime)
		}
	}
}

func (t *TimingWheel) Start() {
	t.wg.Wrap(func() {
		//轮询到达执行时间的任务
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
				t.advanceClock(b.Expiration())
				//执行bucket中到达执行时间的任务 未到达执行时间的任务重新加入
				b.Flush(t.addOrRun)
			case <-t.exitC:
				return
			}
		}
	})
}

func (t *TimingWheel) Stop() {
	close(t.exitC)
	t.wg.Wait()
}

func (t *TimingWheel) AfterFunc(d time.Duration, f func()) *Task {
	return t.Add(d, f)
}

func (t *TimingWheel) Add(d time.Duration, f func()) *Task {
	task := &Task{
		expiration: timeToMs(time.Now().UTC().Add(d)),
		f:          f,
	}
	t.addOrRun(task)
	return task
}

func (t *TimingWheel) Set(task *Task, d time.Duration) {
	task.Cancel()
	task.expiration = timeToMs(time.Now().UTC().Add(d))
	t.addOrRun(task)
}

type Scheduler interface {
	// Next returns the next execution time after the given (previous) time.
	// It will return a zero time if no next time is scheduled.
	//
	// All times must be UTC.
	Next(time.Time) time.Time
}

func (t *TimingWheel) ScheduleFunc(s Scheduler, f func()) (task *Task) {
	expiration := s.Next(time.Now().UTC())
	if expiration.IsZero() {
		// No time is scheduled, return nil.
		return
	}

	task = &Task{
		expiration: timeToMs(expiration),
		f: func() {
			// Schedule the task to execute at the next time if possible.
			expiration := s.Next(msToTime(task.expiration))
			if !expiration.IsZero() {
				task.expiration = timeToMs(expiration)
				t.addOrRun(task)
			}

			// Actually execute the task.
			f()
		},
	}
	t.addOrRun(task)

	return
}
