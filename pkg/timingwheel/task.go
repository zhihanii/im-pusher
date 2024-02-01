package timingwheel

import (
	"container/list"
	"sync/atomic"
	"unsafe"
)

type Task struct {
	expiration int64
	f          func()

	b unsafe.Pointer // type: *bucket

	element *list.Element
}

// Cancel 取消任务
func (t *Task) Cancel() {
	b := t.getBucket()
	if b != nil {
		b.Remove(t)
	}
}

func (t *Task) getBucket() *bucket {
	return (*bucket)(atomic.LoadPointer(&t.b))
}

func (t *Task) setBucket(b *bucket) {
	atomic.StorePointer(&t.b, unsafe.Pointer(b))
}
