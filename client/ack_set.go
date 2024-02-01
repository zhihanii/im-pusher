package client

import "bytes"

type ackSet struct {
	keys  [][]byte
	count int
}

func newAckSet() *ackSet {
	return &ackSet{}
}

func (s *ackSet) add(key []byte) error {
	s.keys = append(s.keys, key)
	s.count++
	return nil
}

func (s *ackSet) reset() {
	s.keys = [][]byte{}
	s.count = 0
}

func (s *ackSet) empty() bool {
	return s.count == 0
}

func (s *ackSet) bytes() []byte {
	return bytes.Join(s.keys, []byte{';'})
}
