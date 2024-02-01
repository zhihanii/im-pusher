package client

type queue struct {
	list *arrayList
}

func newQueue() *queue {
	return &queue{list: newArrayList()}
}

func (q *queue) Enqueue(v *element) {
	q.list.Add(v)
}

func (q *queue) Dequeue() (v *element, ok bool) {
	v, ok = q.list.Get(0)
	if ok {
		q.list.Remove(0)
	}
	return
}

func (q *queue) Peek() (v *element, ok bool) {
	return q.list.Get(0)
}

func (q *queue) Empty() bool {
	return q.list.Empty()
}

func (q *queue) Size() int {
	return q.list.Size()
}

func (q *queue) Clear() {
	q.list.Clear()
}

func (q *queue) Values() []*element {
	return q.list.Values()
}
