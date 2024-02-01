package timingwheel

import (
	"container/list"
	"sync"
	"sync/atomic"
)

type bucket struct {
	expiration int64

	mutex sync.Mutex
	tasks *list.List //任务列表
}

func newBucket() *bucket {
	return &bucket{
		tasks:      list.New(),
		expiration: -1,
	}
}

func (b *bucket) Expiration() int64 {
	return atomic.LoadInt64(&b.expiration)
}

// SetExpiration 设置expiration 同时与旧值比较 若不同于旧值则返回true
// 用于避免在同一轮中, 同一个bucket被多次放入DelayQueue
func (b *bucket) SetExpiration(expiration int64) bool {
	return atomic.SwapInt64(&b.expiration, expiration) != expiration
}

// Add 添加任务
func (b *bucket) Add(task *Task) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	//将任务放入列表尾部
	e := b.tasks.PushBack(task)
	task.setBucket(b)
	task.element = e
}

// Remove 删除任务
func (b *bucket) Remove(task *Task) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.remove(task)
}

func (b *bucket) remove(task *Task) bool {
	if task.getBucket() != b { //任务不在该bucket中
		return false
	}
	//从任务列表中删除
	b.tasks.Remove(task.element)
	task.setBucket(nil)
	task.element = nil
	return true
}

func (b *bucket) Flush(reinsert func(*Task)) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	var (
		next *list.Element
		task *Task
	)
	//遍历任务列表
	for e := b.tasks.Front(); e != nil; {
		next = e.Next()

		task = e.Value.(*Task)
		b.remove(task) //将任务移出bucket
		reinsert(task) //重新加入到时间轮

		e = next
	}

	b.SetExpiration(-1)
}

func (b *bucket) FlushWithError(reinsert func(*Task) error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	var (
		next *list.Element
		task *Task
		err  error
	)
	//遍历任务列表
	for e := b.tasks.Front(); e != nil; {
		next = e.Next()

		task = e.Value.(*Task)
		b.remove(task) //删除任务
		err = reinsert(task)
		if err != nil {

		}

		e = next
	}

	b.SetExpiration(-1)
}
