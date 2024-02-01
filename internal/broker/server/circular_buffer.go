package server

import (
	"github.com/zhihanii/im-pusher/api/protocol"
	"github.com/zhihanii/im-pusher/internal/broker/errors"
)

type CircularBuffer struct {
	rp   uint64
	num  uint64
	mask uint64
	wp   uint64
	//Message使用完后要将Data置为nil, 避免memory leak
	data []protocol.Message
}

func NewCircularBuffer(num int) *CircularBuffer {
	b := new(CircularBuffer)
	b.init(uint64(num))
	return b
}

func (b *CircularBuffer) Init(num int) {
	b.init(uint64(num))
}

func (b *CircularBuffer) init(num uint64) {
	if num&(num-1) != 0 { //若num不是2的幂
		for num&(num-1) != 0 { //将num规整为2的幂
			num &= num - 1
		}
		num <<= 1 //由于规整后的num比原来的小, 所以左移1位使num比原来大
	}
	//log.Printf("num:%d\n", num)
	b.data = make([]protocol.Message, num)
	b.num = num
	b.mask = b.num - 1
}

func (b *CircularBuffer) Get() (m *protocol.Message, err error) {
	if b.rp == b.wp {
		return nil, errors.ErrCircularBufferEmpty
	}
	m = &b.data[b.rp&b.mask]
	return
}

func (b *CircularBuffer) GetAdv() {
	b.rp++
}

func (b *CircularBuffer) Set() (m *protocol.Message, err error) {
	if b.wp-b.rp >= b.num {
		return nil, errors.ErrCircularBufferFull
	}
	m = &b.data[b.wp&b.mask]
	return
}

func (b *CircularBuffer) SetAdv() {
	b.wp++
}

func (b *CircularBuffer) Reset() {
	b.rp = 0
	b.wp = 0
}
