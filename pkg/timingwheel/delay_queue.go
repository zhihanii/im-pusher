package timingwheel

import (
	"container/heap"
	"sync"
	"sync/atomic"
	"time"
)

type item struct {
	Value    any
	Priority int64
	Index    int
}

// 优先队列 小根堆
type priorityQueue []*item

func newPriorityQueue(capacity int) priorityQueue {
	return make(priorityQueue, 0, capacity)
}

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *priorityQueue) Push(x any) {
	n := len(*pq)
	c := cap(*pq)
	if n+1 > c { //超过最大容量则扩容
		npq := make(priorityQueue, n, c*2)
		copy(npq, *pq)
		*pq = npq
	}
	*pq = (*pq)[0 : n+1]
	item := x.(*item)
	item.Index = n
	(*pq)[n] = item
}

func (pq *priorityQueue) Pop() any {
	n := len(*pq)
	c := cap(*pq)
	if n < (c/2) && c > 25 { //缩容
		npq := make(priorityQueue, n, c/2)
		copy(npq, *pq)
		*pq = npq
	}
	item := (*pq)[n-1]
	item.Index = -1
	*pq = (*pq)[0 : n-1]
	return item
}

// PeekAndShift 若堆顶元素的Priority<=max 则取出堆顶元素
// 获取不到元素的两种情况：一是堆中没有元素 二是堆顶元素的Priority>max
func (pq *priorityQueue) PeekAndShift(max int64) (*item, int64) {
	if pq.Len() == 0 {
		return nil, 0
	}

	item := (*pq)[0]
	if item.Priority > max {
		return nil, item.Priority - max
	}
	heap.Remove(pq, 0)

	return item, 0
}

type DelayQueue struct {
	//向外传递到达执行时间的bucket
	C chan any

	mutex sync.Mutex
	pq    priorityQueue

	// Similar to the sleeping state of runtime.timers.
	sleeping int32
	wakeupC  chan struct{}
}

func NewDelayQueue(size int) *DelayQueue {
	return &DelayQueue{
		C:       make(chan any),
		pq:      newPriorityQueue(size),
		wakeupC: make(chan struct{}),
	}
}

// Offer inserts the element into the current queue.
func (dq *DelayQueue) Offer(e any, expiration int64) {
	item := &item{Value: e, Priority: expiration}

	dq.mutex.Lock()
	heap.Push(&dq.pq, item)
	index := item.Index
	dq.mutex.Unlock()

	if index == 0 { //新插入的元素位于堆顶
		// A new item with the earliest expiration is added.
		//需要转变sleeping为0 并发送wakeup通知
		if atomic.CompareAndSwapInt32(&dq.sleeping, 1, 0) {
			dq.wakeupC <- struct{}{}
		}
	}
}

// Poll starts an infinite loop, in which it continually waits for an element
// to expire and then send the expired element to the channel C.
// 轮询检查堆顶的bucket是否到达执行时间
func (dq *DelayQueue) Poll(exitC chan struct{}, nowF func() int64) {
	for {
		now := nowF()

		dq.mutex.Lock()
		item, delta := dq.pq.PeekAndShift(now)
		if item == nil {
			//没有获取到符合条件的元素
			// No items left or at least one item is pending.

			// We must ensure the atomicity of the whole operation, which is
			// composed of the above PeekAndShift and the following StoreInt32,
			// to avoid possible race conditions between Offer and Poll.
			//转变sleeping状态:sleeping变为1
			atomic.StoreInt32(&dq.sleeping, 1)
		}
		dq.mutex.Unlock()

		if item == nil {
			if delta == 0 { //堆中没有元素
				// No items left.
				select {
				//等待有新元素添加到堆中
				case <-dq.wakeupC:
					// Wait until a new item is added.
					continue
				case <-exitC:
					goto exit
				}
			} else if delta > 0 { //堆顶元素的Priority>now
				// At least one item is pending.
				select {
				//若有新元素出现在堆顶 则会收到wakeup通知
				case <-dq.wakeupC:
					// A new item with an "earlier" expiration than the current "earliest" one is added.
					continue
				//否则等待delta时间 堆顶元素的Priority将<=now
				case <-time.After(time.Duration(delta) * time.Millisecond):
					// The current "earliest" item expires.

					// Reset the sleeping state since there's no need to receive from wakeupC.
					//将sleeping变为0
					//若此时DelayQueue的Offer方法被调用且新元素位于堆顶 则if条件为true, 需要消耗掉wakeup通知, 避免干扰
					if atomic.SwapInt32(&dq.sleeping, 0) == 0 {
						// A caller of Offer() is being blocked on sending to wakeupC,
						// drain wakeupC to unblock the caller.
						<-dq.wakeupC
					}
					continue
				case <-exitC:
					goto exit
				}
			}
		}

		select {
		//向外传递到达执行时间的bucket
		case dq.C <- item.Value:
			// The expired element has been sent out successfully.
		case <-exitC:
			goto exit
		}
	}

exit:
	// Reset the states
	atomic.StoreInt32(&dq.sleeping, 0)
}
