package set

import "sync"

var itemExists = struct{}{}

type Set struct {
	mux   sync.RWMutex
	items map[uint64]struct{}
}

func New(values ...uint64) *Set {
	s := &Set{
		items: make(map[uint64]struct{}),
	}
	for _, v := range values {
		s.items[v] = itemExists
	}
	return s
}

func (s *Set) Remove(items ...uint64) {
	s.mux.Lock()
	for _, item := range items {
		delete(s.items, item)
	}
	s.mux.Unlock()
}

func (s *Set) Size() int {
	s.mux.RLock()
	res := len(s.items)
	s.mux.RUnlock()
	return res
}

func (s *Set) Values() []uint64 {
	s.mux.RLock()
	res := make([]uint64, len(s.items))
	for item := range s.items {
		res = append(res, item)
	}
	s.mux.RUnlock()
	return res
}
