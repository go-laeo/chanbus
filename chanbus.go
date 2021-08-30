package chanbus

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var (
	ErrSendTimeout = errors.New("send timeout")
)

type CancelFunc func()

type Chanbus interface {
	Send(v interface{})
	SendTimeout(v interface{}, timeout time.Duration) error
	SendContext(ctx context.Context, v interface{}) error
	Derive(size uint) (<-chan interface{}, CancelFunc)
	Close()
}

type chanbus struct {
	ch  chan interface{}
	to  time.Duration
	chs atomic.Value
}

func New(size uint, timeout time.Duration) Chanbus {
	cb := chanbus{
		ch:  make(chan interface{}, size),
		to:  timeout,
		chs: atomic.Value{},
	}

	cb.chs.Store(ChanList{})

	go cb.forwarding()

	return &cb
}

func (cb *chanbus) forwarding() {
	for v := range cb.ch {
		for _, ch := range cb.chs.Load().(ChanList) {
			select {
			case <-time.After(cb.to):
				continue
			case ch <- v:
			}
		}
	}
}

func (cb *chanbus) Send(v interface{}) {
	cb.ch <- v
}

func (cb *chanbus) SendTimeout(v interface{}, timeout time.Duration) error {
	select {
	case <-time.After(timeout):
		return ErrSendTimeout
	case cb.ch <- v:
		return nil
	}
}

func (cb *chanbus) SendContext(ctx context.Context, v interface{}) error {
	select {
	case <-ctx.Done():
		return ErrSendTimeout
	case cb.ch <- v:
		return nil
	}
}

func (cb *chanbus) Derive(size uint) (<-chan interface{}, CancelFunc) {
	ch := make(chan interface{}, size)
	v := cb.chs.Load().(ChanList)
	v = append(v, ch)
	cb.chs.Store(v)

	return ch, func() {
		v := cb.chs.Load().(ChanList)
		n := v.IndexOf(ch)
		if n > -1 {
			v = append(v[:n], v[n+1:]...)
			cb.chs.Store(v)
		}
	}
}

func (cb *chanbus) Close() {
	close(cb.ch)
	for _, ch := range cb.chs.Load().(ChanList) {
		close(ch)
	}
}
