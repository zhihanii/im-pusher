package client

import "github.com/zhihanii/im-pusher/pkg/timingwheel"

const (
	growthFactor = float32(2.0)  // 扩容因子 两倍扩容
	shrinkFactor = float32(0.25) // 缩容因子 若当前元素个数小于等于当前容量的0.25, 则缩容
)

type element struct {
	sequence uint32
	task     *timingwheel.Task
	success  chan struct{}
}

type arrayList struct {
	elements []*element
	size     int
}

func newArrayList(values ...*element) *arrayList {
	list := &arrayList{}
	if len(values) > 0 {
		list.Add(values...)
	}
	return list
}

func (l *arrayList) Add(values ...*element) {
	l.grow(len(values))
	for _, v := range values {
		l.elements[l.size] = v
		l.size++
	}
}

func (l *arrayList) Get(i int) (*element, bool) {
	if !l.withinRange(i) {
		return nil, false
	}
	if l.size == 0 {
		return nil, false
	}
	return l.elements[i], true
}

func (l *arrayList) Remove(i int) {
	if !l.withinRange(i) {
		return
	}
	if l.size == 0 {
		return
	}
	l.elements[i] = nil
	//将位置从i开始的元素往后移动1个单位
	copy(l.elements[i:], l.elements[i+1:l.size])
	l.size--
	l.shrink()
}

func (l *arrayList) Contains(v *element) bool {
	for i := 0; i < l.size; i++ {
		if l.elements[i] == v {
			return true
		}
	}
	return false
}

func (l *arrayList) Values() []*element {
	values := make([]*element, l.size, l.size)
	copy(values, l.elements[:l.size])
	return values
}

func (l *arrayList) IndexOf(v *element) int {
	if l.size == 0 {
		return -1
	}
	for i, e := range l.elements {
		if e == v {
			return i
		}
	}
	return -1
}

func (l *arrayList) Empty() bool {
	return l.size == 0
}

func (l *arrayList) Size() int {
	return l.size
}

func (l *arrayList) Clear() {
	l.size = 0
	l.elements = []*element{}
}

func (l *arrayList) Swap(i, j int) {
	if l.withinRange(i) && l.withinRange(j) {
		l.elements[i], l.elements[j] = l.elements[j], l.elements[i]
	}
}

func (l *arrayList) Insert(i int, values ...*element) {
	if !l.withinRange(i) {
		if i == l.size {
			l.Add(values...)
		}
		return
	}
	length := len(values)
	l.grow(length)
	l.size += length
	//将位置从i开始的元素往后移动length个单位
	copy(l.elements[i+length:], l.elements[i:l.size-length])
	//插入元素
	copy(l.elements[i:], values)
}

func (l *arrayList) Set(i int, v *element) {
	if !l.withinRange(i) {
		if i == l.size {
			l.Add(v)
		}
		return
	}
	l.elements[i] = v
}

func (l *arrayList) withinRange(i int) bool {
	return i >= 0 && i < l.size
}

// 调整容量
func (l *arrayList) resize(capacity int) {
	newElements := make([]*element, capacity, capacity)
	copy(newElements, l.elements)
	l.elements = newElements
}

// 扩容
func (l *arrayList) grow(n int) {
	currentCapacity := cap(l.elements)
	if l.size+n >= currentCapacity {
		newCapacity := int(growthFactor * float32(currentCapacity+n))
		l.resize(newCapacity)
	}
}

// 缩容
func (l *arrayList) shrink() {
	if shrinkFactor == 0.0 {
		return
	}
	currentCapacity := cap(l.elements)
	if l.size <= int(float32(currentCapacity)*shrinkFactor) {
		l.resize(l.size)
	}
}
