package muckity

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

type key int

const (
	rootKey    key = iota // root object pointer, of Type
	configKey             // Config interface
	worldKey              // World interface
	storageKey            // Storage interface
	systemKey             // System interface
)

type emptyMktyCtx int

var _ Context = new(emptyMktyCtx)

func (emptyMktyCtx) Deadline() (deadline time.Time, ok bool) { return }
func (emptyMktyCtx) Done() <-chan struct{}                   { return nil }
func (emptyMktyCtx) Err() error                              { return nil }
func (emptyMktyCtx) Value(key interface{}) interface{}       { return nil }
func (emptyMktyCtx) Config() Config                          { return nil }
func (emptyMktyCtx) World() World                            { return nil }
func (emptyMktyCtx) CallingSystem() System                   { return nil }
func (emptyMktyCtx) Name() string                            { return "emptyContext" }
func (emptyMktyCtx) Type() string                            { return "muckity:context" }

var (
	background = new(emptyMktyCtx)
	todo       = new(emptyMktyCtx)
)

func Background() Context { return background }
func TODO() Context       { return todo }

type cancelCtx struct {
	Context
	done chan struct{}
	err  error
	mu   sync.Mutex
}

var _ Context = &cancelCtx{}

func (ctx *cancelCtx) Done() <-chan struct{} { return ctx.done }
func (ctx *cancelCtx) Err() error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	return ctx.err
}

var Canceled = errors.New("context canceled")

type CancelFunc func()

func WithCancel(parent Context) (Context, CancelFunc) {
	ctx := &cancelCtx{
		parent,
		make(chan struct{}),
		nil,
		sync.Mutex{},
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
	deadline time.Time
}

var _ Context = &deadlineCtx{}

func (ctx *deadlineCtx) Deadline() (deadline time.Time, ok bool) {
	return ctx.deadline, true
}

var DeadlineExceeded = errors.New("deadline exceeded")

func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc) {
	cctx, cancel := WithCancel(parent)

	ctx := &deadlineCtx{
		cctx.(*cancelCtx),
		deadline,
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

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

type valueCtx struct {
	Context
	value, key interface{}
}

var _ Context = &valueCtx{}

func (ctx *valueCtx) Value(key interface{}) interface{} {
	if key == ctx.key {
		return ctx.value
	}
	return ctx.Value(key)
}

func WithValue(parent Context, key, value interface{}) Context {
	if key == nil {
		panic("key is nil")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &valueCtx{
		parent,
		value,
		key,
	}
}

func (ctx *valueCtx) Config() Config {
	return ctx.Value(configKey).(Config)
}

func (ctx *valueCtx) World() World {
	return ctx.Value(worldKey).(World)
}

func (ctx *valueCtx) CallingSystem() System {
	return ctx.Value(systemKey).(System)
}

func WithConfig(parent Context, value Config) Context {
	return WithValue(parent, configKey, value)
}

func WithWorld(parent Context, value World) Context {
	return WithValue(parent, worldKey, value)
}

func WithSystem(parent Context, value System) Context {
	return WithValue(parent, systemKey, value)
}

var rootCtx Context

func rootContext() Context {
	if rootCtx == nil {
		rootCtx = new(emptyMktyCtx)
	}
	return rootCtx
}

// GetContext returns a singleton context object, or an empty root context pointing at
// a default, unexported Background() context. Pass doOnce: true if you want the singleton
func GetContext(doOnce bool, ctx ...interface{}) Context {
	if doOnce {
		once.Do(func() {
			rootCtx = rootContext()
		})
	}
	return rootContext()
}
