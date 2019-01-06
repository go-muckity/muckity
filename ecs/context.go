package ecs

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

type emptyMktyCtx int

var _ MuckityContext = new(emptyMktyCtx)

func (emptyMktyCtx) Deadline() (deadline time.Time, ok bool)  	{ return }
func (emptyMktyCtx) Done() <-chan struct{}						{ return nil }
func (emptyMktyCtx) Err() error									{ return nil }
func (emptyMktyCtx) Value(key interface{}) interface{}			{ return nil }

var (
	background = new(emptyMktyCtx)
	todo = new(emptyMktyCtx)
)

func Background() MuckityContext	{ return background }
func TODO() MuckityContext			{ return todo }

type cancelCtx struct {
	MuckityContext
	done chan struct{}
	err error
	mu sync.Mutex
}

var _ MuckityContext = &cancelCtx{}

func (ctx *cancelCtx) Done() <-chan struct{}						{ return ctx.done }
func (ctx *cancelCtx) Err() error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	return ctx.err
}

var Canceled = errors.New("context canceled")

type CancelFunc func()

func WithCancel (parent MuckityContext) (MuckityContext, CancelFunc) {
	ctx := &cancelCtx{
		MuckityContext: parent,
		done: make(chan struct{}),
	}

	cancel := func() { ctx.cancel(Canceled) }

	go func() {
		select {
		case <-parent.Done():
			ctx.cancel(parent.Err())
		case <-ctx.Done():

		}
	}()

	return ctx, cancel
}

func (ctx *cancelCtx) cancel(err error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.err != nil {
		return
	}

	ctx.err = err
	close(ctx.done)
}

type deadlineCtx struct {
	*cancelCtx
	deadline 	time.Time
}

var _ MuckityContext = &deadlineCtx{}

func (ctx *deadlineCtx) Deadline() (deadline time.Time, ok bool) {
	return ctx.deadline, true
}

var DeadlineExceeded = errors.New("deadline exceeded")

func WithDeadline(parent MuckityContext, deadline time.Time) (MuckityContext, CancelFunc) {
	cctx, cancel := WithCancel(parent)

	ctx := &deadlineCtx{
		cancelCtx: cctx.(*cancelCtx),
		deadline: deadline,
	}

	t := time.AfterFunc(time.Until(deadline), func() {
		ctx.cancel(DeadlineExceeded)
	})

	stop := func() {
		t.Stop()
		cancel()
	}

	return ctx, stop
}

func WithTimeout(parent MuckityContext, timeout time.Duration) (MuckityContext, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

type valueCtx struct {
	MuckityContext
	value, key interface{}
}

var _ MuckityContext = &valueCtx{}

func (ctx *valueCtx) Value(key interface{}) interface{} {
	if key == ctx.key {
		return ctx.value
	}
	return ctx.MuckityContext.Value(key)
}

func WithValue(parent MuckityContext, key, value interface{}) MuckityContext {
	if key == nil {
		panic("key is nil")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &valueCtx{
		MuckityContext: parent,
		key: key,
		value: value,
	}
}